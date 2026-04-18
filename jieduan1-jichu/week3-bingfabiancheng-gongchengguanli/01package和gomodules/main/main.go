package main

import ( // import组合导入，不用挨个写import
	"fmt"
	// 1. 相对路径导入（当前目录下），当存在go.mod 此方法就不可用
	// course "../user" // course是模块别名
	// 2. 模块路径导入（需先初始化 go.mod）
	// 注意：Go模块路径不支持中文，需将目录重命名为英文
	// 示例：go mod init golearn-project
	// 然后使用：import "golearn-project/week03/01package/user"
	// 其中 week03/01package/user 是 user 包相对于 go.mod 所在目录的路径
	course "golearn-01package-gomodules/user" // course是模块别名
	// . "golearn-01package-gomodules/user" 这种方式少用，可读性不好，通过点号打散，后面引用时，不用加前缀，方法名直接用
	// _ "golearn-01package-gomodules/user" // 这是匿名导入，不使用它里面的方法，只触发它模块内部的初始化方法
)

func main() {
	// 跨模块复用结构体
	c := course.Course{Name: "Go并发编程"}
	fmt.Println(c.GetName(), course.GetCourse())
}
