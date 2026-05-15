package main

import (
	"OldPackageTest/nacos_test/config"
	"encoding/json"
	"fmt"
	"time"

	// nacos相关的包
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// 这段代码的作用：从 Nacos 远程拉取配置文件（user-web.json），然后转成 Go 结构体，打印出来。
func main() {
	// 连接 Nacos 服务器地址
	sc := []constant.ServerConfig{
		{
			IpAddr: "192.168.1.103",
			Port:   8848,
		},
	}
	// Nacos 客户端配置
	// NamespaceId：Nacos 的命名空间（用来隔离环境）
	// 日志、缓存目录：Nacos SDK 自动生成
	cc := constant.ClientConfig{
		NamespaceId:         "c1872978-d51c-4188-a497-4e0cd20b97d5", // 如果需要支持多namespace，我们可以配多个client,它们有不同的NamespaceId
		TimeoutMs:           5000,                                   // 访问 Nacos 超时时间（毫秒）
		NotLoadCacheAtStart: true,                                   // 启动时不读取本地缓存配置，true：每次启动都去 Nacos 拉最新的，不读本地缓存，false：优先读本地缓存，再去同步 nacos
		LogDir:              "tmp/nacos/log",                        // Nacos 客户端日志存放目录，日志：连接记录、报错、调试信息
		CacheDir:            "tmp/nacos/cache",                      // 本地缓存目录，Nacos 会把远程配置缓存到这里，即使 Nacos 挂了，程序还能用缓存启动
		RotateTime:          "1h",                                   // 日志切割时间，每 1 小时生成一个新日志文件
		MaxAge:              3,                                      // 日志保留天数
		LogLevel:            "debug",                                // 日志级别，debug：最详细（开发用），info：简单信息，warn/error：只打印警告、错误
	}
	// 和 Nacos 建立连接，拿到一个客户端对象，用来读取配置。
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}
	// 	DataId: "user-web.json"
	// 你在 Nacos 里创建的配置文件名
	// Group: "dev"
	// 分组（开发环境）
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "user-web.json",
		Group:  "dev"})

	if err != nil {
		panic(err)
	}
	//fmt.Println(content) //字符串 - yaml
	serverConfig := config.ServerConfig{}
	// 把 JSON 字符串转成 Go 结构体
	//想要将一个json字符串转换成struct，需要去设置这个struct的tag，否则结果永远是空
	json.Unmarshal([]byte(content), &serverConfig)
	fmt.Println(serverConfig)

	// 监听 Nacos 配置变化 → 配置改了，程序立刻自动收到最新内容，不用重启！
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: "user-web.json", // 监听哪个配置
		Group:  "dev",           // 监听哪个分组
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("配置文件变化")
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})
	time.Sleep(3000 * time.Second)

}
