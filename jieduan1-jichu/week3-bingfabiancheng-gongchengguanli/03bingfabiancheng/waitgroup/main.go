package main

import (
	"fmt"
	"sync"
	"time"
)

// 学会使用wait group去等待协程结束

func asyncPrint(str string) {
	// time.Sleep(time.Second)
	// println(str)

	for { // 隔一秒就会打印一次
		time.Sleep(time.Second)
		println(str)
	}
}
func main() {
	var wg sync.WaitGroup
	// 我要监控多少个goroutine执行结束
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(i int) {
			defer wg.Done()
			fmt.Println(i) // 这种是顺序打印，且不会重复，数字正确
			// wg.Done()      // 为了防止忘写，可以在协程函数顶部改成写法defer wg.Done()
		}(i)
	}
	fmt.Println("main")
	// time.Sleep(10 * time.Second) 这种写法不优雅，靠睡眠保证子goroutine执行完，借助wait group来精确控制goroutine结束时间
	wg.Wait()

	// wait group主要用于goroutine的执行等待，Add方法必须与Done方法配套成对使用

}
