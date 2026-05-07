package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// `User` 属于 `Company`，`CompanyID` 是外键
type User struct {
	gorm.Model
	Name      string
	CompanyID int     //数据库中存储的字段company_id
	Company   Company // 多对一写法：预加载关系字段，没在数据库里存储
}

type Company struct {
	ID   int
	Name string
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

	var user User
	// 使用默认方法：db.First(&user) 默认不会查询关联的表

	// 使用下面的方法就可以关联查询了
	// 01- 使用preload方法
	//db.Preload("Company").First(&user)
	// 02- 使用joins方法
	db.Joins("Company").First(&user)
	// 03- 取值时看看能不能取出关联的Company的信息
	fmt.Println(user.Name, user.Company.ID)
}
