package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/**
 * 锁住要是解决对共享资源的竞争访问冲突
 */

var total int32
var wg sync.WaitGroup

// 必须使用同一把锁才能保证数据安全，千万不要去复制锁，就失去了锁的效果
// var lock sync.Mutex

func add() {
	defer wg.Done()
	for i := 0; i < 1000000; i++ {
		atomic.AddInt32(&total, 1) // 除了锁，还可以用原子操作方法去进行数值
		// lock.Lock()
		// total++
		// lock.Unlock()
	}
}

func sub() {
	defer wg.Done()
	for i := 0; i < 1000000; i++ {
		atomic.AddInt32(&total, -1)
		// lock.Lock()
		// total--
		// lock.Unlock()
	}
}

func main() {
	wg.Add(2)
	go add()
	go sub()
	wg.Wait()
	fmt.Println(total)

	// 当前单一运行go add()【1000000】或go aub()【-1000000】 的协程是正常的，没问题
	// 但是，如果两个协程都一起运行，就会产生数据资源竞争，导致结果不准确【-350224】每次值都不一样，就需要加锁了
}
