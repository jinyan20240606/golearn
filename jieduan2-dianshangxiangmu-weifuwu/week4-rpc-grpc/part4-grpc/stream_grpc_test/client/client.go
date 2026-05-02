package main

import (
	"context"
	"fmt"
	"golearn/part4-grpc/stream_grpc_test/proto"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("127.0.0.1:50052", grpc.WithInsecure())

	if err != nil {
		panic("连接失败" + err.Error())
	}
	defer conn.Close()

	c := proto.NewGreeterClient(conn)
	// 服务端流模式
	res, _ := c.GetStream(context.Background(), &proto.StreamReqData{Data: "jinyan"})
	for {
		data, err := res.Recv()
		if err != nil {
			break
		}
		println(data.Data)

	}

	// 客户端流模式
	cliStr, _ := c.PutStream(context.Background())
	for i := 0; i < 10; i++ {
		cliStr.Send(&proto.StreamReqData{Data: "jinyan" + strconv.Itoa(i)})
	}

	// 双向流模式通信
	allStr, _ := c.AllStream(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			data, _ := allStr.Recv()
			fmt.Println("收到服务端数据：", data.Data)
		}
	}()

	go func() {
		defer wg.Done()
		for {
			_ = allStr.Send(&proto.StreamReqData{
				Data: "我是客户端，发送数据",
			})
			time.Sleep(time.Second)
		}
	}()

	wg.Wait()
}
