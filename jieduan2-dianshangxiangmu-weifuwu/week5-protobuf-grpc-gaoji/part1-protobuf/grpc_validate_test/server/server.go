package main

import (
	"context"
	"fmt"
	"net"

	"golearn/part1-protobuf/grpc_validate_test/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	proto.UnimplementedGreeterServer
}

func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{Message: "hello " + request.Name}, nil
}

// 定义一个验证接口，protoc-gen-validate 生成的代码会实现这个接口
type Validator interface {
	Validate() error
}

func main() {
	// 定义拦截器：在方法处理前进行参数验证
	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// 判断 req 是否实现了 Validator 接口
		if v, ok := req.(Validator); ok {
			if err := v.Validate(); err != nil {
				// codes 包中定义的状态码
				// 使用status 包中的 Error 函数，创建grpc中的Error对象传进状态码
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
		}
		return handler(ctx, req)
	}

	opt := grpc.UnaryInterceptor(interceptor)
	g := grpc.NewServer(opt)

	proto.RegisterGreeterServer(g, &Server{})

	lis, err := net.Listen("tcp", "0.0.0.0:8088")
	if err != nil {
		panic("监听失败" + err.Error())
	}

	fmt.Println("启动 grpc server: 0.0.0.0:8088")
	if err = g.Serve(lis); err != nil {
		panic("启动失败" + err.Error())
	}
}
