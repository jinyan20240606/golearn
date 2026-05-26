package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func main() {
	// 有推和拉两个模式
	// 拉模式是：客户端不停的向主服务器上拉，看有没有新数据，不停轮询，有一些性能浪费
	// 推模式是：服务器有了数据之后，服务器会主动推到这边来，更加省资源---- 目前用的最多
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"192.168.0.104:9876"}),
		// 消费者组（Consumer Group）= 让该组下多个消费者实例，共同消费同一个 Topic，实现并发处理、负载均衡、高可用。
		// 把当前这个消费者，加入到名叫 mxshop 的消费者组里进行负载均衡。
		consumer.WithGroupName("mxshop"),
	)

	// ... 是类型的一部分，它专门用来表示「可变参数类型」，msgs就是一个切片类型[]*primitive.MessageExt
	if err := c.Subscribe("imooc1", consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			fmt.Printf("获取到值： %v \n", msgs[i])
		}
		return consumer.ConsumeSuccess, nil
	}); err != nil {
		fmt.Println("读取消息失败")
	}
	_ = c.Start()
	//注意最重要的一点：不能让主goroutine退出
	time.Sleep(time.Hour)
	_ = c.Shutdown()
}
