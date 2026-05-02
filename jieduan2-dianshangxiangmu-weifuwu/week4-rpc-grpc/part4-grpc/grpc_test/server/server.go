package main

import (
	"context"
	"net"

	"golearn/part4-grpc/grpc_test/proto"

	"google.golang.org/grpc"
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
	// 1. 实例化一个grpc的server
	g := grpc.NewServer()

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
