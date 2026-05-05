package main

import (
	"context"
	"fmt"

	"golearn/part1-protobuf/grpc_validate_test/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("127.0.0.1:8088",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := proto.NewGreeterClient(conn)

	// 测试1: 正常请求
	fmt.Println("=== 测试1: 正常请求 (name='张三', age=25) ===")
	resp, err := c.SayHello(context.Background(), &proto.HelloRequest{
		Name: "张三",
		Age:  25,
	})
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Println("响应:", resp.Message)
	}

	// 测试2: name为空（违反 min_len: 2）
	fmt.Println("\n=== 测试2: name为空字符串 (违反 min_len: 2) ===")
	resp, err = c.SayHello(context.Background(), &proto.HelloRequest{
		Name: "",
		Age:  25,
	})
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Println("响应:", resp.Message)
	}

	// 测试3: age为负数（违反 gt: 0）
	fmt.Println("\n=== 测试3: age为负数 (违反 gt: 0) ===")
	resp, err = c.SayHello(context.Background(), &proto.HelloRequest{
		Name: "李四",
		Age:  -1,
	})
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Println("响应:", resp.Message)
	}

	// 测试4: age超过150（违反 lte: 150）
	fmt.Println("\n=== 测试4: age=200 (违反 lte: 150) ===")
	resp, err = c.SayHello(context.Background(), &proto.HelloRequest{
		Name: "王五",
		Age:  200,
	})
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Println("响应:", resp.Message)
	}
}
