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
}
