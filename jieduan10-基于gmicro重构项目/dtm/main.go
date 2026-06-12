package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v3"

	"github.com/dtm-labs/client/dtmcli"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
)

type UserAccount struct {
	ID             int     `gorm:"column:id;primary_key"`
	UserId         int     `gorm:"user_id"`
	Balance        float64 `gorm:"balance"`
	TradingBalance float64 `gorm:"trading_balance"`
}

func (UserAccount) TableName() string {
	return "user_account"
}

var lock sync.Mutex

// 转入和转出的时候，都要加锁，否则会出现并发问题
// SAGA 模式下「账户余额调整」通用接口，同时承担正向扣款 / 补偿退款能力
// 入参：
// db *sql.Tx：数据库事务（由 SAGA 协调器统一管控）
// uid：用户 ID
// amount：变动金额
// amount < 0：正向扣款（扣余额）
// amount > 0：补偿退款（加余额）
// 逻辑：
// 全局互斥锁加锁，防止并发修改同一用户余额
// 扣款前校验余额是否充足
// 执行 update 完成余额变更
// 定位：
// 可作为 SAGA 流程里正向动作或补偿动作，一套接口复用。
func SagaAdjustBalance(db *sql.Tx, uid int, amount float64) error {
	// db *sql.Tx 本身就是数据库事务对象，QueryRow、Exec 全部在当前事务中执行
	lock.Lock()
	defer lock.Unlock()

	if amount < 0 {
		var balance float64
		db.QueryRow("select balance from dtm.user_account where user_id = ?", uid).Scan(&balance)
		if balance < -amount {
			return fmt.Errorf("余额不足")
		}
	}
	_, err := db.Exec("update dtm.user_account set balance = balance + ? where user_id = ?", amount, uid)
	if err != nil {
		return err
	}
	return nil
}

var db *gorm.DB

// 实例化gorm
func initDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"root",
		"127.0.0.1",
		"3306",
		"dtm")
	newLogger := glog.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		glog.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  glog.Info,   // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,       // 禁用彩色打印
		},
	)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}
	return nil
}

// MustBarrierFromGin 1
// 用于 Gin 接口中快速获取 DTM 子事务屏障对象。
// DTM 在调用分支接口时，会通过 URL Query 参数 传递事务标识（gid/branchid/op 等），BarrierFromQuery 从请求参数解析出这些信息，生成屏障实例
func MustBarrierFromGin(c *gin.Context) *dtmcli.BranchBarrier {
	ti, err := dtmcli.BarrierFromQuery(c.Request.URL.Query())
	fmt.Println(err)
	return ti
}

// 服务发现， 库存服务有5个
func main() {
	err := initDB()
	if err != nil {
		panic(err)
	}

	// 开始执行gin接口
	r := gin.Default()
	// 转入接口
	r.POST("/SagaBTransIn", func(c *gin.Context) {
		// 获取子事务屏障
		barrier := MustBarrierFromGin(c)
		tx := db.Begin()
		sourceTx := tx.Statement.ConnPool.(*sql.Tx)
		err := barrier.Call(sourceTx, func(tx1 *sql.Tx) error {
			fmt.Println("开始转入")
			userID := 1
			// 转入100
			err := SagaAdjustBalance(sourceTx, userID, 100)
			if err != nil {
				fmt.Printf("转入失败:%s\r\n", err.Error())
				return err
			}
			return nil
		})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
			return
		}

		return
	})
	// 转入补偿接口
	r.POST("/SagaBTransInCom", func(c *gin.Context) {
		// 这块也得获取子事务屏障，代码后补
		fmt.Println("转入失败， 开始补偿")
		userID := 1
		err := SagaAdjustBalance(db, userID, -100)
		if err != nil {
			fmt.Printf("转入补偿失败:%s\r\n", err.Error())
			return
		}
		fmt.Println("转入补偿成功")
	})
	// 转出接口
	r.POST("/SagaBTransOut", func(c *gin.Context) {
		barrier := MustBarrierFromGin(c)
		tx := db.Begin()
		sourceTx := tx.Statement.ConnPool.(*sql.Tx)

		err := barrier.Call(sourceTx, func(tx1 *sql.Tx) error {
			fmt.Println("开始转出")
			userID := 3
			err := SagaAdjustBalance(sourceTx, userID, -100)
			if err != nil {
				// 切记失败时，一定要返回状态码：否则dtm服务不知道你这个接口成功还是失败
				if err.Error() == "余额不足" {
					// 返回这个状态码代表表明失败，直接触发补偿
					c.JSON(http.StatusConflict, gin.H{})
				}
				// 切记失败时，一定要返回状态码：否则dtm服务不知道你这个接口成功还是失败
				fmt.Printf("转出失败:%s\r\n", err.Error())
				// 代表ongoing状态：表明重试
				c.JSON(500, gin.H{"msg": err.Error()})
			}
			fmt.Println("转出成功")
			return nil
		})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
			return
		}
		return
	})
	// 转出补偿接口
	r.POST("/SagaBTransOutCom", func(c *gin.Context) {
		fmt.Println("转出失败， 开始补偿")
		userID := 3
		// 业务上要更严谨，没有转出不能补偿
		err := SagaAdjustBalance(db, userID, 100)
		if err != nil {
			fmt.Printf("转出补偿失败:%s\r\n", err.Error())
			return
		}
		fmt.Println("转出补偿成功")
	})
	// 开启saga事务
	r.GET("start", func(c *gin.Context) {
		req := gin.H{}
		// DTM 服务端的 API 入口地址
		dmtServer := "http://127.0.0.1:36789/api/dtmsvr"
		qsBusi := "http://127.0.0.1:8089"
		saga := dtmcli.NewSaga(dmtServer, shortuuid.New()).
			// 顺序：user3先转出到user1，user1再转入
			// 添加一个TransOut的子事务，正向操作为url: qsBusi+"/TransOut"， 逆向操作为url: qsBusi+"/TransOutCom"
			// 3参为请求参数
			Add(qsBusi+"/SagaBTransOut", qsBusi+"/SagaBTransOutCom", req).
			// 添加一个TransIn的子事务，正向操作为url: qsBusi+"/TransOut"， 逆向操作为url: qsBusi+"/TransInCom"
			// 3参为请求参数
			Add(qsBusi+"/SagaBTransIn", qsBusi+"/SagaBTransInCom", req)
		// 提交saga事务，dtm会完成所有的子事务/回滚所有的子事务
		saga.WaitResult = true // 等待事务的结果，否则默认是异步的，下面拿不到结果
		err := saga.Submit()
		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
		}
		c.JSON(200, gin.H{"message": "ok"})
	})

	r.Run(":8089")
}
