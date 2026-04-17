package main

// / 接口数据类型

// 1. 接口的定义语法 = 方法的集合（只定义要做什么，不定义怎么做）
//    - 只写方法签名，不写实现，不写字段

// type 接口名 interface {
//     方法名1(参数) 返回值
//     方法名2(参数) 返回值
//     // ...
// }
// 2. 任何类型，只要拥有接口里的所有方法，就自动实现了这个接口。，不需要写 implements！不需要继承！
// 3. 只要结构体中包含即可不是全等匹配，就可以满足接口类型的识别（如接口只需要一个say方法，结构体中包含一个say方法还有其他不相干方法，那么这个结构体就可以实现MyWriter1接口）

// 5-1go 语言的接口，鸭子类型 -- 见readme.md记录

// 5-2接口的定义：接口只是定义方法的声明，不是具体的实现
type Duck interface {
	Quack()
	Fly()
}

// 结构体实现接口，很简单：只要实现了接口中声明的所有方法，就算实现了接口，不需要显式声明实现了哪个接口
type Person struct {
	Name string
}

func (p Person) Quack() {
	println(p.Name + " is quacking")
}
func (p Person) Fly() {
	println(p.Name + " is flying")
}

// 实现接口的检验标准：看能不能把一个接口类型的变量赋值为一个结构体类型的变量，如果能赋值成功，那么这个结构体就实现了这个接口

func main() {
	var duck Duck
	duck = Person{Name: "John"}
	duck.Quack()
	duck.Fly()
}
