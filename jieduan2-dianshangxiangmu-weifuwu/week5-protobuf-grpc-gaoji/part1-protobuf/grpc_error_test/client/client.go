package main

import (
	"context"
	"fmt"
	"golearn/part1-protobuf/grpc_error_test/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func main() {

	interceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...) // 客户端中 放行方法是invoker

		fmt.Printf("请求耗时: %v\n", time.Since(start)) // since:从start到现在花了多长时间，%v是占位符，自动展示合适的变量值

		return err

	}

	opt := grpc.WithUnaryInterceptor(interceptor)

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, opt)

	if conn, err := grpc.Dial("127.0.0.1:8088", opts...); err != nil {
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

		// 设置客户端超时3秒
		ctx, _ = context.WithTimeout(ctx, time.Second*3)

		r, err := c.SayHello(ctx, &proto.HelloRequest{Name: "jinyan"})

		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				// ok报错，说明不是一个status Error的错误
				panic("status error 解析失败")
			}
			fmt.Println(st.Message())
			fmt.Println(st.Code())

			panic("调用失败" + err.Error())
		}

		println(r.Message)
	}
}
