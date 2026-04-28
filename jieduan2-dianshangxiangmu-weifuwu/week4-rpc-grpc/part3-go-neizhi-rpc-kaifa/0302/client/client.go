package main

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	// 1、建立一个连接,用Dial发起连接

	conn, err := net.Dial("tcp", "localhost:1234")

	if err != nil {
		panic(err)
	}

	var reply string                                               // 虽然只声明没赋值，有默认值，有内存占用的，并不是一个单纯的地址变量
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn)) // 2、创建一个客户端
	// 底层发送的格式就是json字符串
	// 同步阻塞风格：client.Call() 会阻塞主协程，✅ 服务端不返回，后面代码永远不执行，其他协程会不影响
	errrr := client.Call("HelloService.Hello", "你好 bobby", &reply)

	if errrr != nil {
		panic(errrr)
	}
	println(reply)

}
