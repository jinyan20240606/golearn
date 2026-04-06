package main

import (
	"fmt"
)

// 讲解map的数据结构
/**
 * map是一个key-value的无序集合，主要是查询方便
 */
func main() {
	// 本节主要讲map：很常用 TODO 1-11
	println("hello world")
	// 定义 3种方式：
	var map1 map[string]int     // 定义并没有初始化
	map1 = make(map[string]int) // 使用make函数初始化，make是内置函数，可以初始化map、slice、channel等引用类型的数据结构，make函数会分配并初始化内存空间，并返回一个对应类型的引用。对于map来说，make函数会创建一个空的map，并为其分配内存空间，使其可以存储键值对。使用make函数初始化map后，就可以向其中添加键值对了。
	var map2 = map[string]int{  // 使用字面量初始化，写空花括号也算初始化
		"hello": 1,
		"world": 2,
	}
	// 记住：map的key类型哪些类型可以作为key，哪些不能作为key，绝大多数都可以做value
	fmt.Println(map2)
	// 取值放值
	map1["hello"] = 1
	println(map1["hello"])
	// 删除一个元素，删除一个不存在的元素也不会报错的
	delete(map1, "hello")
	fmt.Println(map1["hello"])
	// map的遍历
	// map是无序的，不保证每次遍历打印都是相同的顺序，slice和数组是有序的
	for k, v := range map1 {
		fmt.Println(k, v)
	}

	// go语言中有一个空的类型叫：nil，（类似于js的null）nil可以用来表示map、slice、channel等引用类型的零值，表示它们没有被初始化或者没有分配内存空间。当你声明一个map但没有使用make函数初始化它时，这个map的值就是nil。对于一个nil map，你不能直接向它添加键值对，因为它没有底层的数据结构来存储这些数据。如果你尝试向一个nil map添加键值对，会导致运行时错误（panic）。因此，在使用map之前，必须先使用make函数或字面量语法来初始化它，以确保它不是nil。
	var map3 map[string]int // 这是定义map3类型，但没有初始化，所以它的值是nil，加个空花括号也算初始化加值也不报错
	fmt.Println(map3)       // 输出map[]，表示map3是一个nil map
	// map3["hello"] = 1 // 这行代码会导致运行时错误，因为map3默认是nil，无法添加键值对

	// map必须初始化才能用，1.map[string]int{} 2.make(map[string]int),但是slice、channel、array、struct这些可以不初始化就可以用，默认值为nil或者空值
	var m1 []string
	fmt.Println(m1 == nil) // 打印为m1为[]，为true
	m1 = append(m1, "hello")

	// 判断map中是否有某个key：map[key] => 总共返回2个参数，第一个是值(不存在的key的值也会返回"")，第二个是bool值，表示key是否存在
	if _, ok := map1["hello"]; ok { // if内可以写逻辑代码用分号隔开
		fmt.Println("hello exists")
	} else {
		fmt.Println("hello does not exist")
	}

	// 很重要的提示：map不是线程安全的，并发编程的时候同时对一个map操作是要报错的，如果你要操作的话
	// 得使用sync.Map（sync包下的map）
	// var m sync.Map
}
