package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"mxshop/pkg/errors"
	"time"

	_ "mxshop/app/pkg/code"
	pb "mxshop/cmd/order/rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50052", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		// --------------
		// gRPC 串联解析自定义错误体系的写法
		// --------------
		// 从 gRPC error 解析出 withCode error
		s := errors.FromGrpcError(err)
		coder := errors.ParseCoder(s)
		fmt.Println(coder.Code())

		// --------------
		// gRPC 原生解析错误的体系写法
		// --------------
		st, ok := status.FromError(err)
		if !ok {
			log.Fatalf("解析grpc错误失败: %v", err)
		}

		// 👇 这就是 gRPC 原生能拿到的所有信息
		fmt.Println("gRPC的内置错误码:", st.Code()) // codes.NotFound = 5
		fmt.Println("错误消息:", st.Message())    // user not found
		fmt.Println("详细信息:", st.Details())
	}
	log.Printf("Greeting: %s", r.GetMessage())

}
