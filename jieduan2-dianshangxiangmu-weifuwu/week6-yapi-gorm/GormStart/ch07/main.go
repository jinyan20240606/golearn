package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	ID           uint   // GORM 默认会把字段名为 ID 的字段，自动当作主键！你不需要写任何标签（如 gorm:"primaryKey"），只要字段名叫 ID，它就自动是主键
	MyName       string `gorm:"column:name"`
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

	//查询方式条件有三种 1. string 2. struct 3. map
	var user User
	db.First(&user)

	//1. 通过save方法更新：一条数据的多个列
	user.MyName = "bobby test"
	user.Age = 100
	user.ID = 0    // 主键为零值 → 新增一条数据行（INSERT 行），主键有值 → 更新已有数据行（UPDATE 行）
	db.Save(&user) //save方法是一个集create(新增数据行)和update于一体的操作

	//2. 通过update方法更新：更新单个列
	db.Model(&User{}).Where("name = ?", "bobby").Update("name", "hello") // update1参是列，2参是值
	//3. 通过updates方法更新：更新多个列
	db.Model(&User{}).Where("name = ?", "bobby").Updates(User{MyName: "hello", Age: 100})
	//4. 通过updates和map更新：更新多个列
	db.Model(&User{}).Where("id = ?", 1).Updates(map[string]any{
		"my_name": "bobby test",
		"age":     100,
	})

	// 区别
	// Save()：会覆盖所有字段（包括零值）
	// Updates()：只覆盖你赋值的字段（更安全）：
	// // // // // gorm中所有方法如Updates()，save等，中传入一个数据库表不存在的字段 **，只会直接报错，绝对不会自动新增列
	// // // // // Save = 新增或修改「数据行」，Save ≠ 新增数据表「列」

	// 更新或忽略选定的字段
	db.Model(&user).Select("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
	// UPDATE users SET name='hello' WHERE id=111;
	db.Model(&user).Omit("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
	// UPDATE users SET age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;

	// 选择 Struct 的字段（会选中零值的字段）
	db.Model(&user).Select("Name", "Age").Updates(User{Name: "new_name", Age: 0})
	// UPDATE users SET name='new_name', age=0 WHERE id=111;

	// 选择所有字段（选择包括零值字段的所有字段）
	db.Model(&user).Select("*").Updates(User{Name: "jinzhu", Role: "admin", Age: 0})

	// 选择除 Role 外的所有字段（包括零值字段的所有字段）
	db.Model(&user).Select("*").Omit("Role").Updates(User{Name: "jinzhu", Role: "admin", Age: 0})

}
