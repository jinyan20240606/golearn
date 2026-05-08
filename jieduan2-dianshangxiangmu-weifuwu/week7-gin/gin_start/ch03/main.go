package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// 路由分组
	goodsGroup := router.Group("/goods") // 统一前缀
	{                                    // 大花括号只是代码块分组用，go的基础语法没有任何实际功能影响，不开启新作用域，里面外面变量互通
		// Go 里这种单独写的 {} 完全不影响作用域。JS 里这种单独写的 {} 会影响块级作用域
		goodsGroup.GET("/list", goodsList)
		// 同方法如get的 上方/list 和下方的/:id 会冲突，因为id动态路由可以写成/list冲突,加前缀或改名
		goodsGroup.GET("/:id/:action/add", goodsDetail) //获取商品id为1的详细信息 模式
		goodsGroup.POST("/add", createGoods)
	}

	router.Run(":8083")
}

func createGoods(c *gin.Context) {

}

func goodsDetail(c *gin.Context) {
	id := c.Param("id")
	action := c.Param("action")
	c.JSON(http.StatusOK, gin.H{
		"id":     id,
		"action": action,
	})
}

func goodsList(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"name": "goodsList",
	})
}
