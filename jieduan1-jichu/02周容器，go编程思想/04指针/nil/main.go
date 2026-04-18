package main

import "fmt"

// nil 在go中的细节
// nil代表某些类型的数据零值
/**
  不同类型的数据的零值不一样
* bool的零值：false
  * int的零值：0
* float的零值：0.0
* string的零值：""
* map的零值：nil，slice的零值：nil， channel的零值：nil， interface的零值：nil，function 也是nil
* pointer的零值：nil // 指针的零值是nil，直接访问会报空指针错误
* struct的零值：{。。。} // 结构体的零值是一个所有字段都为零值的实例（
*/

func main() {
	var a int = 1
	var b int = 2
	var c int = a + b
	println(c)

	type Person struct {
	}

	var ps []Person // 一般叫做nil的slice
	if ps == nil {  // slice的零值是nil，所以条件成立
		fmt.Println("ps is nil")
	}
	/**
	// 切片在go底层是一个结构体：
	type slice struct {
		ptr *数组元素   // 指向底层数组的指针
		len int        // 长度
		cap int        // 容量
	}

	 * 第一种：var ps []Person ===> 这叫 nil 切片（nil slice）
	  - 底层发生了什么？
		只声明，不分配底层数组
		切片结构体里的 ptr = nil，len=0，cap=0
		整个切片结构体等于零值 nil
	*/

	// 这个一般叫做：空的slice
	/*
		第二种：var ps2 = make([]Person, 0)
		底层发生了什么？
			make 会分配一个底层数组（哪怕长度 0）
			切片结构体：
			ptr ≠ nil（指向一个空数组）
			len=0
			cap=0
			切片本身不是 nil！

		**/

	var ps2 = make([]Person, 0) // 使用make函数初始化一个空的slice，ps2不是nil了，而是一个长度为0的slice
	if ps2 == nil {             // 条件不成立了，因为ps2不是nil了，而是一个长度为0的slice
		fmt.Println("ps2 is nil")
	} else {
		fmt.Println("ps2 is not nil")
	}

	// 同样的：map的也有 nil的map 和 空的map
	// 	map 在 Go 底层也是个指针，所以：
	// nil map = 指针为 nil
	// empty map = 指针不为 nil，指向一个空的 hmap 结构体

	// 2. nil map

	var m map[string]int
	// 底层
	// 只声明，没有调用 makemap
	// 变量 m 本质就是个空指针 = nil
	// 行为
	// 可以读、可以遍历（不会崩）
	fmt.Println(m["a"]) // ✅ 读一个不存在的键，返回值类型的零值（如int 的零值 = 0）
	// 不能写！

	m["a"] = 1 // ❌ panic: assignment to entry in nil map

	// 3. empty map（空 map）
	m := make(map[string]int, 0)
	// 或
	m := map[string]int{}
	// 底层
	// 调用了 makemap
	// 分配了 hmap 结构体
	// 指针不是 nil，只是里面没有键值对
	// 行为
	// 可读、可写、可遍历

	m["a"] = 1 // ✅ 完全正常

	// 所以为了安全起见，尽量使用 make去初始化

}
