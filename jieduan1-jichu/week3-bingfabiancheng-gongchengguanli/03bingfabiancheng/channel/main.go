package main

// channel，channel类型是整体类型

import (
	"fmt"
)

// 有缓冲时的传输解法
func main() {
	fmt.Println("hello world")

	// 语法解释：chan string 是一个整体的带子类型的channel类型，msg 的类型是：chan string（表明一个通道类型，通道传递的值是string类型）
	var msg chan string // channel的默认值就是nil

	// make可以初始化map，slice也可以初始化我们的hannel
	msg = make(chan string, 1) // 创建一个容量为1的channel,如果容量为0，则channel是阻塞的，则会永远放不进去，会触发死锁
	// channel 有缓冲和无缓冲的channel是不一样的，容量大于0就是有缓冲，为0的就是无缓冲的

	// 符号-左侧读，右侧写, 记住放值和取值的符号方向
	msg <- "hello" // 放值到channel中, 通道符号发送和接受都是一样的 <-

	data := <-msg // 从channel中取值
	// 若此时chan容量0阻塞的，执行完会阻塞会触发死锁
	// 死锁原因：先发送，然后通道容量为0-阻塞，永远没有人接受触发死锁

	if msg != nil {
		fmt.Println(data)
	}

}

// 当无缓冲时的传输解法：若想成功发送和接收，则需要先启动一个协程去接受，然后才能发送
func main1() {

	var msg chan string
	// 错误 msg <- "hello" ❌ 发送写在了启动协程之前，主协程直接卡死，子协程永远跑不起来。
	msg = make(chan string, 0)
	// ✅ 第一步：先启动协程（让它等着接收）
	go func(msg chan string) {
		fmt.Println("hello world")
		data := <-msg
		fmt.Println(data)
	}(msg)
	// ✅ 第二步：再发送（这时候有人等着，不会卡死）
	msg <- "hello"

	// 为什么用协程能解决？涉及到Go有一种happen-before的机制
}
