package main

import (
	"fmt"
	"sync"
	"time"
)

// 学习context

var wg sync.WaitGroup

// 我们有新的需求，在主goroutine通知子goroutine中主动退出子程序：最原始的方式就是使用全局变量实现

var stop bool

func cpuInfo() {
	defer wg.Done()
	for {
		if stop {
			break
		}
		time.Sleep(1 * time.Second)
		fmt.Println("cpu的信息")
	}

}
func main() {
	// 渐进式的方式学习context
	// 如实现一个需求：写一个goroutine监控的cpu的信息

	wg.Add(1)
	go cpuInfo()

	time.Sleep(5 * time.Second)
	stop = true

	wg.Wait() // 这目前是永远等待的，因为cpuInfo在无限循环
	fmt.Println("main goroutine done")
}
