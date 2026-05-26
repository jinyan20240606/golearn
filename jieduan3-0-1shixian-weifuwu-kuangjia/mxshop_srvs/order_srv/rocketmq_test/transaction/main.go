package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

// 事务监听器
// 必须实现 2 个方法：
// 1. ExecuteLocalTransaction 执行本地事务
// 2. CheckLocalTransaction    事务回查（MQ主动问你业务成没成功）
type OrderListener struct{}

// 一、执行本地事务（核心方法）
// 作用：发送半消息成功后，执行本地业务逻辑（订单/库存/支付）
func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	fmt.Println("开始执行本地逻辑")
	// 模拟业务执行（比如：创建订单、扣减库存、支付）
	time.Sleep(time.Second * 3)
	// // 测试本地执行成功
	// return primitive.CommitTransaction

	// 测试本地执行失败
	fmt.Println("执行本地逻辑失败")
	//本地执行逻辑无缘无故失败 代码异常 宕机
	// 这里返回 UnknowState（未知状态）
	// 代表：本地事务执行结果不确定（可能挂了、异常、没返回）
	// 触发：RocketMQ 后续会主动回查你的业务！
	return primitive.UnknowState
}

// 二、事务回查方法（兜底保证）--- 它不是立即执行的，当ExecuteLocalTransaction长时间没有返回状态或者返回UnknowState，才会触发这个回查方法
// 隔一段时间就会调用这个，一直得到这个方法的响应，回查会反复、定时、重试，直到拿到明确结果：提交 or 回滚 **

// 什么时候触发？
// 上面的 ExecuteLocalTransaction 没有明确返回 Commit/Rollback
// 比如：程序崩溃、网络断了、超时没响应 → MQ 主动问你
// 这个兜底保证的严谨之处还在于：当这个回查方法也挂掉了，它回查不进来，当你下次正常启动后，还依然重新触发你这个回调方法
func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	fmt.Println("rocketmq的消息回查")
	// 模拟去数据库查询订单状态
	time.Sleep(time.Second * 15)
	// 回查结果：本地事务执行成功！
	// 返回 Commit → 让 MQ 把消息投递给消费者
	return primitive.CommitMessageState
}

func main() {
	// 新建一个事务的Producer，1参需要实现2个钩子方法，2参是nameserver地址
	p, err := rocketmq.NewTransactionProducer(
		&OrderListener{}, // 事务监听器
		producer.WithNameServer([]string{"192.168.0.104:9876"}),
	)
	if err != nil {
		panic("生成producer失败")
	}

	if err = p.Start(); err != nil {
		panic("启动producer失败")
	}

	// 3. 发送【事务消息】（半消息）
	// 此时消息不会给消费者，只会暂存在 MQ
	res, err := p.SendMessageInTransaction(context.Background(), primitive.NewMessage("TransTopic", []byte("this is transaction message2")))
	// ExecuteLocalTransaction 必须执行之后，这个SendMessageInTransaction结果才会返回
	// 它内部不是发完就返回，而是严格按这个顺序阻塞执行：
	// 1. 发送 Half 消息（prepare）到 MQ
	// 	消息先存到 MQ
	// 	消费者看不见
	// 2. 阻塞等待 → 调用你的 ExecuteLocalTransaction
	// 	必须等这个方法执行完
	// 	必须等它返回：Commit / Rollback / UnknowState
	// 3. 拿到本地事务的返回结果后
	// 4. SendMessageInTransaction 才会返回 res/err
	if err != nil {
		fmt.Printf("发送失败: %s\n", err)
	} else {
		fmt.Printf("发送成功: %s\n", res.String())
	}

	time.Sleep(time.Hour) // 这块必须多睡一会测试，流出时间测试事务回查
	if err = p.Shutdown(); err != nil {
		panic("关闭producer失败")
	}
}
