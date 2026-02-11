package main

// 字符转化的标准内置包
import "strconv"

// 类型的转换
func main() {
	var a int = 10
	var a1 = uint8(a) // int 转 uint8，需要显式转换
	var b float64 = 5.5
	var b1 = int64(b) // float64 转 int64，会舍弃小数部分

	// 将 int 类型转换为 float64 类型
	var c float64 = float64(a) + b
	println(c, a1, b1)
	println(a + int(b))

	// 类型别名
	type myInt int32
	// 赋值类型别名时，必须用别名转，或者具体对应的值，如var a2 myInt = int32(10)就报错
	// var a2 myInt = myInt(10) 或者下面的赋值
	var a2 myInt = 10
	var b2 float64 = float64(a2)
	println(b2)

	// 字符串转数字
	var str = "123"
	// Itoa()函数将数字转换为字符串，I代表init，a代表char
	// Itoa该函数只返回一个参数，为转换后的数字，不可能出一异常的
	// Atoi该函数返回2个参数，1为转换后的数字，2为错误信息，当错误信息不为空时，说明转换失败出异常了
	var num, err = strconv.Atoi(str)
	if err != nil {
		println("转换失败")
	} else {
		println(num)
	}

	// 字符串转换为float类型，转化为bool类型

}

// 类型转换注意事项：
// 1. 显式转换：Go 语言中的类型转换必须是显式的，不能隐式转换。例如，不能直接将 int 类型赋值给 float64 类型，必须使用转换函数。
// 2. 精度损失：在将浮点数转换为整数时，可能会发生精度损失，因为小数部分会被截断。
// 3. 不同类型之间的转换：某些类型之间的转换可能会导致数据丢失或溢出，例如将较大的整数类型转换为较小的整数类型时。
