package main

import (
	"fmt"

	perrors "github.com/pkg/errors" // 导入错误处理包
)

/*
go的error和其他语言的try catch不一样， go语言将错误和异常分开，其他语言认为错误和异常是一回事，都用trycatch处理
go中认为error是一种值，不算异常，每次调用函数时都会返回2个值一个是正常的返回值，一个是error类型的错误值，go中的异常认为是panic那种类型
*/

// 除法函数
func divFunc(a, b int) (int, error) {
	if b == 0 {
		return 0, perrors.New("b can't be zero")
	}
	return a / b, nil
}

func main() {
	var a, b = 1, 0
	ret, err := divFunc(a, b)
	if err != nil {
		// %+v	带详细结构打印（如果有字段 / 堆栈会展开）
		fmt.Printf("ret is %d, err is %+v\n", ret, err)
	}
}
