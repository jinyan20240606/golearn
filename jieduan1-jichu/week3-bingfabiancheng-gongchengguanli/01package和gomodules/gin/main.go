package main

import (
	"net/http"
	// 写个第三方包时，默认会飘红，点击它的自动检测依赖安装，然后就可以使用了，会自动维护在go.mod文件中
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080

}
