package main

// 这种写法的 4 大痛点（你学习必须记住）
// ① CPU 空转爆表（疯狂浪费）
// 读变量要加锁写变量要加锁代码巨啰嗦一不小心就死锁 / 数据竞争
// ② 主协程不能睡觉，只能死等
// 你不能让主协程 sleep，否则反应迟钝；不睡就一直耗 CPU。
// ③ 本质是：暴力轮询 polling
// 不是事件通知，是一遍遍问 “好了吗？好了吗？”
// ④ 有数据竞争，并发不安全
// 两个 goroutine 写，主 goroutine 读，没有任何同步。
import (
	"fmt"
	"sync"
	"time"
)

// 1. 全局变量（共享内存）
var task1Done bool
var task2Done bool

// 2. 必须加锁！否则有数据竞争！
var lock sync.Mutex

func main() {
	// 协程1
	go func() {
		time.Sleep(3 * time.Second)

		lock.Lock() // 写之前必须加锁
		task1Done = true
		lock.Unlock()
	}()

	// 协程2
	go func() {
		time.Sleep(1 * time.Second)

		lock.Lock() // 写之前必须加锁
		task2Done = true
		lock.Unlock()
	}()

	// 3. 主协程：死循环轮询（暴力检查）
	for {
		lock.Lock() // 读也要加锁！
		d1 := task1Done
		d2 := task2Done
		lock.Unlock()

		if d1 {
			fmt.Println("任务1先完成")
			break
		}
		if d2 {
			fmt.Println("任务2先完成")
			break
		}

		// 不休息 CPU 直接跑满！
		// 休息了 → 不实时！
	}
}
