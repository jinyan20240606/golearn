// jwt+cache 认证相关

// JWT 自带 kid → 代码拿 kid 去 redis 查 → 查到才算合法
// 这就是有状态、可拉黑、可踢人、可失效的 JWT！
package auth

import (
	"fmt"
	"time"

	"mxshop/gmicro/code"
	"mxshop/gmicro/server/restserver/middlewares"
	"mxshop/pkg/common/core"

	"mxshop/pkg/errors"

	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
)

// 定义错误常量 ==============================
// 错误：Token 中缺少 kid 字段
var (
	ErrMissingKID    = errors.New("Invalid token format: missing kid field in claims")
	ErrMissingSecret = errors.New("Can not obtain secret information from cache") // 错误：从缓存获取密钥失败
)

// Secret 结构体
// 存放 JWT 密钥的核心信息
type Secret struct {
	Username string // 用户名
	ID       string // 密钥ID（kid）
	Key      string // 签名密钥（加密串）
	Expires  int64  // 过期时间
}

// CacheStrategy
// 定义了【基于缓存的JWT认证策略】
// 也就是：有状态JWT（不是纯无状态JWT）
type CacheStrategy struct {
	get func(kid string) (Secret, error) // 传入 kid，从缓存/Redis 获取 Secret
}

// 确保 CacheStrategy 实现了 AuthStrategy 接口
var _ middlewares.AuthStrategy = &CacheStrategy{}

// 入口方法：NewCacheStrategy 创建缓存认证策略
// 参数：外部缓存服务动态注入，一个获取密钥的函数（从缓存/redis取）
func NewCacheStrategy(get func(kid string) (Secret, error)) CacheStrategy {
	return CacheStrategy{get}
}

// AuthFunc 核心方法
// 这是 Gin 的中间件函数，所有需要登录的接口都会走这里
func (cache CacheStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头获取 Authorization
		header := c.Request.Header.Get("Authorization")

		// 如果请求头为空，直接返回错误
		if len(header) == 0 {
			core.WriteResponse(c, errors.WithCode(code.ErrMissingHeader, "Authorization header cannot be empty."), nil)
			c.Abort() // 中断请求
			return
		}

		var rawJWT string
		// 解析请求头，格式：Bearer token
		// 把 token 部分解析到 rawJWT 变量
		fmt.Sscanf(header, "Bearer %s", &rawJWT)

		// 定义保存密钥的变量
		var secret Secret

		// 定义 jwt 载荷
		claims := &jwt.MapClaims{}

		// ==================== 核心：校验 JWT ====================
		// ParseWithClaims：解析并校验 JWT
		// 1参 ：要解析的 token
		// 2参 ：&claims{}, // 解析后放这里
		// 3参：关键回调函数,
		// 这个函数内部做的事
		// ① 自动把 JWT 拆成 3 段
		// 	Header: eyJhbGciOiJIUzI1NiIsImtpZCI6IjEyMzQ1NiJ9
		// Payload: eyJ1c2VySWQiOjF9
		// Signature: abc123456
		// ② 自动把 Header 解密成 JSON
		// {
		// 	"alg": "HS256",
		// 	"kid": "123456"  // 👈 看到没！你们的关键 kid
		// 	}
		// ③ 进入关键回调函数
		parsedT, err := jwt.ParseWithClaims(rawJWT, claims, func(token *jwt.Token) (interface{}, error) {

			// 校验签名算法是否是 HMAC（防止算法攻击）
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// 从 JWT Header 中获取 kid（关键！）
			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, ErrMissingKID // 没有 kid 直接报错
			}

			// ==================== 重点：这里就是有状态认证 ====================
			// 纯JWT：只校验签名，只要解密成功就通过
			// 缓存JWT：必须根据 kid 去缓存/redis 查密钥
			// 如果查不到 → 认证失败（拉黑/失效/踢人 都靠这一步）
			var err error
			secret, err = cache.get(kid)
			if err != nil {
				return nil, ErrMissingSecret // 缓存中没有这个密钥，认证失败
			}

			// 返回密钥，用于校验 JWT 签名
			return []byte(secret.Key), nil
		}, jwt.WithAudience(AuthzAudience))

		// 如果解析失败 或 token 无效
		if err != nil || !parsedT.Valid {
			core.WriteResponse(c, errors.WithCode(code.ErrSignatureInvalid, err.Error()), nil)
			c.Abort()

			return
		}

		// ==================== 检查密钥是否过期 ====================
		// 这里判断的是：密钥本身是否过期（不是token过期）
		if KeyExpired(secret.Expires) {
			tm := time.Unix(secret.Expires, 0).Format("2006-01-02 15:04:05")
			core.WriteResponse(c, errors.WithCode(code.ErrExpired, "expired at: %s", tm), nil)
			c.Abort()

			return
		}

		// 认证全部通过！
		// 把用户名存入 gin context，后面接口可以直接取
		c.Set(middlewares.UsernameKey, secret.Username)

		// 放行，继续执行后面的接口逻辑
		c.Next()
	}
}

// KeyExpired 判断密钥是否过期
// expires 是密钥的过期时间戳
func KeyExpired(expires int64) bool {
	if expires >= 1 {
		// 当前时间 晚于 密钥过期时间 → 已过期
		return time.Now().After(time.Unix(expires, 0))
	}
	// 没有设置过期时间，直接返回不过期
	return false
}
