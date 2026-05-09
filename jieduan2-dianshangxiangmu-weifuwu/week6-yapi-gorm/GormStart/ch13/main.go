package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"gorm.io/gorm/schema"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Language struct {
	gorm.Model
	Name    string
	AddTime sql.NullTime //每个记录创建的时候通过钩子自动加上当前时间加入到AddTime中
	// AddTime1 time.Time 这种默认的时间类型有个坑，当没有给字段添加值时，它默认会设置为零值自动变成 零值：0001-01-01 00:00:00，零值又不符合合法的时间格式，就导致数据库操作失败
	// 所以采用sql.NullTime类型，零值可以为null，数据库操作就不会失败了，实质是设置列字段的类型为时间和允许为null：add_time datetime NULL
	// 还有一种写法：*time.Time 这种指针类型的时间，零值是 nil，数据库操作也不会失败，实质也是设置列字段的类型为时间和允许为null：add_time datetime NULL
}

// 02--- BeforeCreate创建之前的钩子
func (l *Language) BeforeCreate(tx *gorm.DB) (err error) {
	l.AddTime = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	return
}

// 01----在gorm中可以通过给某一个struct添加TableName方法来自定义表名
//func (Language) TableName() string{
//	return "my_language"
//}

/*
常见场景：
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
	//NamingStrategy和Tablename不能同时配置，如果同时配置，会以TableName为准
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{ // 全局统一加前缀
			TablePrefix: "mxshop_",
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Language{})
	db.Create(&Language{
		Name: "python",
	})
}
