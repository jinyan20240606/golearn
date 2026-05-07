package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func pong(c *gin.Context) {
	// type H map[string]interface{}  // gin.H是gin框架定义的一个map类型，方便我们在返回json数据的时候使用，key是字符串，value是any类型
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
	// 直接写下面的类型也不会报错
	// c.JSON(http.StatusOK, map[string]string{
	// 	"message": "pong",
	// })
}
func main() {
	//实例化一个gin的server对象
	r := gin.Default()
	r.GET("/ping", pong)
	r.Run(":8083") // listen and serve on 0.0.0.0:8080
}
