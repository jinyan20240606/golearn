package models

// 这个不是数据库的模型定义，是请求参数的模型定义，主要是为了生成JWT Token时使用的载荷结构体定义
import (
	"github.com/dgrijalva/jwt-go"
)

// 自定义 JWT 载荷结构体，存放用户信息 + 标准 JWT 声明，用于生成和解析 Token。
type CustomClaims struct {
	// 加入一些自定义的信息
	ID          uint
	NickName    string
	AuthorityId uint
	// StandardClaims 结构体实现了 Claims 接口继承了  Valid() 方法
	jwt.StandardClaims
}
