package router

import (
	"mxshop-api/goods-web/middlewares"

	"github.com/gin-gonic/gin"

	"mxshop-api/goods-web/api/goods"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	// 给 /goods/* 所有接口统一加链路追踪
	//每个请求都会生成一个根 Span
	//链路结构清晰
	GoodsRouter := Router.Group("goods").Use(middlewares.Trace())
	// 一旦设计修改数据库都需要加上权限校验，查询类的可以不用鉴权
	{
		GoodsRouter.GET("", goods.List)                                                            //商品列表
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New)          //改接口需要管理员权限
		GoodsRouter.GET("/:id", goods.Detail)                                                      //获取商品的详情
		GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Delete) //删除商品
		GoodsRouter.GET("/:id/stocks", goods.Stocks)                                               //获取商品的库存

		GoodsRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)         // 全量更新用PUT方法
		GoodsRouter.PATCH("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.UpdateStatus) // 部分更新用PATCH方法
	}
}
