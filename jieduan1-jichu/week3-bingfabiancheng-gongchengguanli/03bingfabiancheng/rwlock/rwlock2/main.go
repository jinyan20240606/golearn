package main

// 读写锁：实例演示下读写锁的阻止效果

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// var num int
	var rwlock sync.RWMutex
	var wg sync.WaitGroup
	// 一共3个方法：读锁，写锁，解锁

	wg.Add(2)
	// 写的goroutine
	go func() {
		time.Sleep(1 * time.Second)
		defer wg.Done()
		rwlock.Lock() // 写锁，写锁会防止别的写锁获取，和读锁获取
		defer rwlock.Unlock()
		fmt.Println("带写锁的写请求")
		// num = 12
		time.Sleep(5 * time.Second)
	}()

	// for i := 0; i < 5; i++ { // 加也行不加也行，就是并发读协程的数量

	// 读的goroutine
	go func() {
		defer wg.Done()
		for { // 它是一个无限循环，让每个读协程一直不停读
			rwlock.RLock()                     // 读锁，读锁不会阻止别人的读
			time.Sleep(500 * time.Millisecond) // 每半秒拿一次读锁
			fmt.Println("带读锁的读请求")
			// fmt.Println(num)
			rwlock.RUnlock() // for循环中，就不应该写defer了，放在最下面
		}

	}()

	// }

	wg.Wait()

	/**
	上面的打印效果：

	带读锁的读请求 - 最开始每半秒触发读一次
	带读锁的读请求
	带写锁的写请求 - 1秒后触发，待了5秒钟，后，才继续触发读的请求
	带读锁的读请求
	带读锁的读请求
	带读锁的读请求
	带读锁的读请求
	带读锁的读请求
	带读锁的读请求
	*/
}
