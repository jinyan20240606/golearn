package admin

import (
	"mxshop/app/pkg/options"
	"mxshop/gmicro/server/restserver/middlewares"
	"mxshop/gmicro/server/restserver/middlewares/auth"

	"github.com/gin-gonic/gin"

	ginjwt "github.com/appleboy/gin-jwt/v2"
)

func newJWTAuth(opts *options.JwtOptions) middlewares.AuthStrategy {
	gjwt, _ := ginjwt.New(&ginjwt.GinJWTMiddleware{
		Realm:            opts.Realm,
		SigningAlgorithm: "HS256",
		Key:              []byte(opts.Key), // 👈 密钥就在这里！ 密钥从外部传进来，存在结构体里，可以设置过期时间
		Timeout:          opts.Timeout,
		MaxRefresh:       opts.MaxRefresh,
		LogoutResponse: func(c *gin.Context, code int) {
			c.JSON(code, nil)
		},
		IdentityHandler: claimHandlerFun,
		IdentityKey:     middlewares.KeyUserID,
		TokenLookup:     "header: Authorization:, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
	})
	return auth.NewJWTStrategy(*gjwt)
}

func claimHandlerFun(c *gin.Context) interface{} {
	claims := ginjwt.ExtractClaims(c)
	c.Set(middlewares.KeyUserID, claims[middlewares.KeyUserID])
	return claims[ginjwt.IdentityKey]
}
