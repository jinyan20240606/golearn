package main

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
	fmt.Println(a33, b33) // 输出：65 20013，输出的是对应的数值，而不是字符本身，要输出字符本身需要使用 fmt.Printf("%c", a33) 或 fmt.Printf("%c", b33)或String(a33) 或 String(b33)
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
