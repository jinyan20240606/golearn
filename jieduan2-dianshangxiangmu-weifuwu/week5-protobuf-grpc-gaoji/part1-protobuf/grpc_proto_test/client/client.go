package main

import (
	"context"
	"fmt"
	"golearn/part1-protobuf/grpc_proto_test/proto"

	"google.golang.org/grpc"

	"github.com/golang/protobuf/ptypes/empty"
)

func main() {
	if conn, err := grpc.Dial("127.0.0.1:50053", grpc.WithInsecure()); err != nil {
		panic("连接失败" + err.Error())
	} else {
		defer conn.Close()
		// 这里调用服务端的方法

		fmt.Println(empty.Empty{})

		c := proto.NewGreeterClient(conn)
		r, err := c.SayHello(context.Background(), &proto.HelloRequest{Name: "jinyan"})

		if err != nil {
			panic("调用失败" + err.Error())
		}

		println(r.Message)
	}
}
