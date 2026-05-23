package model

// 商品搜索结构体
type EsGoods struct { // 基本上与mysql中的商品结构体一一对应，不需要全部一一对应，只保留搜索查询相关的字段
	ID         int32 `json:"id"`
	CategoryID int32 `json:"category_id"`
	BrandsID   int32 `json:"brands_id"`
	OnSale     bool  `json:"on_sale"`
	ShipFree   bool  `json:"ship_free"`
	IsNew      bool  `json:"is_new"`
	IsHot      bool  `json:"is_hot"`

	Name        string  `json:"name"`
	ClickNum    int32   `json:"click_num"`
	SoldNum     int32   `json:"sold_num"`
	FavNum      int32   `json:"fav_num"`
	MarketPrice float32 `json:"market_price"`
	GoodsBrief  string  `json:"goods_brief"`
	ShopPrice   float32 `json:"shop_price"`
}

func (EsGoods) GetIndexName() string {
	return "goods"
}

// es的商品mapping结构体，不用es自动的
func (EsGoods) GetMapping() string {
	goodsMapping := `
	{
		"mappings" : {
			"properties" : {
				"brands_id" : {
					"type" : "integer"
				},
				"category_id" : {
					"type" : "integer"
				},
				"click_num" : {
					"type" : "integer"
				},
				"fav_num" : {
					"type" : "integer"
				},
				"id" : {
					"type" : "integer"
				},
				"is_hot" : {
					"type" : "boolean"
				},
				"is_new" : {
					"type" : "boolean"
				},
				"market_price" : {
					"type" : "float"
				},
				"name" : {
					"type" : "text",
					"analyzer":"ik_max_word"
				},
				"goods_brief" : {
					"type" : "text",
					"analyzer":"ik_max_word"
				},
				"on_sale" : {
					"type" : "boolean"
				},
				"ship_free" : {
					"type" : "boolean"
				},
				"shop_price" : {
					"type" : "float"
				},
				"sold_num" : {
					"type" : "long"
				}
			}
		}
	}`
	return goodsMapping
}
