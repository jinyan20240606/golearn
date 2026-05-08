package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"OldPackageTest/gin_start/ch06/proto"
)

func main() {
	router := gin.Default()

	router.GET("/moreJSON", moreJSON)
	router.GET("/someProtoBuf", returnProto)

	router.Run(":8083")
}

func returnProto(c *gin.Context) {
	course := []string{"python", "go", "微服务"}
	user := &proto.Teacher{
		Name:   "bobby",
		Course: course,
	}
	// 02. 返回protobuf，，，，---- 这接口返给前端的是原始的protobuf数据，前端是无法直接解析的，除非前端也用protobuf解析，否则就是一串乱码，适合内部服务之间通信
	c.ProtoBuf(http.StatusOK, user)
}

func moreJSON(c *gin.Context) {
	var msg struct { //
		Name    string `json:"user"` // 自动转移成user
		Message string
		Number  int
	}
	msg.Name = "bobby"
	msg.Message = "这是一个测试json"
	msg.Number = 20

	// 01. 返回json，json也是可以配置的，默认是map[string]interface{}，也可以直接传结构体，gin会自动转换成json格式

	c.JSON(http.StatusOK, msg)
}
