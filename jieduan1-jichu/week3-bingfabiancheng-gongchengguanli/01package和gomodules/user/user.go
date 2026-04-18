package course

import "fmt"

type Course struct {
	Name string
}

func (c Course) GetName() string {
	return c.Name
}

// 这是内置魔术方法，名字必须是init，被导入时自动执行init方法
func init() {
	fmt.Println("初始化了init")
}
