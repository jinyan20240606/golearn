package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// 学习context:上次是使用比较优雅的channel方式，我们还可以更优雅的实现：使用context

var wg sync.WaitGroup

// 我们有新的需求，在主goroutine通知子goroutine中主动退出子程序

func cpuInfo(ctx context.Context) { // 改成用参数传递的方式 --> O1
	// 打印祖先上下文的追踪id
	fmt.Printf("traceid: %s\r\n", ctx.Value("traceid"))
	// 可以记录一些日志，这次请求是哪个请求打印的
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("退出cpu监控")
			return
		default: // 不加default，select会阻塞，永远不会执行seelct后面的语句
			time.Sleep(1 * time.Second)
			fmt.Println("cpu的信息")
		}
		fmt.Println("select语句执行完了")
	}

}
func main() {
	// 渐进式的方式学习context
	// 如实现一个需求：写一个goroutine监控的cpu的信息

	wg.Add(1)
	// 方式1、-创建取消上下文：context有根据顶级上下文构造出子功能context
	// ctx1是父cancel上下文
	// ctx1, cancel1 := context.WithCancel(context.Background()) // 创建一个可取消的上下文，爸爸是根节点 context。
	// // ctx2是子cancel上下文
	// ctx2, _ := context.WithCancel(ctx1)
	// go cpuInfo(ctx2) // 我传递了子cancel2上下文，然后下面用父的cancel1取消，也是可以的有父到子的传递性：继承父 ctx 特性，父取消 → 全部子 ctx 连锁取消

	// 方式2、创建一个主动超时的上下文
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	// time.Sleep(5 * time.Second)
	// cancel1() // 设置主动超时后，就不用写手动的取消函数了

	// 方式3、创建一个定时器超时的上下文使用WithDeadline，在时间点定时取消，略讲

	// 方式4、使用WithValue；上下文中通过传值，	构造函数仅返回一个值
	valueCtx := context.WithValue(ctx, "traceid", "traceidval-1234") // 二参是key，三参是value
	go cpuInfo(valueCtx)

	wg.Wait() // 这目前是永远等待的，因为cpuInfo在无限循环
	fmt.Println("main goroutine done")
}
