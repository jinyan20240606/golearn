package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 用来接收参数的校验结构体
type Person struct {
	// 标签规则逻辑：
	// 从 URL 路径里取名字叫 id 的参数，绑定到结构体 ID 字段
	// 绑定规则：required：必须传,uuid：必须是合法的 UUID 格式
	ID   string `uri:"id" binding:"required,uuid"`
	Name string `uri:"name" binding:"required"`
}

func main() {
	router := gin.Default()
	router.GET("/:name/:id", func(c *gin.Context) {
		var person Person
		// 校验参数
		if err := c.ShouldBindUri(&person); err != nil {
			c.Status(404)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"name": person.Name,
			"id":   person.ID,
		})
	})
	router.Run(":8083")
}
