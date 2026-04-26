package main

import "fmt"

// 在go语言中，通过结构体+interface 就能实现类似面向对象的复杂特性，go语言中没有面向对象的一些概念，但是能够实现类似面向对象特性

// 结构体的定义中，花括号里面写的都是字段和类型（不能写方法、不能写逻辑、不能写函数体），外面用接收器绑定具体方法
type Person struct {
	name string
	age  int
	add  func() // 结构体中也可以定义函数类型的字段，函数类型也是一种类型，可以当作字段的类型
}

// 如何为结构体绑定个不同的方法呢？使用接收器方式：func (接收者 StructType) 方法名(参数列表) 返回值类型 { 方法体 }
// 这个接收器：有2种形态
//   - 值接收器：func (p Person) print() {} // 这种方式是值传递，方法内部对接收者的修改不会影响到调用者，默认的print方法就是值接收器，方法内部对接收者的修改不会影响到调用者
//     ///////// - 值传递缺点：在方法内部如修改了p.age后，但是调用者的p.age不会改变，因为是值传递｜|还有大对象下不断拷贝性能问题
//     ///////// - 适用场景：当方法内部不需要修改接收者的字段时，或者接收者是一个小结构体时，使用值接收器比较合适，因为值传递会有一定的性能开销，如果结构体比较大不断进行拷贝，使用值传递会有较大的性能开销
//     ///////// - 如果值传递的是一个指针类型的字段，那么在方法内部修改这个指针类型的字段会影响到调用者，因为指针类型的字段是引用传递的（p1 := &Person{name: "张三", age: 18}）
//   - 指针接收器：func (p *Person) print() {} // 这种方式是指针传递，方法内部对接收者的修改会影响到调用者
//
// ---1不能定义相同方法名称的不同接收器的方法，否则编译器无法区分，无法编译通过
// ---2！！！！接收器必须定义外层包级别，不能定义在函数内部！！！！！
// ---3 接收器中p的命名写法一版规范是结构体的首字母小写，两个单词是两个首字母小写

// 结构体的指针和值类型的不同初始化方式:
// var mrw MyReadWriter = SreadWriter{}// 这个是值类型的结构体初始化
// 一般情况下我们都会用指针初始化，除非明确你的逻辑适合拷贝方式的值类型才能用
// var mrw MyReadWriter = &SreadWriter{} // 这个是指针类型的结构体初始化
// 值类型的初始化时，方法接收器必须写值类型
// 指针类型的初始化时，方法接收器可以写指针类型或者值类型【因为都是兼容的】
// 方法接收器写指针类型时，结构体初始化必须用指针类型初始化

func (p Person) print1(p1 []byte) {
	fmt.Print(p1)
	// return fmt.Sprintf("Person{name:%s, age:%d}", p.name, p.age)
}

func main() {
	// 结构体的不同初始化方式
	// 01
	p1 := Person{name: "张三", age: 18}
	// 02
	p2 := Person{"张三", 18, func() { fmt.Println("hello") }} // 结构体字面量，字段顺序必须和结构体定义的顺序一致

	var persons []Person
	persons = append(persons, p1, p2, Person{"张三", 18, nil})

	person2 := []Person{{"张三", 18, nil}, {"张三", 18, nil}} // 结构体字面量，字段顺序必须和结构体定义的顺序一致
	p2.add()
	fmt.Println(persons)
	fmt.Println(person2)

	// 03 还有一种结构体的初始化方式：什么都不填
	// 1. 初始化时不用赋值，不填就是它的零值
	type MyWriter interface {
		print1(p []byte)
	}

	type writerCloser struct {
		Person       // 匿名结构体字段
		age      int // 普通字段
		MyWriter     // 匿名接口类型字段
	}
	// var wc writerCloser             // ✅ 完全可以
	wc := writerCloser{
		age: 18,
		// Person: Person{name: "张三", age: 18},
		MyWriter: Person{name: "张三", age: 18}, // 用结构体实现接口，不实现接口，就默认是零值
	} // ✅ 完全可以
	// 或者
	fmt.Println(wc.MyWriter, wc.MyWriter == nil)                          // <nil> true ✅ 接口的零值是nil，而下面结构体的零值是每个字段的零值
	fmt.Println("-----58----", wc.age, wc.Person, wc.name, wc.Person.age) // 18 { 0 <nil>}  0
	// Go 中属于 void 类型。fmt.Println 要求传入具体的值参数，无法接受“无值”作为参数
	wc.MyWriter.print1([]byte{2}) // ✅ 由于print1方法无返回值，所以不能放在Println中打印这个函数结果会报错
	// wc.print1([]byte{2})          // 和上行一样，匿名会自动提升

	// 04 直接赋值初始化
	var p Person
	p.name = "张三"
	// p.age = 18 // 不赋值默认就是0
	fmt.Println(p)

	// 05 匿名结构体：声明结构体的同时实例化
	var p3 = struct {
		name string
		age  int
	}{
		name: "张三",
		age:  18,
	}
	fmt.Println(p3)

	// 结构体的嵌套有2种方式：
	// 01 有名字段嵌套
	type Person1 struct {
		p   Person
		age int
	}
	s := Person1{
		p:   Person{name: "张三", age: 18},
		age: 20,
	}
	s.p.age = 19
	fmt.Println(s.p.name, s.age) // 张三 20 (读值时外部同名会覆盖嵌套的内部同名的字段)

	// 02 匿名字段嵌套：匿名字段就是默认字段名与类型名相同的一种简写方式，直接写类型就行了，匿名字段的类型必须是唯一的，不能有重复的类型，否则编译器无法区分
	//////// 注意：
	//  - 只要是【命名类型】或【命名类型的指针】都能作为结构体匿名字段
	//       - 命名类型：自己定义的、有名字的类型 如：type Person struct {   Name string} 或 type MyInt int        // 命名类型
	//       - 命名结构体，命名接口都可以
	//  - 作用：如结构体嵌入接口 = 强制这个结构体必须实现该接口
	type MyInt int
	type Person2 struct {
		Person // 匿名就是默认字段与类型相同，后面可以访问时直接访问person2.Person.xxx
		// *Person // 合法：命名类型的指针（匿名）
		age int
		MyInt
		// int         // 虽然能跑，但不推荐
		// string      // 虽然能跑，但不推荐
		// []int       // 错误！切片不能匿名
		// map[string]int // 错误！
		// func()       // 错误！
		// interface{}  // 错误！
		// chan int     // 错误！
	}
	// 初始化01：时需要初始化匿名字段
	s2 := Person2{
		Person: Person{name: "张三", age: 18},
		age:    20, // 外部的设了同名字段，会覆盖掉内部的
	}
	// 初始化02: Go 仅对「单一匿名字段」的简写语法：匿名嵌入类型 MyWriter，简写初始化可以直接写值，不用写字段名
	// 当结构体里多个字段 / 多个匿名嵌入，初始化：必须写 类型名：值 显式指定，不写就编译报错
	//         - 如见05接口/multi/main-struct.go中的写法
	// var mw MyWriter = writerCloser{
	// 	fileWriter{},
	// }
	// 赋值的时候：直接设置内部的字段即可，外部的同名字段会覆盖掉内部的
	// 规则 1：匿名嵌入 → 内部字段自动提升
	// 规则 2：同名字段，外层优先，遮挡内层：s2.age 👉 永远优先找 Person2 自己的 age，想访问内部被遮挡的 age，必须写全名
	//        - 如修改s2.age 其实是修改外层的age，不是Person里，要想修改内部，只能显示修改：s2.Person.age = 19
	s2.Person = Person{name: "李四", age: 19}
	s2.name = "李三"
	fmt.Println(s2.name, s2.age) // 李三 20

	// 有名字段和匿名字段的区别：
	// 01 ： 调用方式不同
	//   - 有名字段：必须写 .p.xxx
	//   - 匿名字段：会提升可以直接 .xxx（像自己的一样，是 Go 的 “组合式继承”，这是唯一的go中实现继承的方式），
	//        - 当然直接用.p.xxx也是可以的，只是多了一种继承的用法
	// 02 ：方法 / 字段提升
	//  - 匿名字段 = 自动提升：内部结构体的 所有字段 → 提升到外层，内部结构体的 所有方法 → 提升到外层
	//        - 提升后，但是原始方式访问也可以
	//  - 有名字段 = 不提升，通过字段访问
	// 03 : 初始化写法不同
	//// p1 := Person1{
	////  // 有名字段（正常写）
	//// 	p:   Person{Name: "小明"},
	//// 	age: 18,
	//// }
	//// p1 := Person1{
	////  // 匿名字段（必须写类型名），
	//// 	Person: Person{Name: "小明"},
	////  // 当仅有一个匿名字段时，可以简写直接写值不用写字段如下：
	////  // Person{Name: "小明"},
	//// 	age:    18,
	//// }

	// 调用结构体新绑定的方法
	var p14 = Person{name: "张三", age: 20}
	p14.print1([]byte("hello")) // [104 101 108 108 111]%

	// 判断两个结构体的值是否相等
	// 需要结构体的每个字段的值都相等才相等，结构体的字段类型必须支持比较操作（==、!=），如果结构体中有一个字段的类型不支持比较操作，那么这个结构体类型也不支持比较操作，就无法使用==、!=来比较两个结构体的值是否相等了
	// fmt.Println(p1 == p2)
	// fmt.Println(p1 != p2)

}
