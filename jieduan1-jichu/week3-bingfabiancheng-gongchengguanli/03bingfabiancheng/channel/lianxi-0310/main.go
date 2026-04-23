package main

import (
	"fmt"
	"time"
)

var number, letter = make(chan bool), make(chan bool)

// 必须得用2个channel给两个方法交互，只用1个双向的话，会死循环

func printNum() {
	i := 1
	for {
		// 怎么做到交叉打印----需要等待另外一个channel通知我
		<-number // 左侧读个值，相当于从number中接收信号
		fmt.Printf("%d%d", i, i+1)
		i += 2
		letter <- true // 右侧写个值，相当于向letter通道中发送信号
	}
}

func printLetter() {
	i := 0
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for {
		// 怎么做到交叉打印----需要等待另外一个channel通知我
		<-letter // 左侧读个值，相当于从letter中接收信号
		if i >= len(str) {
			return
		}
		fmt.Print(str[i : i+2]) // 这是切片写法切片：从 i 开始，取 2 个字符
		// string(str[0]) // str[0]取出来是Ascii码值，需要转换成字符
		i += 2
		number <- true // 右侧写个值，相当于向number通道中发送信号
	}
}

func main() {
	// 链习题：使用2个goroutine实现交叉打印序列，一个goroutine打印数字，另外一个goroutine打印字母，效果如下：12AB34CD56EF78GH...

	go printNum()
	go printLetter()

	number <- true // 通知第一个goroutine开始执行

	time.Sleep(3 * time.Millisecond)

}
