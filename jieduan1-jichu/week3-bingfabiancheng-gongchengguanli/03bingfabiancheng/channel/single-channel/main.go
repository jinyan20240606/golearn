package main

import (
	"fmt"
	"time"
)

// 主要讲单向channel的使用场景和用法

// 应用单向channel
// 生产者：参数只能接收一个只写的单向channel
func producer(out chan<- int) {
	for i := 0; i < 10; i++ {
		out <- 1
	}
	close(out)

}

// 消费者：只能读取的只读单向channel
func consumer(in <-chan int) {
	// in <- 1 // 不能写入数据会报错
	for num := range in {
		fmt.Println("num=%d\r\n", num)
	}
}

func main() {
	// 定义前提：单向 channel 必须依赖双向 channel 转换而来，不能直接创建！
	// 定义单向channel的写法：符号-左侧只读，右侧只写，不带符号就是双向

	// 写法1:

	// var ch1 chan int       // 必须先建一个双向channel，然后根据双向channel创建单向channel
	// var ch2 chan<- float64 // 单向channel：只能发送(写入通道)的channel，只能写入float64的数据
	// var ch3 <-chan int     // 单向channel：只能接收(读取管道)的channel，只能读取int数据

	// 写法2:
	c := make(chan int, 3)
	var send chan<- int = c // 定义只能发送的单向channel
	var read <-chan int = c // 定义只能读取的单向channel

	send <- 1
	// <-read // 不写变量名不赋值也可以，直接消费管道
	c2 := <-read

	// read <- 2 这样就会报错，违反了单向的规则
	println(c2) // 1

	// 不能将单向channel转化成双向channel
	// d1 := (chan int)(send)

	ch := make(chan int)
	go producer(ch) // 传入双向，函数参数接收按照类型会自动转成单向
	go consumer(ch)

	time.Sleep(time.Second * 3)

}
