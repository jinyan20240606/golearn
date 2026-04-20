package ch11

// 测试的文件名必须是：xxx_test.go
import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

// 01-功能测试函数名必须以 Test 开头，且参数是 (t *testing.T)
// testing.Short()：快速测试模式，跳过慢测试
// t.Error()：断言失败 + 继续执行
// t.Fatal()：断言失败 + 直接终止测试（更常用）
func TestAdd(t *testing.T) {
	re := add(1, 2)
	if re != 3 {
		t.Error("add error")
	}
	t.Log("add success")
}

// 如何在测试代码中根据配置参数跳过一些代码中认为耗时的测试：go test -v -short
func TestAdd2(t *testing.T) {
	// 如果运行时加了 -short 参数，就跳过这个测试
	if testing.Short() {
		t.Skip("skip add2---short模式下跳过")
	}
	re := add(1, 5)
	if re != 6 {
		t.Error("add error")
	}
	t.Log("add success")
}

// 表格驱动测试:就是代码的写法，用定义一个表格结构来描述测试数据，然后根据表格中的数据进行测试，比较方便
func TestAdd3(t *testing.T) {
	var tests = []struct {
		a, b int
		want int
	}{
		{1, 2, 3},
		{2, 2, 4},
		{3, 2, 5},
	}
	for _, test := range tests {
		re := add(test.a, test.b)
		if re != test.want {
			t.Errorf("add(%d, %d) = %d; want %d", test.a, test.b, re, test.want)
		}
	}
}

// 02-性能测试函数名必须以 Benchmark 开头，且参数是 (b *testing.B)
func BenchmarkAdd(bb *testing.B) {
	var a, b, c int
	a = 125
	b = 451
	c = 576
	for i := 0; i < bb.N; i++ { // b.N 表示测试的次
		// Go 性能测试默认没有固定次数默认：自动跑到 ≥1 秒，次数不固定，它是自动调整的，一开始执行 1 次，自动增加次数，直到运行时间 ≥ 1 秒
		if actual := add(a, b); actual != c {
			bb.Error("add error")
		}
		bb.Log("add success", i)
	}
}

const numbers = 100

// 对字符串拼接的不同方法做性能测试
func BenchmarkStringSprintf(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var str string
		for j := 0; j < numbers; j++ {
			// Go 标准库fmt包的格式化函数，返回拼接后的字符串（不打印）
			str = fmt.Sprintf("%s%d", str, j)
		}
	}
	b.StopTimer()
}
func BenchmarkStringAdd(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var str string
		for j := 0; j < numbers; j++ {
			// 整型转字符串
			str += strconv.Itoa(j)
		}
	}
	b.StopTimer()
}

func BenchmarkStringBuilder(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var builder strings.Builder
		for j := 0; j < numbers; j++ {
			builder.WriteString(strconv.Itoa(j))
		}
		_ = builder
	}
	b.StopTimer()
}

// 上面的性能测试结果：命令：go test -bench=. -benchmem
// 对比结论：

// StringBuilder 每次操作仅耗时 525 ns，内存分配 6 次，性能最优。
// StringAdd 次之，耗时 3144 ns，分配 99 次。
// StringSprintf 最慢，耗时 8653 ns，分配 198 次，内存开销最大。
