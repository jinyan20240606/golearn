package main

import (
	"net"
	"net/rpc" // go 内置的rpc包
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
		go rpc.ServeConn(conn)
	}

	// rpc就是把callid映射到函数

	// 这块需要写循环，相关注意事项：
	// 1. Go TCP 最底层：`net.Listen() + Accept()`
	// // 这是操作系统原始 API，accept 是「单次动作」，不是「持续监听服务」
	// // 只负责：开门、等一个人,不自带循环,必须自己写 for {} 才能一直等
	// // 底层本质：极简设计原则（Unix/Linux 设计哲学），每个系统调用，只做一件最小、单一的事，accept：只负责接一位客人，至于「接完一位再接下一位」——交给用户代码自己控制
	// 2. Go HTTP 是高级封装
	// // 这是Go 官方帮你封装好的成品服务器,内部已经写好了 for 循环 + Accept + 协程,你看不到，但它一直在运行,所以你不用写 for
	// 3. Node.js 全是高级封装
	// // Node 不管 TCP 还是 HTTP,底层全自带 事件循环 (EventLoop),你完全不用管循环,所以永远不用写 for
}
