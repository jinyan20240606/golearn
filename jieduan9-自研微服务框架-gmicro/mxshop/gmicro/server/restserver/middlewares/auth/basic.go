// 包 auth：认证相关
package auth

import (
	"encoding/base64" // 用来解码 Basic 认证的 Base64 字符串
	"mxshop/pkg/common/core"
	"strings" // 字符串切割

	"mxshop/gmicro/code"
	"mxshop/gmicro/server/restserver/middlewares"

	"mxshop/pkg/errors"

	"github.com/gin-gonic/gin"
)

// BasicStrategy
// 定义 Basic 认证策略（对应 HTTP Basic Auth）
type BasicStrategy struct {
	// compare 是外部传入的函数：用来校验 用户名+密码 是否正确
	compare func(username string, password string) bool
}

// 确保 BasicStrategy 实现了 AuthStrategy 接口（语法检查，无实际运行逻辑）
var _ middlewares.AuthStrategy = &BasicStrategy{}

// NewBasicStrategy
// 创建一个 Basic 认证策略
// 参数：compare 函数 → 校验账号密码是否正确
func NewBasicStrategy(compare func(username, password string) bool) BasicStrategy {
	return BasicStrategy{
		compare: compare,
	}
}

// AuthFunc
// 核心！这是 Gin 中间件，所有走 Basic 认证的接口都会执行这里
func (b BasicStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. 从请求头获取 Authorization
		// 格式是：Basic YWRtaW46MTIzNDU2
		// 切割成 2 段：["Basic", "YWRtaW46MTIzNDU2"]
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		// 2. 校验格式：必须是两段，第一段必须是 Basic
		if len(auth) != 2 || auth[0] != "Basic" {
			// 格式错误 → 返回 401
			core.WriteResponse(
				c,
				errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong."),
				nil,
			)
			c.Abort() // 中断请求
			return
		}

		// 3. 对第二段内容进行 Base64 解码
		// 例如：YWRtaW46MTIzNDU2 → 解码后 → admin:123456
		payload, _ := base64.StdEncoding.DecodeString(auth[1])

		// 4. 按冒号切割成 用户名 和 密码
		// admin:123456 → ["admin","123456"]
		pair := strings.SplitN(string(payload), ":", 2)

		// 5. 校验是否成功切成两段，并且调用 compare 校验账号密码是否正确
		if len(pair) != 2 || !b.compare(pair[0], pair[1]) {
			// 账号密码错误 → 返回 401
			core.WriteResponse(
				c,
				errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong."),
				nil,
			)
			c.Abort()

			return
		}

		// 6. 认证成功！把用户名存入 gin context
		c.Set(middlewares.UsernameKey, pair[0])

		// 7. 放行，继续执行后面的接口
		c.Next()
	}
}
