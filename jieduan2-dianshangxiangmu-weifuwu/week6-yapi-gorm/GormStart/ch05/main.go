package main

import (
	"database/sql"
	"errors"
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
	// 查询单个对象：first就是指定表名，默认按主键升序查询
	//1- 通过first查询单个数据, 获取第一条记录（主键升序），take方法没有指定排序字段，last方法默认是主键倒序
	var user User
	// db.First(&user)

	//2- 通过主键查询
	//我们下面运行必须把上面First注释，因为不能重复给user赋值，我们不能给user赋值
	result := db.First(&user, []int{1, 2, 3}) // 通过主键查询,id in 1,2,3的数据的第一条
	// 判断是否是同一种error类型
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("未找到")
	}
	fmt.Println(user.ID)

	// 查询全部对象
	//3- 检索全部对象
	var users []User
	result := db.Find(&users)                 // 检索全部对象
	fmt.Println("总共记录:", result.RowsAffected) // 影响行数
	for _, user := range users {
		fmt.Println(user.ID)
	}

}
