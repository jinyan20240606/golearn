package main

import "fmt"

/*
*

	if 布尔表达式 { // 条件这块if (表达式) {} 加不加括号都可以，默认规范不加括号
	   // 执行逻辑
	}
*/
func main() {
	println("if 条件判断")
	age := 18
	country := "中国"
	if age >= 18 && country == "中国" { // 这里是一个复合条件，多个条件之间用 && 连接
		// 或者 if (age >= 18) && (country == "中国") {  .... }
		fmt.Println("成年了")
	} else if age < 18 {
		fmt.Println("未成年")
	} else {
		fmt.Println("未成年")
	}

	// if内可以写逻辑代码用分号隔开
	if age := 18; age >= 18 {
		fmt.Println("成年了")
	} else {
		fmt.Println("未成年")
	} // if 里面可以定义变量，但是只能在if里面使用
	// if 里面定义的变量，只能在if里面使用，不能在if外面使用
}
