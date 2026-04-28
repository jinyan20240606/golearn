package main

import "net/rpc" // 默认go的rpc包，使用gob 进行序列化

func main() {
	// 1、建立一个连接,用Dial发起连接

	client, err := rpc.Dial("tcp", "localhost:1234")

	if err != nil {
		panic(err)
	}

	var reply string // 虽然只声明没赋值，有默认值，有内存占用的，并不是一个单纯的地址变量

	errrr := client.Call("HelloService.Hello", "你好 bobby", &reply)

	if errrr != nil {
		panic(errrr)
	}
	println(reply)

}
