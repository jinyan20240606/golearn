package main

import (
	"fmt"
	"strings"
)

// 字符串的基本操作
func main() {
	println("hello world")

	name := "张三"
	bytes := []rune(name)
	bytes[0] = '李'
	name = string(bytes)
	println(name, len(bytes)) // 输出：李三 2 每个中文占1个rune位置

	// 转义符
	println("hello\"你好\"world")

	// 反引号转义: fmt.Println(`hello"你好"world`)
	println(`hello"w"orld`)
	println("hello\rworld")

	// 格式化输出: 字符串拼接
	fmt.Printf("hello %s %d %.2f结束", "world", 123, 3.1415)             // 输出：hello world 123 3.14结束%
	userMsg := fmt.Sprintf("hello %s %d %.2f结束", "world", 123, 3.1415) // 输出：hello world 123 3.14结束%
	fmt.Println(userMsg)

	// 通过string的builder 来拼接字符串 （高性能的方式，但体验最差）
	var builder strings.Builder
	builder.WriteString("hello")
	builder.WriteString("world")
	fmt.Println(builder.String())

	// 字符串的操作的常用方法
	// 是否包含
	fmt.Println(strings.Contains("hello world", "world")) // true
	// 字符串的长度：len() 是 Go 内置函数，对 string 类型调用时，返回的是字符串底层 UTF-8 编码的字节数，而非字符个数
	/**
	* 不同长度方法关键区别：len() 对字符串的计算逻辑和对切片（如 []rune/[]byte）的计算逻辑不同：
		对 []rune("hello你world") 调用 len()，返回字符个数（5+1+5=11）；
		对 string 调用 len()，返回字节个数（核心语法规则）。

		若需获取字符个数（而非字节数），需借助 []rune 转换后调用 len()（核心语法联动）：
		```go
		// 语法：len([]rune(字符串)) → 字符个数
		fmt.Println(len([]rune("hello你world"))) // 输出 11（5+1+5）
		```
	*/
	fmt.Println(len("hello你world")) // 13 每个中文占3个字节
	// 字符串的切割
	fmt.Println(strings.Split("hello world", " ")) // [hello world]
	// 字符串的替换
	fmt.Println(strings.ReplaceAll("hello world", "world", "go")) // hello go
	// 大小写转换
	fmt.Println(strings.ToUpper("hello world")) // HELLO WORLD
	fmt.Println(strings.ToLower("HELLO WORLD")) // hello world
	// 字符串的前缀和后缀
	fmt.Println(strings.HasPrefix("hello world", "hello")) // true
	fmt.Println(strings.HasSuffix("hello world", "world")) // true
	// 字符串的出现次数
	fmt.Println(strings.Count("hello world", "l")) // 3
	// 字符串的索引, 也是按字节的索引计算的，中文占3个字节的索引长度
	fmt.Println(strings.Index("hello你world", "w"))     // 8
	fmt.Println(strings.LastIndex("hello world", "o")) // 7
	// 去掉左右两边的指定的特殊字符
	fmt.Println(strings.Trim("hello world", " ")) // hello world
	// 去掉左右两边的特殊字符
	fmt.Println(strings.TrimLeft("hello world", " "))  // hello world
	fmt.Println(strings.TrimRight("hello world", " ")) // hello world
	// 去掉左右两边的特殊字符，且指定字符
	fmt.Println(strings.Trim("hello world", "h"))   // ello world
	fmt.Println(strings.Trim("hello world", "d"))   // hello worl
	fmt.Println(strings.Trim("hello world", "ld"))  // hello wor
	fmt.Println(strings.Trim("hello world", "hld")) // ello wor
}
