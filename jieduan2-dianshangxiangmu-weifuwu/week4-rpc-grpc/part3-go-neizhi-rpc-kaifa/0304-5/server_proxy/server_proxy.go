package serverproxy

import (
	"golearn/part3-go-neizhi-rpc-kaifa/0304-5/handler"
	"net/rpc"
)

type HelloServicer interface {
	Hello(request string, reply *string) error
}

// 如何做到解耦：我们关心的是实际注册的函数方法名，并不是结构体本身handler.HelloService，不希望耦合严重，
// 不希望写死具体的结构体类型，要做到高度复用，是应该是外界传什么样的结构体类型都可以，就是一个透传srv的过程，不希望硬编码srv的类型
// ，这就是引入了接口带来的好处
// func RegisterHelloService(srv *handler.HelloService) error { // 这种耦合严重
func RegisterHelloService(srv HelloServicer) error { // 换成接口，解耦
	return rpc.RegisterName(handler.HelloServiceName, srv)
}
