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

// User 拥有并属于多种 language，`user_languages` 是连接表，
type User3 struct {
	gorm.Model
	Languages []Language `gorm:"many2many:user_languages;"` // 定义多对多时，不需要外键字段声明，连接表会同时拥有两个模型的外键
}

type Language struct {
	gorm.Model
	Name string
}

/*
1. 我们自己定义表名是什么
2. 统一的给所有的表名加上一个前缀
*/
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

	db.AutoMigrate(&User3{})

	// 插入数据
	// languages := []Language{}
	// languages = append(languages, Language{Name:"go"})
	// languages = append(languages, Language{Name:"java"})
	// user := User3{
	// 	Languages: languages,
	// }

	// 插入数据
	// db.Create(&user)

	// 预加载的方式：查询数据，既取出用户，又取出语言
	// var user User3
	// db.Preload("Languages").First(&user)
	// for _, language := range user.Languages{
	// 	fmt.Println(language.Name)
	// }

	//如果我已经取出一个用户来了，但是这个用户我们之前没有使用preload来加载对应的Languages
	//不是说用户有language我们就一定要取出来、看业务需求
	var user User3
	db.First(&user)
	var laguages []Language
	// 单独取出来用户后，可以单独通过查询关联方法Association方法来取出关联的Languages：https://gorm.io/zh_CN/docs/associations.html#%E5%85%B3%E8%81%94%E6%A8%A1%E5%BC%8F
	_ = db.Model(&user).Association("Languages").Find(&laguages) // find参数传的变量地址，用于存储查询结果
	for _, language := range laguages {
		fmt.Println(language.Name)
	}
}
