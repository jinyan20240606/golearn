package main

import (
	"io"
	"net/http"
	"net/rpc" // go 内置的rpc包，Gob（Go 私有编码）
	"net/rpc/jsonrpc"
	// go 内置的jsonrpc包，它是 Go 官方提供的：，基于 JSON 编码的 RPC 框架，用 JSON 代替 Gob 编码，实现跨语言 RPC 调用
)

type HelloService struct {
}

func (p *HelloService) Hello(request string, reply *string) error {
	*reply = "hello:" + request // *reply 指针，指针变量 存的是地址，*指针 才是取值 / 赋值
	return nil
}

func main() {
	_ = rpc.RegisterName("HelloService", &HelloService{})
	// 1、实例化一个server，就是开启一个socket-tcp端口的监听
	// 我们不使用net了，因为它可以监听tcp，我们使用http专用的包
	// listener, _ := net.Listen("tcp", ":1234")
	http.HandleFunc("/jsonrpc", func(w http.ResponseWriter, r *http.Request) {
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			ReadCloser: r.Body,
			Writer:     w,
		}
		rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
	})
	http.ListenAndServe(":1234", nil)
}
