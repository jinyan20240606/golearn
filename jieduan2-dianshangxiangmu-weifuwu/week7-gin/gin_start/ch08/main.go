package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func MyLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Set("example", "123456")
		//让原本改执行的逻辑继续执行
		c.Next()

		end := time.Since(t)
		fmt.Printf("耗时:%V\n", end)
		status := c.Writer.Status()
		fmt.Println("状态", status)
	}
}

func Hook404() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()
		status := c.Writer.Status()
		if status == 404 {
			c.JSON(http.StatusOK, gin.H{
				"msg": "页面找不到",
			})
		}
	}
}

func main() {
	router := gin.Default()
	// 全局使用gin内置的中间件
	// 使用logger：自动打印接口请求日志
	// recovery中间件：如果接口代码崩溃（panic），服务器不会直接挂掉，捕获崩溃，返回 500 错误给前端，程序继续运行
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(Hook404())

	authrized := router.Group("/auth") // 使用自定义中间件，仅在该分组下生效
	authrized.Use(MyLogger())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.Run(":8083")
}
