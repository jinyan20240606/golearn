package main

import "fmt"

// defer的应用场景：类似于trycatch中的finally，省去了finally的冗长写法，距离try太远，容易忘写，
// defer可以直接写函数内的主代码挨着，会兜底最后且return前去延迟执行，可以利用这个机制做一些清理工作

func main() {
	// 链接数据库，打开文件，开始锁，无论中间成功或失败如何最后都要记得去关闭数据库，关闭文件，解锁
	// 用defer挨着写在一起，最后执行，避免忘记写，保证了无论中间代码成功失败都能执行到defer的代码
	defer fmt.Println("关闭数据库")
	defer fmt.Println("关闭文件")
	defer func() {
		fmt.Println("关闭锁")
	}()
	fmt.Println("开始处理数据")

	// 多个defer的执行顺序：栈的顺序，先进后出，后进先出
	// 输出结果：
	// 开始处理数据
	// 关闭锁
	// 关闭文件
	// 关闭数据库

	// defer也会在你函数return前，修改你变量的值的
	var deferReturn = func(ret int) int {
		defer func() {
			ret++
		}()
		return ret
	}
	fmt.Println(deferReturn(1)) // 输出结果：2, 因为defer会修改你返回值的值
	return

}
