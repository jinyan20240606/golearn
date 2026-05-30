package main

// map的范型
type Mymap[KEY int | string, VALUE float32 | float64] map[KEY]VALUE

// 结构体的范型
type Man struct {
}
type Woman struct {
}

type Company[T Man | Woman] struct {
	Name string
	CEO  T
}

// 通道的范型
type MyChannel[T int | string] chan T

// 类型嵌套
type WowStruct[T string | int, S []T] struct {
	A T
	B S
}

// 错误用法约束：
// 错误用法1, 类型参数不能单独使用，在 Go 里是非法语法，但在 TypeScript 里是完全合法、正常、常用的！
// type CommonType[T int | string] T

// 错误用法2
// type CommonType[T *int | string] []T   // Go 不允许在类型参数列表直接写指针类型 *int！，范型的语法中会当成乘法号
// 为了解决上述的语法歧义，Go 语言要求在类型约束包含容易产生误解的符号（如 *）时，必须使用 interface{} 将其包裹起，这是go中的特殊规定，不算普通接口类型语法了，因为普通接口类型语法是里面定义方法类型
// 接口语法扩展：interface{ ~int } 表示：不仅包含 int 本身，还包含所有底层类型是 int 的类型（比如 type MyInt int）
type CommonType[T interface{ *int } | string] []T

// 错误用法3 匿名结构体不支持泛型
// 错误用法4 泛型不支持switch判断即类型断言，只能使用反射去获取类型
//  ---- 无法直接对泛型变量使用 switch v.(type) 或 v.(SomeType) 进行类型断言

// 错误用法5 匿名函数不支持泛型

func main() {
	//company := Company[Man]{
	//	Name: "bobby",
	//	CEO: Man{},
	//}

	//company := Company[Woman]{
	//	Name: "bobby",
	//	CEO:  Woman{},
	//}

	//var c MyChannel[string]

	//几种常见的错误

}
