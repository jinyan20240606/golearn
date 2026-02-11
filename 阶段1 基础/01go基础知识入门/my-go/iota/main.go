package main

// iota: 是特殊的常量生成器，默认是int类型，常用于定义枚举类型和位掩码等。它在 const 块中自动递增，简化了常量的定义过程。，可以自动生成一系列相关的常量值，避免手动赋值，并且使代码更具可读性和维护性。
//  ------- 可以认为是可以被编译器修改的常量
// iota 是 Go 语言中的一个预定义标识符，代表常量生成器。它在 const 声明中使用，可以自动生成一系列相关的常量值。iota 的值在每个 const 块中从 0 开始，并且在每次遇到新的 const 声明时重置为 0。

// iota 的使用非常灵活，可以用于生成枚举类型、位掩码等。下面是一些示例：

const (
	// 使用 iota 生成枚举类型
	Red   = iota // Red 的值为 0
	Green        // Green 的值为 1 ----- 后面也默认是前一个的Red值，相当于这个3个常量都定义的iota，然后编译器会在这个常量里自动识别计算
	Blue         // Blue 的值为 2
)

const (
	ERRCODE_SUCCESS = iota // 后续的值自动从0递增
	ERRCODE_FAIL    = iota // FAIL 的值为 1
	ERRCODE_NO_DATA = iota // NO_DATA 的值为 2
)

const (
	// 使用 iota 生成位掩码
	Read    = 1 << iota // Read 的值为 1 (1 << 0)
	Write               // Write 的值为 2 (1 << 1)
	Execute             // Execute 的值为 4 (1 << 2)
)

const (
	// 使用 iota 生成连续的常量值
	A = iota + 1 // A 的值为 1 (0 + 1)
	B            // B 的值为 2 (1 + 1)
	C            // C 的值为 3 (2 + 1)
)

const (
	A1 = iota
	B1
	C1 = "hello"
	D1 = "world"
	D11
	D12 = 100
	E1  = iota
) // 输出：0 1 hello world world 6

// 虽然中间定义了字符串即中断了iota的连续，但 still iota会自动递增，所以C1的值为"hello"，D1的值为"world"，E1的值为4

func main() {
	println(Red, Green, Blue)             // 输出: 0 1 2
	println(Read, Write, Execute)         // 输出: 1 2 4
	println(A1, B1, C1, D1, D11, D12, E1) // 输出: 0 1 hello world world 6
}

// 通过使用 iota，可以简化常量的定义，避免手动赋值，并且使代码更具可读性和维护性。
