package main

import (
	"context"
	"golearn/part1-protobuf/metadata_test/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	if conn, err := grpc.Dial("127.0.0.1:8088", grpc.WithInsecure()); err != nil {
		panic("连接失败" + err.Error())
	} else {
		defer conn.Close()
		// 这里调用服务端的方法
		c := proto.NewGreeterClient(conn)

		// 添加metadata数据发送

		md := metadata.New(map[string]string{ // Go 直接字面量创建一个 map，map[键类型]值类型{健：值}
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
			"version":   "1.0.0",
			"name":      "jinyan",
		})

		ctx := metadata.NewOutgoingContext(context.Background(), md)

		r, err := c.SayHello(ctx, &proto.HelloRequest{Name: "jinyan"})

		if err != nil {
			panic("调用失败" + err.Error())
		}

		println(r.Message)
	}
}
