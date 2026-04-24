package main

import (
	"fmt"
	"sync"
	"time"
)

// 学习context:最原始的方式就是使用全局变量实现,我们可以优雅点，使用channel

var wg sync.WaitGroup

// 我们有新的需求，在主goroutine通知子goroutine中主动退出子程序
// var stop = make(chan struct{}) // 不要写全局channel做信号传递，用参数传递的方式会更加的好 --> O1

func cpuInfo(stop chan struct{}) { // 改成用参数传递的方式 --> O1
	defer wg.Done()
	for {
		select {
		case <-stop:
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
	var stop = make(chan struct{}) // 改成局部变量用参数传递的方式

	wg.Add(1)
	go cpuInfo(stop) // 改成用参数传递的方式

	time.Sleep(5 * time.Second)
	stop <- struct{}{}

	wg.Wait() // 这目前是永远等待的，因为cpuInfo在无限循环
	fmt.Println("main goroutine done")
}
