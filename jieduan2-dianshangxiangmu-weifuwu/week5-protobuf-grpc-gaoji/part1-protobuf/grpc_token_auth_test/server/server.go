package main

import (
	"context"
	"fmt"
	"net"

	"golearn/part1-protobuf/grpc_token_auth_test/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Server struct {
	// grpc自动实现的：你的服务端代码不一定立刻全挂掉，未实现的方法会走默认的 unimplemented 返回
	proto.UnimplementedGreeterServer
}

func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{Message: "hello " + request.Name}, nil
}

func main() {
	// 3步走
	// 定一个server端的grpc的拦截器---- 拦截器中验证元数据
	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		fmt.Println("接收到了新请求----进入拦截器")

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			// 已经开始接触到grpc的错误处理了
			return nil, status.Error(codes.Unauthenticated, "无token认证信息")
		}

		var (
			appid  string
			appkey string
		)

		if val, ok := md["appid"]; ok {
			appid = val[0]
		}

		if va1l, ok := md["appkey"]; ok {
			appkey = va1l[0]
		}

		if appid != "101010" || appkey != "i am key" {
			return nil, status.Error(codes.Unauthenticated, "appid或appkey错误")
		}

		fmt.Println("获取metadata成功")

		if nameSlice, ok := md["name"]; ok {
			fmt.Println("name:", nameSlice)
			for _, name := range nameSlice {
				fmt.Println("name:", name)
			}
		}
		for key, val := range md {
			fmt.Println(key, val)
		}

		res, err := handler(ctx, req) // 服务端中handler是放行方法
		fmt.Println("处理完成，返回响应----进入拦截器")
		return res, err
	}

	opt := grpc.UnaryInterceptor(interceptor)
	// 1. 实例化一个grpc的server
	g := grpc.NewServer(opt /*opt2*/)

	// 2. 注册服务
	proto.RegisterGreeterServer(g, &Server{})

	// 3. 启动监听服务
	lis, err := net.Listen("tcp", "0.0.0.0:8088")

	if err != nil {
		panic("监听失败" + err.Error())
	}

	// 这块不用加for循环，内部是每来一个请求都会用一个协程处理的
	if err = g.Serve(lis); err != nil {
		panic("启动失败" + err.Error())
	}
}
