package main

// 解决的问题：
// 1. 我们把服务端的业务逻辑也抽离到handler里面，保证逻辑分离，业务逻辑和非业务逻辑分离
// 2. 将服务的实例化逻辑也封装到专门的服务代理里封装到 server_proxy中

import (
	handler "golearn/part3-go-neizhi-rpc-kaifa/0304-5/handler"
	serverproxy "golearn/part3-go-neizhi-rpc-kaifa/0304-5/server_proxy"
	"net"
	"net/rpc"         // go 内置的rpc包，Gob（Go 私有编码）
	"net/rpc/jsonrpc" // go 内置的jsonrpc包，它是 Go 官方提供的：，基于 JSON 编码的 RPC 框架，用 JSON 代替 Gob 编码，实现跨语言 RPC 调用
)

func main() {
	// 1、实例化一个server，就是开启一个socket-tcp端口的监听
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	// 2、注册一个服务和handler，与客户端使用的服务名称一致, 使用 server_proxy 封装好的，不需要rpc.RegisterName
	// serviceName := handler.HelloServiceName
	// err = rpc.RegisterName(serviceName, &handler.HelloService{})
	serverproxy.RegisterHelloService(&handler.HelloService{})

	// 3、启动服务
	for {
		conn, err := listener.Accept()
		if err != nil {
			println("Accept error:", err.Error())
			continue
		}
		// 5、处理连接
		go func() {
			defer conn.Close()
			rpc.ServeCodec(jsonrpc.NewServerCodec(conn)) // 替换为JSON序列化协议
		}()
		// 需要加go协程，否则就是串行处理一个接一个处理连接，加协程就能处理并发了
	}
}
