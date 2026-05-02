package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"golearn/part4-grpc/stream_grpc_test/proto"

	"google.golang.org/grpc"
)

const PORT = ":50052"

type server struct {
	proto.UnimplementedGreeterServer
}

// 服务端流模式的rpc
func (s *server) GetStream(req *proto.StreamReqData, res proto.Greeter_GetStreamServer) error {
	// 这个方法：只要调用send一次，客户端就会收到流式消息，所以不能用普通的return消息了，return是一次性的，send可以持续发送，要做到持续发送，就得加for循环
	i := 0
	for {
		i++
		_ = res.Send(&proto.StreamResData{
			Data: req.Data + " " + strconv.Itoa(i) + strconv.FormatInt(time.Now().Unix(), 10),
		})
		time.Sleep(time.Second)
		if i > 5 {
			break
		}
	}
	return nil
}

// 双向流模式
func (s *server) AllStream(allStr proto.Greeter_AllStreamServer) error {
	// 应该使用协程，一个协程 专门接收客户端消息，一个协程 专门发送消息，不能串行，否则会阻塞卡死。

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			data, _ := allStr.Recv()
			fmt.Println("收到客户端数据：", data.Data)
		}
	}()

	go func() {
		defer wg.Done()
		for {
			_ = allStr.Send(&proto.StreamResData{
				Data: "服务器返回数据",
			})
			time.Sleep(time.Second)
		}
	}()

	wg.Wait()

	return nil
}

// 客户端流模式的rpc
func (s *server) PutStream(cliStr proto.Greeter_PutStreamServer) error {
	for {
		req, err := cliStr.Recv()
		if err != nil {
			return err
		}
		fmt.Println(req.Data)
	}
	return nil
}
func main() {
	lis, err := net.Listen("tcp", PORT)

	if err != nil {
		panic("监听失败" + err.Error())
	}

	g := grpc.NewServer()

	proto.RegisterGreeterServer(g, &server{})

	if err = g.Serve(lis); err != nil {
		panic("启动失败" + err.Error())
	}
}
