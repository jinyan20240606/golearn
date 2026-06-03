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
	// 这是自定义的错误体系在grpc中兼容使用
	return nil, errors.ToGrpcError(e)
	// 下面注释是gprc的自带error体系
	return nil, status.Error(codes.NotFound, "user not found")
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
