package main

import (
	"context"
	"golearn/part1-protobuf/metadata_test/proto"

	"google.golang.org/grpc"
)

type customCredential struct {
}

// 重点实现这个方法
func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"appid": "101010", "appkey": "i am key"}, nil
}

func (c customCredential) RequireTransportSecurity() bool {
	return false
}

func main() {

	// 这是实现auth的原生拦截器实现 01
	// interceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// 	start := time.Now()

	// 	// 添加metadata数据发送

	// 	md := metadata.New(map[string]string{ // Go 直接字面量创建一个 map，map[键类型]值类型{健：值}
	// 		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	// 		"appid":     "101010",
	// 		"appkey":    "i am key",
	// 	})

	// 	ctx = metadata.NewOutgoingContext(context.Background(), md)

	// 	err := invoker(ctx, method, req, reply, cc, opts...) // 客户端中 放行方法是invoker

	// 	fmt.Printf("请求耗时: %v\n", time.Since(start)) // since:从start到现在花了多长时间，%v是占位符，自动展示合适的变量值

	// 	return err

	// }

	// go中有一个专门的拦截器去实现，可以把代码写的更加简单 02

	// opt := grpc.WithUnaryInterceptor(interceptor)
	opt := grpc.WithPerRPCCredentials(customCredential{})

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, opt)

	if conn, err := grpc.Dial("127.0.0.1:8088", opts...); err != nil {
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
