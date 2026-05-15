package main

import (
	"OldPackageTest/grpclb_test/proto"
	"context"
	"fmt"
	"log"

	// grpc的consul负载均衡插件
	// 导入后自动注册 consul:// 协议，让gRPC Dial能识别
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important

	"google.golang.org/grpc"
)

func main() {
	// 1. 连接gRPC服务：使用Consul进行【服务发现】+【负载均衡】
	// consul://[consul地址]:[端口]/[服务名]?参数
	conn, err := grpc.Dial(
		// 协议头：consul://
		// 192.168.1.103:8500：Consul注册中心地址
		// user-srv：要调用的服务名称
		// wait=14s：等待服务发现更新的超时时间
		// tag=srv：只筛选tag为srv的服务实例
		"consul://192.168.1.103:8500/user-srv?wait=14s&tag=srv",
		// 关闭TLS安全校验（开发环境用）
		grpc.WithInsecure(),
		// 重点：设置gRPC负载均衡策略为 round_robin（轮询）
		// 10 次请求会轮流分发到不同服务实例：
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// 2. 循环发送10次请求，测试负载均衡轮询效果
	for i := 0; i < 10; i++ {
		// 创建gRPC客户端
		userSrvClient := proto.NewUserClient(conn)
		// 调用远程服务：获取用户列表
		rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
			Pn:    1,
			PSize: 2,
		})
		if err != nil {
			panic(err)
		}
		for index, data := range rsp.Data {
			fmt.Println(index, data)
		}
	}

}
