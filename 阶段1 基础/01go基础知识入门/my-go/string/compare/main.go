package main

// 字符串的比较
func main() {
	print("hello world")

	a := "hello"
	b := "hello"

	if a == b { // 注意等于符号是2个=，不等于是 !=
		println("a == b")
	} else {
		println("a != b")
	}

	// 字符串的大小比较，字符之间的逐个码点大小比较
	if a < b {
		println("a < b")
	} else {
		println("a >= b")
	}
	// 字符串的长度比较
	if len(a) < len(b) {
		println("len(a) < len(b)")
	} else {
		println("len(a) >= len(b)")
	}
}
