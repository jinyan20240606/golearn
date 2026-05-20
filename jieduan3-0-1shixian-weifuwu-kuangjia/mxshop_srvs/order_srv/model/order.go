package model

import "time"

// 记录用户加购的商品，用户下单前所有商品都放这里。
type ShoppingCart struct {
	BaseModel
	User int32 `gorm:"type:int;index"` //查 “我的购物车” 时，直接 where user=? 快速查。
	// 该表一行记录对应一个商品
	Goods int32 `gorm:"type:int;index"` //商品 ID，加索引
	// 加索引的原则：我们需要查询时候才会加， 1. 会影响插入性能 2. 会占用磁盘
	// 不查询就不加索引
	Nums    int32 `gorm:"type:int"` // 购买数量
	Checked bool  //是否选中（用于下单时勾选哪些商品）
}

// 指定表名
func (ShoppingCart) TableName() string {
	return "shoppingcart"
}

// 订单主表：张订单的：用户、金额、状态、支付、地址 都在这里。
type OrderInfo struct {
	BaseModel

	User int32 `gorm:"type:int;index"`
	// 平台内部唯一订单编号必须加索引用于客服查询、支付回调、订单操作
	OrderSn string `gorm:"type:varchar(30);index"` //订单号，我们平台自己生成的订单号
	PayType string `gorm:"type:varchar(20) comment 'alipay(支付宝)， wechat(微信)'"`

	//status大家可以考虑使用iota来做，待支付、支付成功、超时关闭、交易结束，订单所有操作都围绕状态流转
	Status     string `gorm:"type:varchar(20)  comment 'PAYING(待支付), TRADE_SUCCESS(成功)， TRADE_CLOSED(超时关闭), WAIT_BUYER_PAY(交易创建), TRADE_FINISHED(交易结束)'"`
	TradeNo    string `gorm:"type:varchar(100) comment '交易号'"` //交易号就是支付宝的订单号 查账
	OrderMount float32
	PayTime    *time.Time `gorm:"type:datetime"`

	Address      string `gorm:"type:varchar(100)"`
	SignerName   string `gorm:"type:varchar(20)"`
	SingerMobile string `gorm:"type:varchar(11)"`
	Post         string `gorm:"type:varchar(20)"`
}

func (OrderInfo) TableName() string {
	return "orderinfo"
}

// OrderGoods 订单商品表 --- 订单商品需要单独拿一张表来村
type OrderGoods struct {
	BaseModel

	Order int32 `gorm:"type:int;index"`
	Goods int32 `gorm:"type:int;index"`

	//为什么把商品的信息保存下来了？？？ ， 故意要字段冗余， 因为高并发系统中我们一般都不会遵循三范式  1、做镜像 记录，后面产生纠纷，能拿到用户当时下单的商品快照，2、也是为了提高性能，减少跨服务之间的查询，有些时候故意要字段冗余
	GoodsName  string `gorm:"type:varchar(100);index"`
	GoodsImage string `gorm:"type:varchar(200)"`
	GoodsPrice float32
	Nums       int32 `gorm:"type:int"`
}

func (OrderGoods) TableName() string {
	return "ordergoods"
}
