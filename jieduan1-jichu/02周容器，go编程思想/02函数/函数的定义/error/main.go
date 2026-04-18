package main

import (
	"errors"
	"fmt"
)

// go 语言中的错误处理：与其他语言有很大的不同

// error panic recover

// error：
// 一个函数可能出错，其他语言是trycatch去抱住这个函数
// go语言中没有trycatch，函数返回错误，错误处理是通过函数的返回值来处理的，函数的返回值可以是一个error类型，调用者通过判断这个error是否为nil来确定函数是否执行成功
// go设计者认为必须要处理这个error错误，--- 也叫防御性编程

// panic是一个内置函数：这个函数会让你的程序退出，不推荐随便使用panic
// ----- 一版在哪用到：比如一个服务的启动过程中，有些依赖服务必须要准备好，mysql联通，redis联通，如果这些依赖服务没有准备好，那么这个服务就没有什么意义了，这时候可以使用panic来让程序退出，避免程序在一个不正常的状态下运行，导致更严重的问题

// recover内置函数：这个函数会捕获panic，然后返回一个error

func A() (int, error) {
	// panic("自定义错误") // 这个函数会让你的程序退出，后面的语句不会执行
	defer func() {
		// recover内置函数：这个函数会捕获被动和主动的panic，然后返回一个error（如代码写的语法错误能捕获到）
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	// 1. defer必须写在最前面recover必须写在defer里才能生效，才能捕获到后面代码的panic，否则defer写在后面就捕获不到前面代码的直接panic了，后面代码都不执行了
	// 2、recover处理异常后，程序不会退出了，继续往下执行
	// 3、多个defer会行成栈，后进的先执行

	return 0, errors.New("自定义错误1")
}
func main() {
	if _, err := A(); err != nil {
		fmt.Println(err)
	}

}
