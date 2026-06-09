package auth

import (
	// 主药使用的第3方包
	"mxshop/gmicro/server/restserver/middlewares"

	ginjwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// AuthzAudience defines the value of jwt audience field.
// jwt 固定配置（受众，随便写的）
const AuthzAudience = "mxshop.imooc.com"

// JWTStrategy defines jwt bearer authentication strategy.
// 包装了第三方的 GinJWTMiddleware
type JWTStrategy struct {
	ginjwt.GinJWTMiddleware // 👈 这里就是别人写好的 JWT 中间件
}

// 确保实现了你们项目的 AuthStrategy 接口
var _ middlewares.AuthStrategy = &JWTStrategy{}

// NewJWTStrategy create jwt bearer strategy with GinJWTMiddleware.
// 创建对象，把第三方中间件传进来
func NewJWTStrategy(gjwt ginjwt.GinJWTMiddleware) JWTStrategy {
	return JWTStrategy{gjwt}
}

// AuthFunc defines jwt bearer strategy as the gin authentication middleware.
// 作为 gin 中间件生效，直接调用第三方库的 MiddlewareFunc()
func (j JWTStrategy) AuthFunc() gin.HandlerFunc {
	return j.MiddlewareFunc()
}

// 密钥不在这个文件里！！！
// 密钥是在初始化 gin-jwt 库的时候传进去的！
// 见 jieduan9-自研微服务框架-gmicro/mxshop/app/mxshop/api/router.go的`jwtAuth := newJWTAuth(cfg.Jwt)`方法 这个启动restserver的web-server文件
