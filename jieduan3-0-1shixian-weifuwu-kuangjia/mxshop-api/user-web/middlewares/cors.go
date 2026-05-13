package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		// 处理options请求：设置响应头告诉浏览器这些请求头下是允许浏览器非简单请求访问的

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			// 直接返回204状态码，告诉浏览器这个options请求是成功的，后续的真正请求就可以继续了
			// 204：无内容、请求成功，但没有返回数据。
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
}
