package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"mxshop/app/pkg/code"
	"mxshop/pkg/errors"
	"net"

	pb "mxshop/cmd/order/rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	port = flag.Int("port", 50052, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	e := errors.WithCode(code.ErrUserNotFound, "user not found")
	// 方式1：我们的方案 - 手动调用 ToGrpcError 转换（因为 withCode 没有实现 GRPCStatus() 接口）
	return nil, errors.ToGrpcError(e)
	// 方式2：gRPC 原生 error 体系（只有 grpcCode + 纯文本，没有业务码）
	return nil, status.Error(codes.NotFound, "user not found")
	// 方式3（如果 withCode 实现了 GRPCStatus() 接口，像 Kratos 那样）：
	// 直接返回即可，gRPC 框架会自动调用 e.GRPCStatus() 转换，不需要手动 ToGrpcError
	// return nil, e
	//return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
