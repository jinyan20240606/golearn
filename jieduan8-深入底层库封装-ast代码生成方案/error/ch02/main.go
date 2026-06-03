package main

import (
	"errors"
	"fmt"

	"github.com/pkg/errors"
)

type DBError struct {
	msg string
}

func (d *DBError) Error() string {
	return d.msg
}

var ErrNameEmpty = &DBError{"name can't be empty"}

func (s *Student) SetName(name string) error {
	if name == "" || s.Name == "" {
		return ErrNameEmpty
	}
	s.Name = name
	return nil
}

type Student struct {
	Name string
	Age  int
}

func NewStudent() (*Student, error) {
	stu := &Student{
		Age: 18,
	}
	err := stu.SetName("")
	if err != nil {
		// 我想加上我自定义的错误，又想保留原始的调用栈，不能使用new会丢栈
		return stu, errors.Wrap(err, "set name faild")
	}
	return stu, nil
}

func main() {
	_, e := NewStudent()
	//%+v 显示错误信息和堆栈信息，比如文件名 行号等
	var perr *DBError
	// if errors.Is(e, ErrNameEmpty) {
	// 	fmt.Println("match")
	// }
	// As 关心错误的类型，怕奴蛋， Is 关心错误的值
	if errors.As(e, &perr) {
		fmt.Println("match")
	}
}
