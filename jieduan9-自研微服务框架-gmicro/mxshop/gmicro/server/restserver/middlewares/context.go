package middlewares

import "github.com/gin-gonic/gin"

const (
	UsernameKey = "username"
	KeyUserID   = "userid"
	UserIP      = "ip"
)

// 为每个请求添加上下文, django
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 为每个请求都从c中获取信息字段设置到上下文中，自定义的字段很有用
		// c.Set(UsernameKey, "admin")
		//TODO 大家自己去扩展
		c.Next()
	}
}
