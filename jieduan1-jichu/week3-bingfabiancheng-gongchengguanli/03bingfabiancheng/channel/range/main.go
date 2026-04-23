package main

// 讲
import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("开始----")

	var msg chan string // channel的默认值就是nil

	msg = make(chan string, 1) // 创建一个容量为1的channel,如果容量为0，则channel是阻塞的，则会永远放不进去，会触发死锁
	go func(msg chan string) {
		for data := range msg {
			// forrange语法支持对channel遍历，自动取值， 会一直循环，直到channel关闭
			// 使用forrange时会一直循环取值，没有值了还会一直等待（所以看不到打印all done），直到channel关闭
			fmt.Println(data)
		}
		// 不用forrange，得多次手动取值
		// data := <-msg
		// fmt.Println(data)
		// data = <-msg
		// fmt.Println(data)
		fmt.Println("all done")

	}(msg)
	msg <- "1" // 放值到channel中, 通道符号发送和接受都是一样的 <-
	msg <- "2"
	close(msg)
	// close(msg) // 关闭channel，这个与其他语言有区别，其他语言好像不能关闭一个队列
	// 		Java：BlockingQueue 没有关闭功能
	// Python：Queue 没有关闭功能
	// C++：队列也没有关闭
	// 只有 Go：可以主动关闭通道，并让所有接收者感知！
	// 这是 Go 为了优雅退出、广播通知专门设计的超级功能

	// 关闭后的，就不能再往通道中发送数据了，已经关闭的channel可以继续取值，但是不能再发送数据了

	time.Sleep(3 * time.Second) // 简单阻塞几秒，避免主协程提前退出，一般会用waitgroup解决

}
