package main

import (
	"fmt"
	"time"
)

func asyncPrint(str string) {
	// time.Sleep(time.Second)
	// println(str)

	for { // 隔一秒就会打印一次
		time.Sleep(time.Second)
		println(str)
	}
}
func main() {
	// 任何一个函数都可以使用go关键字启动一个协程就变成了异步的，放在底层的Goroutine中去执行
	// 主死随从：主程序结束了，从异步程序还没来得及执行就随着主程结束而结束
	// go asyncPrint("hello world") // 写法1
	// // go协程的匿名函数写法，启动goroutine // 写法2
	// go func() {
	// 	time.Sleep(time.Second)
	// 	println("hello world")
	// }()
	// 循环启动多个协程 // 写法3
	for i := 0; i < 5; i++ {
		// 1. 闭包问题：一个函数中引用了外面另外一个作用域变量就产生闭包
		// 2. for循环的问题：每次for循环时，i变量会重用，
		// 解法1:tmp := i
		go func() {
			// for循环套异步的代码这种写法一定会有个bug：你以为输出 0 1 2 3 4，实际不是，而可能是乱的数字：43102
			// 原因：因为协程是异步的代码，里捕获的是变量 i 的地址（引用），不是副本，异步代码不是立即执行，是底层的gmp统一调度执行的，所以不会严格打印01234
			// 解决方法：创建一个副本
			fmt.Println(i)
			// 解法1 的 fmt.Println(tmp) // 这种解法1，只能保证每次打印的数是不会重复，但是不会保证打印的数是顺序的，因为异步调度触发时机是无法保证的。
		}()

		// （最终的）解法2:使用传参：参数是值传递
		go func(i int) {
			fmt.Println(i) // 这种是顺序打印，且不会重复，数字正确
		}(i)
	}
	fmt.Println("main")
	time.Sleep(10 * time.Second) // 这个睡眠时间可以保证主死随从的从先执行完，
	// 执行后，先打印main，后再打印hello world

}
