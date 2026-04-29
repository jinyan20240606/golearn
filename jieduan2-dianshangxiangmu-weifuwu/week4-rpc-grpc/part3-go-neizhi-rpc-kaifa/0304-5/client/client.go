package main

// 要解决的问题
// 1. 实现代理化调用，实现zerorpc那样，直接调rpc注册的方法名，无需关注底层的非业务相关的服务链接逻辑，封装到client_proxy中

import (
	clientproxy "golearn/part3-go-neizhi-rpc-kaifa/0304-5/client_proxy"
)

func main() {
	// 1、建立一个连接,用jsonrpc.Dial发起连接
	// client, err := jsonrpc.Dial("tcp", "localhost:1234")
	// 使用代理建立链接
	client := clientproxy.NewHelloServiceClient("jsonrpc", "localhost:1234")

	var reply string
	errrr := client.Hello("你好 bobby", &reply)

	if errrr != nil {
		panic(errrr)
	}
	println(reply)

}
