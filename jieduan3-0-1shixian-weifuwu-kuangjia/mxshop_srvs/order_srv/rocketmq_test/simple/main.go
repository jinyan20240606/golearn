package main

// 发送普通消息的示例
import (
	"context"
	"fmt"

	// 你引入的这三个包，是 Apache 官方 Go 客户端
	"github.com/apache/rocketmq-client-go/v2"
	// 基础工具包，创建消息 NewMessage()，消息结构体、错误码、配置常量
	"github.com/apache/rocketmq-client-go/v2/primitive"
	// 生产者专用包，配置生产者，同步 / 异步 / 事务发送
	"github.com/apache/rocketmq-client-go/v2/producer"
)

// 用 Go 语言创建一个 RocketMQ 生产者（Producer），同步发送一条普通消息到 MQ。
func main() {
	// 创建一个 RocketMQ 生产者实例，NameServer的默认端口为9876，NameServer：MQ 的地址管理器（相当于通讯录）
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.0.104:9876"}))
	if err != nil {
		panic("生成producer失败")
	}

	// 启动客户端，建立与 MQ 的连接，启动失败直接崩溃（panic）
	if err = p.Start(); err != nil {
		panic("启动producer失败")
	}

	// SendSync = 同步发送，发出去 → 等待 MQ 返回确认 → 才继续执行，，，发消息一定要发到对应的topic名字里去，NewMessage1参就是topic名字，2参是消息体
	res, err := p.SendSync(context.Background(), primitive.NewMessage("imooc1", []byte("this is imooc1")))
	if err != nil {
		fmt.Printf("发送失败: %s\n", err)
	} else {
		fmt.Printf("发送成功: %s\n", res.String())
	}

	if err = p.Shutdown(); err != nil {
		panic("关闭producer失败")
	}

}
