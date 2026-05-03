package main

import (
	"context"
	"fmt"
	"golearn/part1-protobuf/grpc_proto_test/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/golang/protobuf/ptypes/empty"
)

func main() {
	if conn, err := grpc.Dial("127.0.0.1:50053", grpc.WithInsecure()); err != nil {
		panic("连接失败" + err.Error())
	} else {
		defer conn.Close()
		// 这里调用服务端的方法

		fmt.Println(empty.Empty{}, proto.Pong{})
		// 在client中如何使用嵌套的message-Result消息结构体类型进行实例化，查看源码中它定义成了单独的名字，直接用单独的名字用就行
		// fmt.Println(proto.HelloReply_Result{})

		c := proto.NewGreeterClient(conn)
		r, err := c.SayHello(context.Background(), &proto.HelloRequest{
			Name: "jinyan",
			// 为切片中的每个元素添加结构体类型和花括号，--- 一般在类型中加*号代表作为指针类型，一般在值中加&，代表取指针这个值准备赋给某个指针类型的变量值，不是类型
			Data: []*proto.HelloRequest_Result{
				{
					Name: "jinyan",
					Age:  17,
				},
			},
			Age: 18,
			G:   proto.Gender_MALE,
			Mp: map[string]string{
				"jinyan": "18",
			},
			// 使用时间戳类型
			AddTime: timestamppb.New(time.Now()),
		})

		if err != nil {
			panic("调用失败" + err.Error())
		}

		println(r.Message)
	}
}
