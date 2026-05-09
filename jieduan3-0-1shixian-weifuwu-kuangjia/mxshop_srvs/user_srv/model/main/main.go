package main

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mxshop_srvs/user_srv/model"
	"os"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// md5加密的方法例子
func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}

func main() {
	dsn := "root:root@tcp(192.168.0.104:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"

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
		// 关闭表名自动加 s 复数，默认：GORM 会自动给表名加 s
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		// 自定义日志，用来打印 SQL 语句、慢查询、错误日志
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	// PBKDF2 加密算法的参数（PBKDF2 是行业标准，比 MD5 安全 10000 倍。）
	// 16：salt 长度（随机盐）
	// 100：迭代次数（越慢越安全）
	// 32：生成的密钥长度
	// sha512.New：哈希算法（SHA512，超级安全）
	options := &password.Options{16, 100, 32, sha512.New}
	// 把明文密码 admin123加密
	salt, encodedPwd := password.Encode("admin123", options)
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	fmt.Println(newPassword)

	// 循环生成10个用户：一次性插入 10 条用户数据，密码都是 admin123（加密存储）
	for i := 0; i < 10; i++ {
		user := model.User{
			NickName: fmt.Sprintf("bobby%d", i),
			Mobile:   fmt.Sprintf("1878222222%d", i),
			Password: newPassword,
		}
		db.Save(&user)
	}

	////设置全局的logger，这个logger在我们执行每个sql语句的时候会打印每一行sql
	////sql才是最重要的，本着这个原则我尽量的给大家看到每个api背后的sql语句是什么
	//
	////定义一个表结构， 将表结构直接生成对应的表 - migrations
	//// 迁移 schema，会创建一个user用户信息表，注释不用这个的话，直接上面db.save也会默认创建表
	//_ = db.AutoMigrate(&model.User{}) //此处应该有sql语句

	// md5加密 方法演示
	fmt.Println(genMd5("xxxxx_123456"))
	//将用户的密码变一下 随机字符串+用户密码
	//暴力破解 123456 111111 000000 彩虹表 盐值
	//e10adc3949ba59abbe56e057f20f883e
	//e10adc3949ba59abbe56e057f20f883e

	// Using custom options
	//options := &password.Options{16, 100, 32, sha512.New}
	//salt, encodedPwd := password.Encode("generic password", options)
	//newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	//fmt.Println(len(newPassword))
	//fmt.Println(newPassword)
	//
	//passwordInfo := strings.Split(newPassword, "$")
	//fmt.Println(passwordInfo)
	//check := password.Verify("generic password", passwordInfo[2], passwordInfo[3], options)
	//fmt.Println(check) // true
}
