package main

import "fmt"

// 结构体实现接口的方式
//   简单写法：--- 一个结构体实现一个接口，可以直接在结构体的接收器定义接口所需的方法
//   嵌套写法：--- 还可以在结构体定义接口匿名字段，在结构体初始化时赋值给接口匿名字段一个符合该接口的结构体

// 提示：一个结构体可以实现多个接口，一个接口也可以被多个结构体实现

// 简单写法如下：
type MyWriter1 interface {
	Write(p []byte)
}

type MyCloser1 interface {
	Close()
}

type MyReadWriteCloser1 struct {
}

func (rwc MyReadWriteCloser1) Write(p []byte) {
	fmt.Println(2)
}
func (rwc MyReadWriteCloser1) Close() {
	fmt.Println(3)
}

func main() {

	var mw MyWriter1 = MyReadWriteCloser1{}
	var mc MyCloser1 = MyReadWriteCloser1{}
	fmt.Println(mw, mc)
}
