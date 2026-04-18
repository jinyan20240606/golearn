package main

import (
	"fmt"
	"time"
)

// go中，只有for循环，没有while循环和do while循环，for循环可以实现这两种循环的功能
func main() {
	print("讲解 for循环")
	var i1 int = 5
	time.Sleep(2 * time.Second)
	// 初始化变量语句；条件判断；后置语句
	for i := 0; i < 5; i++ {
		println(i)
	}
	// for啥都不写，就相当于while(true) 了
	// for {
	// 	println("hello world")
	// }

	for i1 > 3 { // 模拟的while循环和跳出循环
		println(i1)
		i1--

	}

	// 用for循环求1-100的累加和
	var sum int = 0
	for i := 1; i <= 100; i++ {
		sum += i
	}
	println(sum) // 5050

	// for循环实现9乘法表
	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			print(j, "*", i, "=", i*j, "\t") // \t表示制表符Tab
		}
		println()
	}

	// for循环还有一种用法：for range 遍历数组、切片、map、字符串，channel，相当于js中的forEach
	/*
		// 语法
		for key, value := range arr {
			// value是对应的值的拷贝，修改value不会修改原来的值
			println(key, value)
		}
	*/
	name := "hello world"
	for i, v := range name {
		println(i, v, string(v)) // v是字符的ASCII码（英文对应ASCII码，中文对应unicode码点） ---要想打印字符本身     println(i, string(v)) // v是字符
		fmt.Printf("%c\n", v)    // 也可以使用fmt.Printf()函数，%c表示按字符格式输出
	}
	// key不想使用时，可以使用_占位符来忽略它，要不未使用定义的变量会语法报错
	for _, v := range name {
		println(v, string(v)) // v是字符的ASCII码（英文对应ASCII码，中文对应unicode码点） ---要想打印字符本身     println(string(v)) // v是字符
	}
	// 只写一个值时，默认就代表下标
	name1 := "hello 世界"
	name1Rune := []rune(name1)
	println("测试62----")
	for v := range name1 {
		fmt.Println(v, string(name1[v]))
		// v是下标，name[v]是对应的字符,当name含中文时，index是字节的索引，直接用name[index] 展示中文是会乱码的，
		// 你只能用默认的for的key，value形式获取中文字符
		// 或者 转成rune数组进行遍历，用rune[index]获取中文字符
	}
	for i := 0; i < len(name1Rune); i++ {
		fmt.Println(i, string(name1Rune[i]))
	}

	// for循环的退出：break 和 continue
	// break：跳出当前循环，继续执行循环后面的代码
	for i := 0; i < 10; i++ {
		if i == 5 {
			break // 跳出循环
		}
		println(i)
	}
	println("循环结束")

	// continue：跳过当前循环的剩余代码，直接进入下一次循环
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue // 跳过偶数，继续下一次循环
		}
		println(i) // 打印奇数
	}

	// goto语句 跳出嵌套的循环
	// goto语句可以直接跳转到指定的标签处，标签是一个标识符，后面跟一个冒号，goto语句会跳转到这个标签处继续执行代码，一般不常用，容易导致代码混乱，不建议使用，，，
	// 一般用于：当程序出现错误时，使用goto语句统一跳转到错误处理的代码处进行处理，可以避免代码重复，提高代码的可读性和维护性
	for i := 0; i < 10; i++ {
		println("i:", i)
		for j := 0; j < 10; j++ {
			if i == 5 && j == 5 {
				// 此处使用break只能跳出内层循环，外层循环仍然会继续执行
				goto over // 直接跳出所有循环，跳到over标签处继续执行
			}
			println(i, j)
		}
	}
over:
	println("循环结束")
}
