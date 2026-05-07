package main

import (
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
	Company   Company // 预加载关系字段，没在数据库里存储
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

	db.AutoMigrate(&User{}) //只写user表就行，会自动创建关联的表，新建了user表和company表，并设置了外键

	db.Create(&User{ // 此时插数据，会报错：因为缺少关联的company表数据
		Name: "bobby2",
	})

	//db.Create(&User{
	//	Name:      "bobby",
	//	Company: Company{
	//		Name:"慕课网",
	//	},
	//})

	db.Create(&User{ // 关联创建：会自动往相关联的表插入数据
		Name: "bobby2", // 此时CompanyID 可以不写，会自动根据Company插入数据
		Company: Company{ // 此时必须写已有的id值，若再写Name，则会往Company表里重复插入数据，传id的话，会关联已有的行记录数据
			ID: 1,
		},
	})

}
