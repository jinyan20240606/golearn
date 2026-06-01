package monkey

import (
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
)

// mock 一个函数, mock 范围更广， 而且不需要事先生成代码， 大家可以结合自己的需求
// // 使用方式1:对函数进行monkey测试：特点填写的是函数变量
func TestCompute(t *testing.T) {
	//动态的补丁技术
	patches := gomonkey.ApplyFunc(networkCompute, func(a, b int) (int, error) {
		return 2, nil
	})
	defer patches.Reset()

	sum, err := Compute(1, 2)
	if err != nil {
		t.Error(err)
	}
	if sum != 3 {
		t.Errorf("sum is %d, want 3", sum)
	}
}

// 使用方式2:对结构体的方法进行monkey测试：特点填写的是方法名字符串必须大写
func TestCompute2(t *testing.T) {
	//动态的补丁技术
	var c *Computer
	patches := gomonkey.ApplyMethod(reflect.TypeOf(c), "NetworkCompute", func(_ *Computer, a, b int) (int, error) {
		return 2, nil
	})
	defer patches.Reset()

	c = &Computer{}
	sum, err := c.Compute(1, 2)
	if err != nil {
		t.Error(err)
	}
	if sum != 3 {
		t.Errorf("sum is %d, want 3", sum)
	}
}

var num = 10

// 使用方式3: 对全局变量进行monkey测试：特点：参数直接传入变量名
func TestGlobalVar(t *testing.T) {
	patches := gomonkey.ApplyGlobalVar(&num, 12)
	defer patches.Reset()

	if num != 10 {
		t.Errorf("expected %v, got: %v", 10, num)
	}
}
