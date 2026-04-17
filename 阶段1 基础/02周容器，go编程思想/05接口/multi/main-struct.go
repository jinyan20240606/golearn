package main

import "fmt"

// 嵌套写法如下：结构体实现接口通过内嵌接口匿名字段实现，赋值字段一个符合接口的子结构体
// 这个写法必须得学会，对于接口解耦很好

// 接口
type MyWriter interface {
	Write(p []byte)
}

type MyCloser interface {
	Close() error // 它是函数类型，且返回值是接口类型，然后接口的零值 = nil，所以 error 类型变量可以 = nil
}

// 2个结构体
type writerCloser struct { // 这是结构体的匿名字段方式定义方式
	MyWriter // 结构体嵌入接口 =相当于 强制这个结构体必须至少包含实现该接口
}

type fileWriter struct {
	filePath string
}

// 接收器
func (fw fileWriter) Write(p []byte) {
	fmt.Println("写string到文件file")
}

func (dw writerCloser) Close() error {
	fmt.Println(3)
	return nil
}

func main1() {

	var mw MyWriter = writerCloser{
		// Go 对「只有一个匿名字段时」的简写语法：匿名嵌入类型 MyWriter，简写初始化可以直接写值，不用写字段名
		// 它默认就对应那个匿名字段的类型：MyWriter
		fileWriter{},
	}
	mw.Write([]byte("hello world"))
	fmt.Println(mw)
}
