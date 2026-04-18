package main

import (
	"fmt"
	// 1. 相对路径导入（当前目录下），当存在go.mod 此方法就不可用
	// course "../user" // course是模块别名
	// 2. 模块路径导入（需先初始化 go.mod）
	course "golearn-01package-gomodules/user"
	// 注意：Go模块路径不支持中文，需将目录重命名为英文
	// 示例：go mod init golearn-project
	// 然后使用：import "golearn-project/week03/01package/user"
	// 其中 week03/01package/user 是 user 包相对于 go.mod 所在目录的路径
	// 3. 本地包替换导入（推荐生产环境）
	// 在 go.mod 中使用 replace 指令将远程包替换为本地包：
	// go.mod 内容：
	//   module golearn-project
	//   go 1.21
	//   replace github.com/yourname/user => ../user
	// 然后使用：import "github.com/yourname/user"
	// 这样代码中用远程路径，实际指向本地目录
)

func main() {
	c := course.Course{Name: "Go并发编程"}
	fmt.Println(c.GetName())
}
