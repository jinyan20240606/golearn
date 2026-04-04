package main

import "fmt"

func main() {
	// switch 语句的语法
	/*
		 switch 表达式 {
		 case 常量1:
			 // 执行逻辑1
		 case 常量2:
			 // 执行逻辑2
		 default:
			 // 默认执行逻辑
		 }
	*/
	var i = 1
	switch i {
	case 1, 3, 5: // case 语句可以同时匹配多个值，多个值之间用逗号分隔
		fmt.Println("i = 1")
	case 2:
		fmt.Println("i = 2")
	default:
		fmt.Println("i > 2")
	}
}
