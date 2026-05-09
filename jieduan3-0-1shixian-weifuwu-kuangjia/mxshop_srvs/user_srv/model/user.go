package model

import (
	"time"

	"gorm.io/gorm"
)

// 每张表都有的 id、创建时间、更新时间、软删除
// 不用每个结构体都写一遍
// 后续所有的表都继承这个基础表模型
type BaseModel struct {
	ID        int32     `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

// 用户信息表，结构
type User struct {
	BaseModel
	// index:idx_mobile：创建索引 idx_mobile，查询更快，unique：唯一，不能重复，type:varchar(11)：数据库类型 varchar，长度 11，not null：不能为空
	Mobile string `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	// varchar (100)：存加密密码（bcrypt 很长），不能为空
	Password string `gorm:"type:varchar(100);not null"`
	NickName string `gorm:"type:varchar(20)"`
	// 时间采用指针类型，目的：可以为 nil → 数据库存 NULL，如果不用指针，没传值会存 0001-01-01 垃圾时间
	Birthday *time.Time `gorm:"type:datetime"`
	// column:gender：数据库字段名指定叫 gender，default:male：默认值 男，comment：字段注释
	Gender string `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女, male表示男'"`
	// default:1 → 默认普通用户，1 = 普通用户，2 = 管理员
	Role int `gorm:"column:role;default:1;type:int comment '1表示普通用户, 2表示管理员'"`
}
