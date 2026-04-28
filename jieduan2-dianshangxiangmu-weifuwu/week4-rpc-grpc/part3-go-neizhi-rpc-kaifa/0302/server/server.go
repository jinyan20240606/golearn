package main

import (
	"net"
	"net/rpc"         // go 内置的rpc包，Gob（Go 私有编码）
	"net/rpc/jsonrpc" // go 内置的jsonrpc包，它是 Go 官方提供的：，基于 JSON 编码的 RPC 框架，用 JSON 代替 Gob 编码，实现跨语言 RPC 调用
)

type HelloService struct {
}

func (p *HelloService) Hello(request string, reply *string) error {
	*reply = "hello:" + request // *reply 指针，指针变量 存的是地址，*指针 才是取值 / 赋值
	return nil
}

func main() {
	// 1、实例化一个server，就是开启一个socket-tcp端口的监听
	listener, _ := net.Listen("tcp", ":1234")
	// 2、注册一个服务和handler
	_ = rpc.RegisterName("HelloService", &HelloService{})
	// 3、启动服务
	for {
		// 4、接收一个连接：当一个新的链接进来后，就会创建一个conn
		// 如果不加for循环：这里只执行【一次】Accept()，处理完这一个连接就main函数结束了
		conn, _ := listener.Accept()
		// 5、处理连接
		// go rpc.ServeConn(conn) // 这次我们就不使用这个函数，使用能替换序列化协议的方法了
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn)) // 替换为JSON序列化协议
		// 需要加go协程，否则就是串行处理一个接一个处理连接，加协程就能处理并发了
	}
}
