package model

// 常量定义
const (
	LEAVING_MESSAGES = iota + 1 // 1 普通留言
	COMPLAINT                   // 2 投诉
	INQUIRY                     // 3 询问
	POST_SALE                   // 4 售后
	WANT_TO_BUY                 // 5 求购
)

type LeavingMessages struct {
	BaseModel

	User        int32  `gorm:"type:int;index"` // 用户ID
	MessageType int32  `gorm:"type:int comment '留言类型: 1(留言),2(投诉),3(询问),4(售后),5(求购)'"`
	Subject     string `gorm:"type:varchar(100)"` // 留言主题

	Message string // 留言内容 // 字符串类型gorm中默认会转成text类型
	File    string `gorm:"type:varchar(200)"` // 附件/图片地址
}

func (LeavingMessages) TableName() string {
	return "leavingmessages"
}

// 收获地址表：一个用户可以有多个地址
type Address struct {
	BaseModel

	User         int32  `gorm:"type:int;index"`
	Province     string `gorm:"type:varchar(10)"`
	City         string `gorm:"type:varchar(10)"`
	District     string `gorm:"type:varchar(20)"`
	Address      string `gorm:"type:varchar(100)"`
	SignerName   string `gorm:"type:varchar(20)"` // 收货人姓名
	SignerMobile string `gorm:"type:varchar(11)"` // 收货人手机号
}

// 用户商品收藏表
type UserFav struct {
	BaseModel

	User  int32 `gorm:"type:int;index:idx_user_goods,unique"` // 联合唯一索引（
	Goods int32 `gorm:"type:int;index:idx_user_goods,unique"` // 联合唯一索引（
	// 设置联合唯一索引的目的：
	//	同一个用户不能重复收藏同一个商品
	//
	// 用户 ID + 商品 ID 组合唯一
}

func (UserFav) TableName() string {
	return "userfav"
}
