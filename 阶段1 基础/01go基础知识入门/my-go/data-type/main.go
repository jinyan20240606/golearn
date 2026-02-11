package main

/**
   数据类型参考文章：https://juejin.cn/post/7225535217339170871
   https://www.runoob.com/go/go-data-types.html
   Go 语言中的数据类型主要分为以下几类：
1. 基本数据类型：
	- 数值类型：整型【整数（int、int8、int16、int32、int64）、无符号整数（uint、uint8、uint16、uint32、uint64）】、浮点数（float32、float64）、复数（complex64、complex128）。
	- 布尔值：表示真假的值，通常用于条件判断。
	- 字符串类型：string 表示单个字符，通常使用单引号括起来。
    - 字符类型：byte、rune。
		- byte：Go 语言的语法规定（type byte = uint8），和 ASCII 码无关 —— 不管是否存储 ASCII 码，byte 本质都是 8 位无符号整数，取值范围 0~255；
			- ASCII 码是 byte 的典型场景：因为 ASCII 码的取值范围（0~127）刚好落在 byte（uint8）的取值范围内（0~255），所以 byte 天然适合存储 ASCII 字符，这是 “场景适配” 而非 “语法原因”
			- ASCII码主要用来存储英文的码点的，中文是无法表示的，所以 byte 不能存储中文字符。
		- rune等同于int32，用于表示UTF-8字符串的Unicode码点。
			- rune 本质上是 int32，表示 UTF-8 编码的字符。它的范围更大，专门表示「Unicode 码点」（全球所有语言的字符编码），包括中文、日文、Emoji 等无法用 byte（uint8）存储的字符


2. 复合数据类型：包括数组（array）、切片（slice）、结构体（struct）和映射（map）。
3. 引用数据类型：包括指针（pointer）和接口（interface）。
4. 函数类型：函数也是一种数据类型，可以作为变量传递和返回。
5. 通道类型：用于在 goroutine 之间进行通信的通道（channel）。

Go 语言还支持用户定义的类型，可以通过 type 关键字创建新的类型别名或结构体类型。
*/

import "fmt"

func main() {
	// 有符号整数
	var a int8 = 100            // 小范围 -128 ~ 127
	var b int64 = 1234567890123 // 大范围
	var c int = 1000            // 通用（64位系统下等价于int64）

	// 无符号整数
	var d uint8 = 255          // 字节最大值
	var e byte = 65            // byte是uint8别名，65对应ASCII的'A'
	fmt.Println(a, b, c, d, e) // 输出：100 1234567890123 1000 255 65

	// 浮点数

	var f1 float32 = 3.1415926         // 精度有限，实际存储为3.1415925
	var f2 float64 = 3.141592653589793 // 高精度
	var f3 float64 = 1.23e5            // 科学计数法，等价于123000.0

	fmt.Printf("float32: %.7f\n", f1)  // 输出：3.1415925（精度丢失）
	fmt.Printf("float64: %.15f\n", f2) // 输出：3.141592653589793（高精度）
	fmt.Println(f3)                    // 输出：123000
	// 注意事项
	// float64 是浮点型首选（精度更高，大部分场景无性能损失）；
	// 浮点型有精度丢失问题，金额计算禁止用 float（需用 decimal 第三方库或 int64 以分为单位）；
	// 不能直接用 == 比较两个浮点数是否相等（需判断差值是否小于极小值，如 1e-9）。

	// 字节类型：针对于字节单独起的别名的意思，相当于c语言中的char类型，本质就是个int类型
	var a33 byte = 'A'    // byte 是 uint8 的别名，'A' 的 ASCII 码是 65
	var b33 rune = '中'    // rune 是 int32 的别名，'中' 的 Unicode 码点是 20013
	fmt.Println(a33, b33) // 输出：65 20013，输出的是对应的数值，而不是字符本身，要输出字符本身需要使用 fmt.Printf("%c", a33) 或 fmt.Printf("%c", b33)
	var c33 byte
	c33 = 'a' + 1    // a的ASCII码加1，a的ASCII码为97
	fmt.Println(c33) // 输出：98,就代表b的ASCII码
}

// 类型转化
func main2() {
	var a int = 100
	var b float64

	b = float64(a) // int 转 float64
	fmt.Println(b) // 输出：100.0

	a = int(b)     // float64 转 int
	fmt.Println(a) // 输出：100
}
