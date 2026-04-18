package main

import "fmt"

// 接口也支持嵌套使用，有利于抽离复用

type MyWriter interface {
	Write(string)
}

type MyReader interface {
	Read() string
}

// 接口组合起来可以凑成个新的接口
type MyReadWriter interface {
	MyWriter
	MyReader
	Close()
}

// 基于上述接口我们用用黄金组合结构体实践一下
type SreadWriter struct {
}

// Close implements [MyReadWriter].
func (s *SreadWriter) Close() {
	fmt.Println("close")
}

// Read implements [MyReadWriter].
func (s *SreadWriter) Read() string {
	fmt.Println("read")
	return ""
}

// Write implements [MyReadWriter].
func (s *SreadWriter) Write(string) {
	fmt.Println("write")
}

func main() {
	// var mrw MyReadWriter = SreadWriter{}// 这个是值类型的结构体初始化
	// 一般情况下我们都会用指针初始化，除非明确你的逻辑适合拷贝方式的值类型才能用
	var mrw MyReadWriter = &SreadWriter{} // 这个是指针类型的结构体初始化
	// 值类型的初始化时，方法接收器必须写值类型
	// 指针类型的初始化时，方法接收器可以写指针类型或者值类型【因为都是兼容的】
	// 方法接收器写指针类型时，结构体初始化必须用指针类型初始化

	mrw.Read()

}
