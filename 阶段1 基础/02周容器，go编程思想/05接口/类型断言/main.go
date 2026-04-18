package main

import (
	"fmt"
	"strings"
)

// / 接口数据类型

// 一个结构体可以实现多个接口，这时候就需要用到断言了，我有时候需要知道这个结构体实现了哪些接口，断言，是否实现了某一个接口
// 区别：
// a.(type):类型判断 -> 只能用在 switch 里 → 用来判断是什么类型
// a.(int): 类型断言 → 把空接口变成具体的 int。。。可以用在任何地方 → 用来取出 / 转换成 int 类型

// func add(a, b int) int {
// 	return a + b
// }

// func add1(a, b int8) int8 {
// 	return a + b
// }
// 这样太麻烦了，需要为每种参数类型定义不同的方法，改成如下比较好

func add(a, b interface{}) int { // 先用空接口代表any，接收任意类型，在内部做类型断言处理
	// 这就是类型断言，判断是不是目标int类型，是则返回(其实就是个语法糖)
	a1, ok := a.(int)
	if !ok {
		panic("参数类型错误")
	}
	b1, ok := b.(int)
	if !ok {
		panic("参数类型错误")
	}
	fmt.Printf("%T\n", a1) // int
	return a1 + b1
}

// 通过switch语法进行类型判断，接收什么类型就返回什么类型的加法处理
func add1(a, b interface{}) interface{} {
	// switch v := a.(type) {  // 👈 这里一步到位--- 这是简写语法，判断并赋值

	// switch ...(type) 只做判断，不做赋值
	switch a.(type) {
	// 此时的a 还是 interface {}，不转成具体类型不能运算。
	case int:
		a1 := a.(int) // 👈 转成具体类型
		b1 := b.(int)
		return a1 + b1
	case int8:
		a1 := a.(int8)
		b1 := b.(int8)
		return int(a1) + int(b1)
	case float64:
		a1 := a.(float64)
		b1 := b.(float64)
		return a1 + b1
	case string:
		a1 := a.(string)
		b1 := b.(string)
		return a1 + b1
	default:
		panic("参数类型错误")
	}
}

func main() {
	a := 1
	b := 2
	fmt.Println(add(a, b)) // 3

	// 字符串情况
	re := add1("nihao", "world")
	// 这取出来还是默认是interface类型，需要再次类型断言转换成具体类型去运算
	res, _ := re.(string)
	fmt.Println(re, strings.Split(res, " "))

}
