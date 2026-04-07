package main

import "fmt"

// 本节介绍函数的定义知识

/*
*

1. go函数支持普通函数，匿名函数，闭包
2. go函数是一等公民--意思就是：（因为动态语言默认都是支持下面特性，但是静态语言有的不支持不算作一等公民）
  - 函数本身可以当作变量即可以作为参数传递给其他函数，也可以作为返回值返回
  - 匿名函数
  - 闭包

3. 函数的定义语法：

	func 函数名(参数列表) (返回值列表) {
		函数体
	}
	多值返回值：func add(a int, b int) (int, int) {
		return a + b, a - b
	}
	单值返回值：func add(a int, b int) int {}
	无返回值：func add(a int, b int) {} // 如果不定义返回参数类型，就代表无返回值
	返回值列表中还可以加变量名：
	func add(a int, b int) (sum int, diff int) {
	    // 也省去了var定义变量的过程，直接用即可
		sum = a + b
		diff = a - b
		return // 直接 return 就可以了，返回值会自动被识别为 sum 和 diff，
		// 写明也可以 return sum, diff //写不写都可以
	}

4. 返回值特点：go中函数是可以返回多值的，返回值列表可以省略，如果省略了，那么函数返回值就是void
其他语言函数只能返回单值，如js、python
*/
// 5. 函数参数传递：go函数中所有类型参数传递都是值传递，即函数内部对参数的修改不会影响到调用者
// 6. 变长参数的传递和接收：go语言支持变长参数，即函数可以接受任意数量的参数
// 变长参数可以是任意类型：string、float64、struct、interface{} 都可以。
// 变长参数的本质是：把传入的多个参数打包成一个切片
// func add(a ...int) int { // 变长参数，参数列表中可以省略参数类型
// a 的类型就是 []int
// 	sum := 0
// 	for _, v := range a {
// 		sum += v
// 	}
// 	return sum
// }

// 参数列表也可以写函数类型，因为函数是一等公民
func calu(opt string, myFunc func(int)) func(int) {
	return func(x int) {
		if opt == "double" {
			myFunc(x * 2)
		} else {
			myFunc(x)
		}
	}
}

func add(a int, b int) int { // 返回单值时，就不用加括号
	// 参数列表，若是相同类型的可以简写：func add(a, b int) int {}
	return a + b
}

func main() {
	fmt.Println(add(1, 2))
	// sum,_ := add(1, 2) // 多值时这样接受
	// fmt.Println(sum)

	// 匿名函数的使用：临时在参数中定义函数不能写名字
	calu("double", func(x int) {
		fmt.Println(x)
	})

	// 或者这样写匿名函数
	myFunc := func(x int) {
		fmt.Println(x)
	}
	calu("double", myFunc)

	// 闭包的使用：
	add := func(a int) func(int) int {
		return func(b int) int {
			return a + b
		}
	}
	add1 := add(1)       // 将1这个传给a给初始化闭包住了
	fmt.Println(add1(2)) // 输出 3
}
