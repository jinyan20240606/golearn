package main

import (
	"context"
	"golearn/part4-grpc/grpc_test/proto"

	"google.golang.org/grpc"
)

func main() {
	if conn, err := grpc.Dial("127.0.0.1:8088", grpc.WithInsecure()); err != nil {
		panic("连接失败" + err.Error())
	} else {
		defer conn.Close()
		// 这里调用服务端的方法
		c := proto.NewGreeterClient(conn)

		r, err := c.SayHello(context.Background(), &proto.HelloRequest{Name: "jinyan"})

		if err != nil {
			panic("调用失败" + err.Error())
		}

		println(r.Message)
	}
}
