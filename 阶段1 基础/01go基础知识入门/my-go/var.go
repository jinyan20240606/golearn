package mygo

// 静态语言的变量
// 1. 定义必须显式声明类型
// 2. 变量必须先定义后使用
// 3. 类型定下来后不能再改变，
// 4. 变量的零值：每种类型都有一个默认的零值，如 int 的零值是 0，string 的零值是 ""，bool 的零值是 false

// 1、全局变量和局部变量
//
//	func main() {
//		// 全局变量
//		var a = 1
//		// 常量
//		const a = 1
//		const b int = 2
//		const c int = a + b
//		// 局部变量
//		func() {
//			var b = 2
//			println(a, b)
//		}()
//	}

// 2、全局变量常量定义在main函数之外
const a = 1
const b int = 2
const c int = a + b

var name = "zhangsan"
var age int = 18
var height float64 = 1.8
var isOk bool = true

// 多个变量定义
var a2, b2, c2 int = 1, 2, 3 // 都是int类型
var a3, b3, c3 = 1, "3", 3

var (
	a1     = 1
	b1 int = 2
	c1 int = a1 + b1
)

func main() {
	// 共有哪几种数据类型：？
	// int 、float64、string、bool、byte、rune、interface{}、struct{}、map[string]int 、[]int、chan int、func() int、error、any、...
	// 3、定义变量的3种方式：
	// 方式1
	var a = 1 // 明确赋值时可以省略类型
	var b int = 2
	var c int = a + b

	// 方式2：简洁变量的声明方式，注意：不能用于全局
	d := 1

	// 方式3
	var e int
	e = 1

	// go 语言中，局部变量定义了且不使用的局部变量会报错，这是与其他语言的区别，不使用的全局变量不会报错
	println(c, d, e)
}

// 4、匿名变量
// 匿名变量的用途：函数返回多个值时，可以忽略掉不需要使用的变量
func getUserInfo() (string, int) {
	return "zhangsan", 18
}
func main2() {
	name, _ := getUserInfo() // 使用匿名变量 _ 来忽略掉不需要使用的变量
	println(name)
}
func main1() {
	var _ int = 1 // 匿名变量，表示这个变量不需要使用，编译器会忽略它的存在。匿名变量在函数参数列表中也很常见，用于表示某个参数不需要使用。
}

// 5、变量的作用域：全局和局部和块级作用域
//     - 变量的生命周期：全局变量在程序运行期间一直存在，局部变量在函数调用结束后就会被销毁
/**
 * 块级作用域
  “块”（Block）是指由 {} 包裹的代码区域，比如：
	函数体 func xxx() { ... }
	条件语句 if { ... }/switch { ... }
	循环语句 for { ... }
	甚至单独的 { ... } 包裹的代码段
	块级作用域：在某个块内声明的变量，仅在该块（及嵌套的子块）内可见，块外部无法访问。
*/
func main4() {
	var a1 = 1
	println(a1) // a1 是局部变量，只能在函数内部访问
	{
		// var a333 = 1
	}
	// println(a333) // a333 是块级变量，只能在它所在的块内访问，块外无法访问
}
