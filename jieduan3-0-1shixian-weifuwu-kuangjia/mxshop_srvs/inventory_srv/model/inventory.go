package model

import (
	"database/sql/driver"
	"encoding/json"
)

// 仓库表
// type Stock struct {
// 	BaseModel
// 	Name    string
// 	Address string
// }

type GoodsDetail struct {
	Goods int32
	Num   int32
}
type GoodsDetailList []GoodsDetail

func (g GoodsDetailList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GoodsDetailList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

type Inventory struct {
	BaseModel
	Goods  int32 `gorm:"type:int;index"` // 指向商品ID，现在是跨微服务了，不可能直接设外键了关联到商品表了，index加索引，方便查询
	Stocks int32 `gorm:"type:int"`       // 库存数量
	// Stock   Stock // 还得与仓库做多对多关联一个仓库有多个商品的库存，一个商品的库存也可能有多个仓库的库存，本节简单点，先不引入仓库表的复杂度了
	Version int32 `gorm:"type:int"` //分布式锁的乐观锁 ，这个字段很重要，要做分布式锁
}

type InventoryNew struct {
	BaseModel
	Goods   int32 `gorm:"type:int;index"`
	Stocks  int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"` //分布式锁的乐观锁
	Freeze  int32 `gorm:"type:int"` //冻结库存 --- 用于分布式事务TCC的实现方案
}

// 仓库服务：这个表的作用：生成出库单、记录出库单，仓库管理员要拿这个单子去办理出库出货
type Delivery struct {
	Goods   int32  `gorm:"type:int;index"`    // 商品ID
	Nums    int32  `gorm:"type:int"`          // 出库数量
	OrderSn string `gorm:"type:varchar(200)"` // 订单号，哪个订单产生的
	Status  string `gorm:"type:varchar(200)"` //1. 表示等待支付 2. 表示支付成功 3. 失败------ T阶段改成1，确认时改成2，回滚时改成3
}

type StockSellDetail struct {
	OrderSn string          `gorm:"type:varchar(200);index:idx_order_sn,unique;"`
	Status  int32           `gorm:"type:varchar(200)"` //1 表示已扣减 2. 表示已归还
	Detail  GoodsDetailList `gorm:"type:varchar(200)"`
}

func (StockSellDetail) TableName() string {
	return "stockselldetail"
}

//type InventoryHistory struct {
//	user int32
//	goods int32 // 商品
//	nums int32 // 具体扣减多少
//	order int32 // 订单
//	status int32 //1. 表示库存是预扣减， 幂等性---多次重复发预扣减，都是只生效一次， 2. 表示已经支付
//}
