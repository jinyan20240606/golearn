package initialize

import (
	"fmt"
	// consul客户端
	"github.com/hashicorp/consul/api"
	// consul 官方 grpc 客户端
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"mxshop-api/user-web/global"
	"mxshop-api/user-web/proto"
)

// 从配置里拿到 Consul 地址 → 通过 Consul 服务发现（通过UserSrvInfo.Name获取对应的用户grpc微服务），连接用户 gRPC 服务 → 生成 gRPC 客户端 → 放到全局变量供项目使用
func InitSrvConn() {
	// 从项目全局配置里，拿到 Consul 的地址、端口 等信息。
	consulInfo := global.ServerConfig.ConsulInfo
	// 通过 Consul 服务发现（通过UserSrvInfo.Name获取对应的用户grpc微服务）
	// 建立 gRPC 连接，这里不是直接连 IP，而是通过 Consul 服务发现去连用户服务。
	// 方式1: 是 gRPC 内置的【自动服务发现】**它自己内部会去 consul 拿地址、监听变化、做负载均衡
	userConn, err := grpc.Dial(
		// 服务发现的协议格式：consul://consul地址:端口/服务名?wait=14s
		// 去 Consul 里找 服务名叫 user-srv 的 gRPC 服务。
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		// 不使用 TLS/HTTPS，明文传输（开发环境用）
		grpc.WithInsecure(),
		// 开启 gRPC 负载均衡，策略是 round_robin（轮询）---- 多个服务实例时，轮流请求，实现负载均衡。
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}

	userSrvClient := proto.NewUserClient(userConn)
	// 放到全局供全项目使用
	global.UserSrvClient = userSrvClient
}

func InitSrvConn2() {
	//从注册中心获取到用户服务的信息
	cfg := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)

	userSrvHost := ""
	userSrvPort := 0
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	// 方式2: 是直接调用 Consul API 【手动服务发现】**你自己写代码去 consul 查 IP、查端口，拿到后自己连。
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ServerConfig.UserSrvInfo.Name))
	//data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.ServerConfig.UserSrvInfo.Name))
	if err != nil {
		panic(err)
	}
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}
	// 说明没有取到
	if userSrvHost == "" {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
		return
	}

	//拨号连接用户grpc服务器 跨域的问题 - 后端解决 也可以前端来解决
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】",
			"msg", err.Error(),
		)
	}
	//1. 后续的用户服务下线了 2. 改端口了 3. 改ip了 ----- 这个暂时不在这解决，由后面的grpc负载均衡来做---它们会自动感知、自动重连、自动换新节点。

	//2. 当前的好处：已经事先创立好了连接，这样后续就不用进行再次tcp的三次握手
	//  --- 提前建立好 gRPC 长连接 → 存到全局 → 后面所有接口都复用这个连接
	//3. 存在的性能问题：一个连接多个groutine共用，性能 - 连接池解决（后续由专门的grpc插件-grpc-connection-pool 解决或者负载均衡也能解决）
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}
