package main

import (
	"errors"
	"fmt"
)

func mPrint(datas ...interface{}) {
	for _, value := range datas {
		fmt.Println(value)

	}
}

type myInfo struct{}

func (mi *myInfo) Error() string {
	return "自定义错误1"
}

func main() {
	var data = []interface{}{
		1,
		"hello",
		true,
	}
	var data1 = []string{"hello", "world"}
	mPrint(data...)
	// 01-1-直接将data1解构传参会报错
	// mPrint(data1...) // 这样不可以报错：在解构下每个参数必须类型对应。。。。如果非解构普通传入，any接口类型是能接受string的则可以

	// 01-2-可以通过间接的方式实现将data1解构传参
	var data2 []interface{}
	for _, value := range data1 {
		data2 = append(data2, value)
	}
	mPrint(data2...)

	// 02- error的本质就是个接口类型，只要实现了error接口，那么就可以作为赋值error类型
	var err error
	err = errors.New("自定义错误")
	fmt.Println(err) // 自定义错误
	err = &myInfo{}
}
