package main

// 这个文件想要运行：你的包名必须是 main，并且还需要一个 main 函数。为什么？因为 Go 程序从 main 函数开始执行。
// 你可以把 main包的main 函数想象成程序的入口点。
// 【重要：】目录与包的关系： 一个目录 = 一个包，目录下所有 .go 文件必须声明相同包名（如 package main）（与python不同）
//         - Go 要求「同一个目录下的所有 .go 文件必须属于同一个包」，尤其是可执行程序（需要 package main + main() 函数），这个规则是硬性的
//         - 同一个目录下只能有一个 main() 函数（多个会冲突）。

// 这是一个简单的 Go 程序，它打印 "Hello, World!" 到控制台。

// import "fmt" 引入 fmt 包，这个包提供了一些用于格式化输出的函数。
// println的ln代表“line”，表示打印一行文本并换行。
import "fmt"

func main() {
	fmt.Println("Hello, World!")
	// 原生的 println() 函数和 fmt.Println() 函数一样，但是它没有格式化输出。
	println("This is my Go project.")
}

/**
   编译和运行
  1. 使用 go build 命令编译 Go 程序,并生成一个可执行文件,然后需要手动执行这个可执行文件。
  go build ./hello.go
  2. 使用 go run 命令运行 Go 程序，它会自动编译并执行程序。省去了中间的可执行文件的生成
  go run ./hello.go
*/

// 每个目录下要只有一个main函数，不要多个main函数，是不符合规范的
