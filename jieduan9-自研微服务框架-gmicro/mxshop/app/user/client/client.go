package main

import (
	"context"
	"fmt"
	v1 "mxshop/api/user/v1"
	"mxshop/gmicro/registry/consul"
	rpc "mxshop/gmicro/server/rpcserver"
	_ "mxshop/gmicro/server/rpcserver/resolver/direct"
	"mxshop/gmicro/server/rpcserver/selector"
	"mxshop/gmicro/server/rpcserver/selector/random"
	"time"

	"github.com/hashicorp/consul/api"
)

func main() {
	// ====================== 直连模式：一行代码直接连接 ============ START ==========
	// 不用 consul
	//不需要传注册中心：rpc.WithDiscovery(r),
	//不用初始化客户端
	//直接写死 IP:PORT

	// conn, err := rpc.DialInsecure(context.Background(),
	// 	// 关键：直连地址，多个用逗号分隔
	// 	rpc.WithEndpoint("direct:///127.0.0.1:8081,127.0.0.1:8082"),
	// )
	// ========================================================== END ============
	//设置全局的负载均衡策略
	// 把你们自研的【随机负载均衡算法】注册到 gRPC 里，替换掉官方 round_robin具体算法！
	selector.SetGlobalSelector(random.NewBuilder())
	// 把你们自研的 selector功能 注册到 gRPC 内部的负载均衡列表里
	rpc.InitBuilder()

	conf := api.DefaultConfig()
	conf.Address = "127.0.0.1:8500"
	conf.Scheme = "http"
	cli, err := api.NewClient(conf)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(true))

	conn, err := rpc.DialInsecure(context.Background(),
		// 不用官方 round_robin。用我们自己写的 selector 负载均衡！
		rpc.WithBalancerName("selector"),
		// 必须传注册中心
		rpc.WithDiscovery(r),
		rpc.WithClientTimeout(time.Second*5000),
		// 传服务发现的地址
		rpc.WithEndpoint("discovery:///mxshop-user-srv"),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	uc := v1.NewUserClient(conn)

	for {
		_, err := uc.GetUserList(context.Background(), &v1.PageInfo{})
		if err != nil {
			panic(err)
		}

		fmt.Println("success")
		time.Sleep(time.Millisecond * 2)
	}

}
