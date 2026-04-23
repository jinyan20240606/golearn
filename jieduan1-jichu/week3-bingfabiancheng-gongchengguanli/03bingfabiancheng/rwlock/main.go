package main

// 读写锁 ：基本使用示例

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var num int
	var rwlock sync.RWMutex
	var wg sync.WaitGroup
	// 一共3个方法：读锁，写锁，解锁

	wg.Add(2)
	// 写的goroutine
	go func() {
		defer wg.Done()
		rwlock.Lock() // 写锁，写锁会防止别的写锁获取，和读锁获取
		defer rwlock.Unlock()
		fmt.Println("写锁")
		num = 12
	}()

	//【A操作】
	time.Sleep(1 * time.Second) // 加上这个就能伪实现顺序执行，了，最终输出12不是0

	// 读的goroutine
	go func() {

		defer wg.Done()
		rwlock.RLock() // 读锁，读锁不会阻止别人的读
		defer rwlock.RUnlock()
		fmt.Println("读锁")
		fmt.Println(num)
	}()
	wg.Wait()

	// 执行完，终端输出打印0.，goroutine执行时，并没有保证协程定义的额先后顺序，所以要保证输出的顺序得借助后面通信的方法还没讲，先在中间sleep以下见【A操作】，也能暂时实现
}
