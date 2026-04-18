package main

import (
	"container/list"
	"fmt"
)

// 讲解list链表的数据结构
/**
* slice和数组的缺点：
1. 数组的长度是固定的，不能动态改变，如果需要动态改变长度，可以使用切片（slice）来代替数组。
 2. slice底层一定是连续的存储空间，如果我们内存没有足够的连续存储空间的话，就无法分配对应大小的slice的
3. 链表是不一样的，链表的每个节点都是独立的内存空间，节点之间通过指针连接起来的，。
	- 优点：所以链表不需要连续的存储空间，可以动态改变长度，适合频繁插入和删除的场景
		- 数组和切片的插入得挪动两侧的存储空间，插入一个新的内存空间，两侧的空间都得移动
	- 缺点：链表需要额外的空间去保存指针，指针占4字节，所以链表需要额外的4字节空间去保存指针。
		- 插入和删除方便，但是查询性能差，需要遍历整个链表和指针指向才能找到目标元素。
		- 数组slice查询性能好，因为数组底层是连续的存储空间，可以通过索引快速找到目标元素。
	    	- 查询性能比slice差：slice的性能比链表更好，因为slice底层是数组，数组的访问速度更快，链表需要通过指针访问，需要不停的通过指针找到下一个元素，性能较差。
4. list是一个包，不是关键字不能直接用，需要先引包
*/
func main() {
	// 本节主要讲list：
	println("hello world")
	// 两种方式声明：
	// list := list.New()
	var list list.List
	list.PushFront("hello1") // 插入头部
	list.PushBack("hello")   // 插入尾部
	list.PushBack(1)
	fmt.Println(list)               // {{0x14000118120 0x14000118150 <nil> <nil>} 2}
	fmt.Println(list.Front().Value) // hello
	// 遍历打印链，不像forrange可以直接使用
	fmt.Println("开始正序遍历")
	// Front() 往前找返回链表的第一个元素，Back() 往后找返回链表的最后一个元素，Next() 返回链表中下一个元素，Prev() 返回链表中上一个元素
	for e := list.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
	fmt.Println("开始倒序遍历")
	for e := list.Back(); e != nil; e = e.Prev() {
		fmt.Println(e.Value)
	}
	// // 不能使用forrange
	// for _, v := range list { // 会报错：cannot range over list (type list.List)
	// 	fmt.Println(v)
	// }

	// InsertBefore(v any, mark *Element) *Element	在mark元素之前插入一个值为v的新元素，并返回新元素的指针。
	// 删除一个元素
	list.Remove(list.Front())

}
