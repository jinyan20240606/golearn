package main

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.0.104:9876"}))
	if err != nil {
		panic("生成producer失败")
	}

	if err = p.Start(); err != nil {
		panic("启动producer失败")
	}

	msg := primitive.NewMessage("imooc1", []byte("this is delay message"))
	// 设置延迟的级别--- 就触发了延时消息
	msg.WithDelayTimeLevel(3)
	res, err := p.SendSync(context.Background(), msg)
	if err != nil {
		fmt.Printf("发送失败: %s\n", err)
	} else {
		fmt.Printf("发送成功: %s\n", res.String())
	}

	if err = p.Shutdown(); err != nil {
		panic("关闭producer失败")
	}

	// 为什么要使用延时消息？--- 常用来解决超时归还的问题
	//支付的时候， 淘宝， 12306， 购票，为了防止下单有人恶意不支付就锁你单 ---- 就需要有个超时归还机制 -- 那么多订单判断超时归还就对应定时轮询执行检查是否超时的一个逻辑

	//我可以去写一个轮询去定时查看哪些订单超时了， 轮询的问题： 1. 多久执行一次轮询 30分钟
	//在12:00执行过一次， 下一次执行就是在 12:30的时候 但是12:01的时候下了单， 12:31就应该超时 13:00时候才能超时
	//那我1分钟执行一次啊， 比如我的订单量没有这么大，1分钟执行一次， 其中29次查询都是无用， 而且你还还会轮询mysql

	//rocketmq的延迟消息， 1. 时间一到就执行，不会包含无用轮询， 2. 消息中包含了订单编号，你只查询这种订单编号

}
