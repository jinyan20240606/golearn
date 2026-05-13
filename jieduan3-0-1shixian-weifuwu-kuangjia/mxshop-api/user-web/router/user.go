package router

// 路由：

import (
	"mxshop-api/user-web/api"
	"mxshop-api/user-web/middlewares"

	"github.com/gin-gonic/gin"
)

// 调用方传递一个全局的gin-group实例，便于复用
func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	{
		// 加上鉴权中间件
		UserRouter.GET("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("pwd_login", api.PassWordLogin)
		UserRouter.POST("register", api.Register)

		UserRouter.GET("detail", middlewares.JWTAuth(), api.GetUserDetail)
		UserRouter.PATCH("update", middlewares.JWTAuth(), api.UpdateUser)
	}
	//服务注册和发现
}
