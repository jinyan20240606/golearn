package main

import (
	"fmt"
)

type Message string

type Greeter struct {
	Message Message
}

type Event struct {
	Greeter Greeter
}

func NewMessage() Message {
	return "Hello there!"
}

func NewGreeter(m Message) Greeter {
	return Greeter{Message: m}
}

func (g Greeter) Greet() Message {
	return g.Message
}

// 每个Provider都是一个构造函数，构造函数的参数是其他Provider的返回值，构造函数的返回值是需要注入的对象
func NewEvent(g Greeter) Event {
	return Event{Greeter: g}
}

func (e Event) Start() {
	msg := e.Greeter.Greet()
	fmt.Println(msg)
}

func main() {
	/*
		构建过程， ioc容器就是要解决这个初始化问题
	*/
	event := initEvent()

	event.Start()
}
