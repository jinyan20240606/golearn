package v1

import "gorm.io/gorm"

// 抽象工厂的实现
type DataFactory interface {
	Goods() GoodsStore
	Categorys() CategoryStore
	Brands() BrandsStore
	Banners() BannerStore
	CategoryBrands() GoodsCategoryBrandStore

	Begin() *gorm.DB
}
