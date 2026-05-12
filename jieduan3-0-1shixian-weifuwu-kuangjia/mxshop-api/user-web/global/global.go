package global

// 全局变量
// 主要存放一些全局变量，统一维护
import (
	"mxshop-api/user-web/config"
	"mxshop-api/user-web/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	Trans ut.Translator
	// 必须要指针类型，因为要变
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	// 通过初始化配置文件时，赋值到这个全局变量，供其他文件直接使用
	NacosConfig *config.NacosConfig = &config.NacosConfig{}

	UserSrvClient proto.UserClient
)
