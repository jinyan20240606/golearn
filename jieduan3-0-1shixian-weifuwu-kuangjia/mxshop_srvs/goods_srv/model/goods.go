package model

import (
	"context"
	"strconv"

	"gorm.io/gorm"

	"mxshop_srvs/goods_srv/global"
)

//类型， 这个字段是否能为null， 这个字段应该设置为可以为null还是设置为空， 0
//实际开发过程中 尽量设置为不为null
//https://zhuanlan.zhihu.com/p/73997266
//这些类型我们使用int32还是int

// 我们这样用嵌套的写法建3级分类的表结构
type Category struct {
	BaseModel
	// Name不能为null，且也不能添加默认值，为null会有一些潜在问题
	Name string `gorm:"type:varchar(20);not null" json:"name"`
	// GORM 自关联的语法
	ParentCategoryID int32 `json:"parent"` // 定义个外键---真正在数据库中存储的外键名称
	// ParentCategory它的唯一作用：告诉 GORM 这个结构体是 “自关联”，让 GORM 知道用哪个字段去关联查询，它本身不生成数据库字段！，GORM 默认命名规范，外键字段名必须是：关联字段名 + ID（ParentCategory + ID = ParentCategoryID）
	ParentCategory *Category `json:"-"` // 外键，自己指向自己时，必须得写成指针，否则会无限递归嵌套，字段只给gorm内部逻辑代码用不存库，不给前端返回，不参与 JSON 序列化
	// 这个字段也是gorm的逻辑字段--后续都是通过gorm预加载赋的相关联值，不建库字段：foreignKey:ParentCategoryID → 子分类用这个字段作为外键关联父的ID即references:ID → 关联的是当前表的主键 ID
	SubCategory []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
	// 我们这里一般都定义成int32，一般不用int，因为后期使用proto和数据库时，都支持int32以上，不支持int，需要频繁转换，所以为了减少转换，直接设置int32
	Level int32 `gorm:"type:int;not null;default:1" json:"level"`
	// 表示是否显示在首页 Tab 栏
	IsTab bool `gorm:"default:false;not null" json:"is_tab"`
}

type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(20);not null"`
	Logo string `gorm:"type:varchar(200);default:'';not null"`
}

// 多对多关系：商品分类和品牌对应关系表
type GoodsCategoryBrand struct {
	BaseModel
	// 创建联合唯一索引
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Category   Category

	BrandsID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Brands   Brands
}

// 这个表就不会自动转化下划线
func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);not null"`
	// 轮博图的先后顺序
	Index int32 `gorm:"type:int;default:1;not null"`
}

// 商品表
type Goods struct {
	// 基础字段：ID、创建时间、更新时间、软删除
	BaseModel
	// 分类关联：商品所属分类
	CategoryID int32    `gorm:"type:int;not null"` // 分类ID
	Category   Category // 分类对象（仅GORM关联使用）
	// 品牌关联：商品所属品牌
	BrandsID int32  `gorm:"type:int;not null"` // 品牌ID
	Brands   Brands // 品牌对象（仅GORM关联使用）

	// 商品状态控制
	OnSale   bool `gorm:"default:false;not null"` // 是否上架（true=上架，false=下架）
	ShipFree bool `gorm:"default:false;not null"` // 是否包邮
	IsNew    bool `gorm:"default:false;not null"` // 是否新品
	IsHot    bool `gorm:"default:false;not null"` // 是否热销（热门商品）

	// 商品基础信息
	Name        string  `gorm:"type:varchar(50);not null"`   // 商品名称
	GoodsSn     string  `gorm:"type:varchar(50);not null"`   // 商品货号（唯一编号）
	ClickNum    int32   `gorm:"type:int;default:0;not null"` // 点击数（浏览量）
	SoldNum     int32   `gorm:"type:int;default:0;not null"` // 销量
	FavNum      int32   `gorm:"type:int;default:0;not null"` // 收藏量
	MarketPrice float32 `gorm:"not null"`                    // 市场标价（原价）--- float32类不用指定数据库中的type，会自动对应成数据库的float32，不像int32那样
	ShopPrice   float32 `gorm:"not null"`                    // 销售价（实际售价）
	GoodsBrief  string  `gorm:"type:varchar(100);not null"`  // 商品简介/短描述
	// 这里用自定义类型GormList，或者用官方提供的 datatypes.JSON，支持数组对象
	Images          GormList `gorm:"type:varchar(1000);not null"` // 商品轮播图列表（JSON数组）
	DescImages      GormList `gorm:"type:varchar(1000);not null"` // 商品详情图列表
	GoodsFrontImage string   `gorm:"type:varchar(200);not null"`  // 商品封面图（主图）
}

// 创建钩子中同步es
func (g *Goods) AfterCreate(tx *gorm.DB) (err error) {
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandsID:    g.BrandsID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}

	_, err = global.EsClient.Index().Index(esModel.GetIndexName()).BodyJson(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// 更新钩子中同步es
func (g *Goods) AfterUpdate(tx *gorm.DB) (err error) {
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandsID:    g.BrandsID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}

	_, err = global.EsClient.Update().Index(esModel.GetIndexName()).
		Doc(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// 删除钩子中同步es
func (g *Goods) AfterDelete(tx *gorm.DB) (err error) {
	_, err = global.EsClient.Delete().Index(EsGoods{}.GetIndexName()).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
