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

// GORM 中自定义数据类型的标准写法，专门用来让 Go 结构体切片自动存到 MySQL 的 JSON 字段
type GoodsDetailList []GoodsDetail

// 当你 Create / Save 时，GORM 自动调用 Value()，把 GoodsDetailList 切片序列化成 JSON 字符串，存入 MySQL 的 JSON 字段
func (g GoodsDetailList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// 当你 First / Find 时，GORM 自动调用 Scan()，把 MySQL 里的 JSON 字节 反序列化，自动填回 GoodsDetailList 结构体
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

// 归还库存用的归还记录对照表
type StockSellDetail struct {
	OrderSn string          `gorm:"type:varchar(200);index:idx_order_sn,unique;"` // 应该建立个唯一索引，后续会查询这个
	Status  int32           `gorm:"type:varchar(200)"`                            //1 表示已扣减 2. 表示已归还
	Detail  GoodsDetailList `gorm:"type:varchar(200)"`
}

func (StockSellDetail) TableName() string {
	return "stockselldetail"
}

// 这个表结构不采用，比较麻烦，一个订单中10个商品你就得插10条数据，因为是以商品的维度的，简单点可以直接采用以订单号为维度，上面的StockSellDetail结构
//type InventoryHistory struct {
//	user int32
//	goods int32 // 商品
//	nums int32 // 具体扣减多少
//	order int32 // 订单
//	status int32 //1. 表示库存是预扣减， 幂等性---多次重复发预扣减，都是只生效一次， 2. 表示已经支付
//}
