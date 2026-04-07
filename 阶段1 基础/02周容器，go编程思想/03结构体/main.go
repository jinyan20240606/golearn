package main

import (
	"fmt"
	"strconv"
)

// 学习go语言中的结构体：类似于其他语言的class，又比class更加轻量级

/**
*  type关键字的作用
*     1. 定义结构体类型
*     2. 定义接口类型
*     3. 定义类型别名 ：代码编译的时候，类型别名会被替换为原始类型
            - type MyInt = int // 定义一个int的类型别名 MyInt，它是 int 的别名，它打印类型还是显示的int
      4. 定义自定义类型：type MyType int // 这就是一个新的类型叫 MyType，它的底层类型是 int，但它和 int 是不同的类型，不能直接赋值，需要进行类型转换
            - 定义函数类型 // type MyFunc func(int) int // 定义一个新的函数类型 MyFunc，它接受一个 int 参数并返回一个 int
			- 给自定义类型绑定方法的标准语法
			  func (接收者 接收者类型) 方法名(参数列表) 返回值类型 {
				 方法体
			  }
	  5. 类型判断
	     - a.(type)是类型断言的一种特殊形式，专门用于在switch语句中判断接口变量的动态类型

*/

type MyInt int // 4的自定义类型
// 给自定义类型绑定方法的标准语法

func (mi MyInt) string() string {
	// Go 标准库函数：数字转字符串
	return strconv.Itoa(int(mi))
}

func main() {
	var i MyInt = 10
	fmt.Println(i.string())
	fmt.Printf("%T", i)

	// 5.类型判断比如例子
	var a interface{} = "abc"
	switch a.(type) { // a.(type)是类型断言的一种特殊形式，专门用于在switch语句中判断接口变量的动态类型
	case int:
		fmt.Println("a is an int")
	case string:
		fmt.Println("a is a string")
	default:
		fmt.Println("a is of unknown type")
	}

}
