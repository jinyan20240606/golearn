package router

import (
	"mxshop-api/goods-web/api/brands"
	"mxshop-api/goods-web/middlewares"

	"github.com/gin-gonic/gin"
)

// 1. 商品的api接口开发完成
// 2. 图片的坑
func InitBrandRouter(Router *gin.RouterGroup) {
	// 品牌列表的CRUD接口
	BrandRouter := Router.Group("brands").Use(middlewares.Trace())
	{
		BrandRouter.GET("", brands.BrandList)          // 品牌列表页
		BrandRouter.DELETE("/:id", brands.DeleteBrand) // 删除品牌
		BrandRouter.POST("", brands.NewBrand)          //新建品牌
		BrandRouter.PUT("/:id", brands.UpdateBrand)    //修改品牌信息
	}

	// 品牌分类的CRUD接口
	CategoryBrandRouter := Router.Group("categorybrands")
	{
		CategoryBrandRouter.GET("", brands.CategoryBrandList)          // 类别品牌列表页
		CategoryBrandRouter.DELETE("/:id", brands.DeleteCategoryBrand) // 删除类别品牌
		CategoryBrandRouter.POST("", brands.NewCategoryBrand)          //新建类别品牌
		CategoryBrandRouter.PUT("/:id", brands.UpdateCategoryBrand)    //修改类别品牌
		CategoryBrandRouter.GET("/:id", brands.GetCategoryBrandList)   //获取分类的品牌
	}
}
