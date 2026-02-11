package main

// 常量类型只可以定义：数值（整数-浮点-数复数）、字符串、布尔值
// 不曾使用的常量：没有强制使用的硬性要求

func main() {
	const pi float32 = 3.14
	const world = "世界"
	const (
		hello = "hello"
		hello2
		world1 = "世界"
		num    = 100
		num2
	)
	// 常量细节：输出：hello hello 世界 100 100。未赋值的常量会自动复制前一个常量的值
	println(pi)
	println(world)
}
