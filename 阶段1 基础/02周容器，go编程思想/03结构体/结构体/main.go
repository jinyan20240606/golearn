package main

import "fmt"

// 在go语言中，通过结构体+interface 就能实现类似面向对象的复杂特性，go语言中没有面向对象的一些概念，但是能够实现类似面向对象特性

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
func (p Person) print1() {
	fmt.Printf("Person{name:%s, age:%d}", p.name, p.age)
	// return fmt.Sprintf("Person{name:%s, age:%d}", p.name, p.age)
}

func main() {
	// 结构体的不同初始化方式
	p1 := Person{name: "张三", age: 18}
	p2 := Person{"张三", 18, func() { fmt.Println("hello") }} // 结构体字面量，字段顺序必须和结构体定义的顺序一致

	var persons []Person
	persons = append(persons, p1, p2, Person{"张三", 18, nil})

	person2 := []Person{{"张三", 18, nil}, {"张三", 18, nil}} // 结构体字面量，字段顺序必须和结构体定义的顺序一致
	p2.add()
	fmt.Println(persons)
	fmt.Println(person2)

	// 直接赋值初始化
	var p Person
	p.name = "张三"
	// p.age = 18 // 不赋值默认就是0
	fmt.Println(p)

	// 匿名结构体：声明结构体的同时实例化
	var p3 = struct {
		name string
		age  int
	}{
		name: "张三",
		age:  18,
	}
	fmt.Println(p3)

	// 结构体的嵌套有2种方式：
	// 01
	type Person1 struct {
		p   Person
		age int
	}
	s := Person1{
		p:   Person{name: "张三", age: 18},
		age: 20,
	}
	s.p.age = 19
	fmt.Println(s.p.name, s.age)

	// 02 匿名字段嵌套：匿名字段就是没有名字的字段，直接写类型就行了，匿名字段的类型必须是唯一的，不能有重复的类型，否则编译器无法区分
	type Person2 struct {
		Person
		age int
	}
	s2 := Person2{ // 初始化时需要初始化匿名字段
		Person: Person{name: "张三", age: 18},
		age:    20, // 外部的设了同名字段，会覆盖掉内部的
	}
	// 赋值的时候：直接设置内部的字段即可，外部的同名字段会覆盖掉内部的
	s2.Person = Person{name: "李四", age: 19}
	s2.name = "李三"
	fmt.Println(s2.name, s2.age) // 张三 20

	// 调用结构体新绑定的方法
	var p14 = Person{name: "张三", age: 20}
	p14.print1()

	// 判断两个结构体的值是否相等
	// 需要结构体的每个字段的值都相等才相等，结构体的字段类型必须支持比较操作（==、!=），如果结构体中有一个字段的类型不支持比较操作，那么这个结构体类型也不支持比较操作，就无法使用==、!=来比较两个结构体的值是否相等了
	// fmt.Println(p1 == p2)
	// fmt.Println(p1 != p2)

}
