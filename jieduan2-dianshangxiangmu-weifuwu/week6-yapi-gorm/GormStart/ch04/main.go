package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	ID           uint
	Name         string
	Email        *string
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullString
	ActivedAt    sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:root@tcp(192.168.0.104:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	// 全局模式
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	//单一的 SQL 语句
	// 构建了3条数据
	var users = []User{{Name: "bobby1"}, {Name: "bobby2"}, {Name: "bobby3"}}
	// 批量方式01- 批量插入users数据，这种没有限制，当数据量比较大的时候，应该用下面的
	db.Create(&users)

	// 批量方式02- 批量插入：为什么不一次性提交所有的 还要分批次，--因为sql语句有长度限制
	db.CreateInBatches(users, 10000) // 每批次上限插入10000条数据

	for _, user := range users {
		fmt.Println(user.ID) // 1,2,3
	}
	// 批量方式03- 根据使用map直接创建
	db.Model(&User{}).Create(map[string]interface{}{
		"Name": "jinzhu", "Age": 18,
	})

	// 关联创建：是一个Create命令会帮你更新到2个表中，
	// 最终动作：插入 users 表 1 条数据，插入 credit_cards 表 1 条数据
	// https://gorm.io/zh_CN/docs/create.html#%E5%85%B3%E8%81%94%E5%88%9B%E5%BB%BA
	type CreditCard struct {
		gorm.Model
		Number string
		UserID uint // 关联 users 表的外键，自动通过约定的命名规范识别的
	}

	type User struct {
		gorm.Model
		Name       string
		CreditCard CreditCard // 关联了CreditCard这张表
	}

	db.Create(&User{
		Name:       "jinzhu",
		CreditCard: CreditCard{Number: "411111111111"},
	})
	// INSERT INTO `users` ...
	// INSERT INTO `credit_cards` ...

	// 默认值
	// 参考文档
}
