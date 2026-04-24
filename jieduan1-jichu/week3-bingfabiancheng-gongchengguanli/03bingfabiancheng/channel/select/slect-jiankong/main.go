package main

import (
	"fmt"
	"time"
)

func g1(ch1 chan struct{}) {
	time.Sleep(2 * time.Second)
	ch1 <- struct{}{} // 空结构体实例化
	// ch1 <- struct{}{} //
}

func g2(ch2 chan struct{}) {
	time.Sleep(3 * time.Second)
	ch2 <- struct{}{} // 空结构体实例化
}

func main() {
	// 纯信号通道值类型的选型：默认会想到bool类型占 1 字节，比较省空间，更常用空结构体占用 0 字节，纯信号通知、不需要传任何数据 → 统一用：chan struct{}
	// 定义2个全局的channel，分别对应两个任务，很多时候不同的goroutine一般是写入不同的channel，不同goroutine不推荐写同一个channel

	ch1 := make(chan struct{}, 1) // 这块channel用协程消费时，写成无缓冲有缓冲都行
	ch2 := make(chan struct{})    // 默认就是0无缓冲

	go g1(ch1)
	go g2(ch2)

	// 🔥 一句 select 搞定：不需要加锁、无全局变量、无轮询
	// 规则只有3条
	// 1. 遍历所有 case，看哪个通道能读 / 能写，找到一个就绪的(超时先就绪或channel信号先就绪等)，就执行它
	// 2. 如果多个同时就绪 → 随机选一个执行
	//     - 不会按代码顺序执行第一个 case，而是随机选一个！Go 官方设计原因：防止总是优先执行前面的 case，导致后面的通道 “饿死”保证公平性
	// 3. default
	// 4. 超时设置
	tc := time.NewTimer(1 * time.Second) // 创建一个定时器，2 秒后，定时器会自动往 tc.C 通道发一个信号，select 捕捉到这个信号 → 判定超时
	select {
	case <-ch1: // 不同的case对应channel取值
		fmt.Println("任务1完成")
	case <-ch2:
		fmt.Println("任务2完成")

	case <-tc.C: // 超时处理
		fmt.Println("任务超时")
		return

	// 不加 default：select 会一直阻塞等待（卡住）永远不会执行select语句后面的语句
	// 加 default：select 永远不会阻塞，进入default后，就相当于进入了select的结束语句，然后继续往下执行select语句后面的语句
	// 协程还在 只要sleep在延时睡眠执行，没立即发信号 → 通道空 → 所有 case 都不满足 → 就会一直触发 default
	default: // default：
		fmt.Println("任务未完成")
	}

	// 可以很快打印出 任务2完成

	// 注意细节：default与超时不要一起用，因为：default 优先级 = 最高！还没有等到超时就直接default了
}

func main1() {
	// 演示 select 对两个channel同时就绪时的随机执行效果
	ch1 := make(chan struct{}, 1)
	ch2 := make(chan struct{}, 1)

	// 这是两个channel同时就绪--同时写值，这时会触发select的随机执行
	ch1 <- struct{}{}
	ch2 <- struct{}{}

	select {
	case <-ch1: // 不同的case对应channel取值
		fmt.Println("任务1完成")
	case <-ch2:
		fmt.Println("任务2完成")
	}
}
