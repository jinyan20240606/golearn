package main

import "fmt"

// 指针的用法和定义
/**
 * 指针就是存 “内存地址” 的变量，它指向真正的值。
	普通变量：存值
	指针变量：存地址

1. 在 Go 里，所有数据类型都有对应的指针类型。int → *int、string → *string、struct → *struct、map → *map、func → *func 全都有。
	只要是 Go 里的值类型，就一定能：用 &x 取地址，得到对应的指针类型 *T
	--- 写法统一：T 对应 *T
2. 只有一个例外：nil 本身没有指针，nil 是一个标志，不是类型，所以没有 *nil。
3. 容易混淆的点：引用类型本身已经 “像指针”，，，像 slice、map、chan、func、interface它们底层就是指针结构，所以平时很少写 *[]int、*map，但语法上依然存在。，var s []int
var p *[]int = &s  // 完全合法


符号规则：
	*数据类型：定义指针数据类型如*int 是 int 类型指针，定义各个数据类的指针
	&变量 = 取地址
	*指针 = 取值
*/
func println(a *int) { // 参数中接收一个 int 类型指针
	fmt.Println(*a) // 函数内 *a：取出指针指向的真实值
}
func change(num *int) {
	*num = 100
}

type Person struct {
	name string
}

// 4-3：通过swap交换指针的位置

// 1. 函数接收指针参数
// func swap1(p1, p2 *Person) { // 参数连续写多个变量，只写最后一个类型，前面所有变量都是这个类型
// 	var temp Person
// 	temp = *p1
// 	*p1 = *p2
// 	*p2 = temp
// }

// 2. 接收器接收指针参数 .接收器也必须声明在包级别，不支持在函数内部定义方法
func (p1 *Person) swap2(num int) int {

	return num
}

// 4-3 通过指针实现交换2个值
func swap(a, b *int) {
	// a, b = b, a // 【实现方式A】不对，因为函数参数都是值传递，指针也是值传递相当于拷贝了新的指针，新副本指针交换交换，所以不生效
	*a, *b = *b, *a // 【实现方式B】正确，a地址的值改成b地址的值，b地址的值改成a地址的值，这样就交换了a和b指向的值
}

func main() {

	// 指针的不同初始化方式
	// var a *int // 声明一个 int 指针（此时=nil，空指针）
	// 第一种初始化方式：
	ps := &Person{name: "里斯"} // 直接用 & 取地址，初始化一个结构体指针

	// 第二种初始化方式：
	var ps1 Person
	a1ww2 := &ps1 // 把 b 的地址赋值给指针 a
	// Go 会自动帮你解引用指针，所以指针可以直接用 . 访问字段，不用手写 *
	fmt.Println(a1ww2.name, ps) // 这不会报空指针错误：string 类型的零值 = 空字符串 ""

	// // ❌ 这才会报空指针错误！
	// var a1ww2 *Person // 只声明指针，没指向任何变量（= nil）
	// fmt.Println(a1ww2.name)
	// // panic: runtime error: invalid memory address or nil pointer dereference

	// 第三种初始化方式：
	// new(结构体) → 返回指针 *Person
	// 分配一块内存，存放空结构体
	// 结构体里的 string 字段默认值 = 空字符串 ""
	// 所以 println(a2.name) → 输出空（什么都不显示）,不会报错
	var a2 = new(Person) //
	print(a2.name)       // 输出 0

	// 初始化的2个关键字：map，channel，slice初始化推荐使用make方法
	// 指针初始化推荐使用new函数，指针要初始化，否则会出现nil pointer error:
	// 指针零值是 nil，直接访问会报空指针错误
	// map必须要初始化

	///。
	// 指针最核心 4 个用法
	// 1. 取地址 &
	x := 10
	p := &x // p 是指针，存 x 的地址
	// 2. 解引用 *
	fmt.Println(*p) // 输出 10
	// 3. 修改指针指向的值
	*p = 20
	fmt.Println(x) // x 变成 20
	// 4. 指针做函数参数（可以修改外部变量）
	a1 := 10
	change(&a1)
	fmt.Println(a1) // 100

	// 绕圈符号：* 和 & 在一起会互相抵消，无论套多少层 &*&*&*，最后都 = c
	var c *int = new(int)
	println(*&c)
	println(&*c)
	println(&*&*c)
	println(&*&*&*c)
	// go中指针与c++2处不同的地方：
	// 1. go中没有指针运算，不能对指针进行加减运算，而c语言是可以的
	//        - go中没有指针类型转换，不能将一个指针类型转换为另一个指针类型，不能将一个指针类型转换为一个整数类型，不能将一个整数类型转换为一个指针类型
	// 2. go语言中对结构体指针的访问会自动解引用
	p0 := &Person{name: "里斯"}

	// p0.name = "33"
	(*p0).name = "33" // 这两种方式是等价的，go语言中对结构体指针的访问会自动解引用，所以可以直接使用点操作符访问结构体字段，而不需要显式地解引用指针)

	// go的指针是一个阉割版，unsafe包里面才有指针运算和指针类型转换等功能，普通的go代码中是没有的，----不是特殊情况一般不使用unsafe包
	// 这也是go语言设计者为了安全性而做出的设计选择，避免了指针运算和指针类型转换带来的安全风险，同时也简化了go语言的使用，让开发者更专注于业务逻辑的实现，而不是指针的操作细节。

	var p1 *Person
	// fmt.Println(p1.name) // 不行会报错，因为p1是nil指针，不能访问它的字段

	p1.swap2(2)

	// 4-3
	a, b := 1, 2
	swap(&a, &b)
	fmt.Println(a, b) //实现方式B正确：2 1 ----- 如果用的实现方式A：结果还是a=1 b=2，没变
}
