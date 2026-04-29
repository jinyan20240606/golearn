package clientproxy

import (
	"golearn/part3-go-neizhi-rpc-kaifa/0304-5/handler"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type HelloServiceStub struct {
	*rpc.Client //匿名结构体字段，可以将方法自动提升，后面可以直接调用Call
}

// 在go语言中，没有类对象，就意味着没有类的初始化功能

func NewHelloServiceClient(protol, address string) *HelloServiceStub {
	client, err := jsonrpc.Dial("tcp", "localhost:1234")

	if err != nil {
		panic(err)
	}
	return &HelloServiceStub{client}
}

func (c *HelloServiceStub) Hello(request string, reply *string) error {
	err := c.Call(handler.HelloServiceName+".Hello", request, reply)
	return err
}
