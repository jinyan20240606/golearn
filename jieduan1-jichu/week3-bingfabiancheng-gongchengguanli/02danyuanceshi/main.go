package ch11

// 单元测试命令：go test
// go test 是一个按照一定的约定和组织的测试代码驱动程序，在包目录中，所有以_test.go结尾的文件都会被go test运行
// 我们写的_test.go文件，源码不用担心内容过多，因为go build 会自动忽略掉以_test.go结尾的文件，不会打包到最后的可执行文件中
// test文件有4类：Test开头的功能测试 Benchmark开头的性能测试 Example开头的示例测试 Fuzz开头的模糊测试
// 现在只关心Test和Benchmark两类，其他后面再学

// go中测试文件是与功能文件写在一起的，不像java那样的放在单独的测试文件目录

// 运行测试命令：在目录下执行`go test`命令
// go test -v：普通运行（会执行所有测试）-------- -v = verbose，意思是：详细输出
// go test -v -short 快速运行（跳过这个测试）
// go test -v -run TestAdd2 只运行这一个测试
// go test ：默认不加v参数，就是非详细模式，有失败的只打印失败的，没有失败的才打印通过的简短结果

// 性能测试命令：性能测试命令必须带 -bench=
// go test -bench=. // 运行所有性能测试
// go test -bench=. -benchmem // 性能测试 + 显示内存分配（最常用）
// go test -bench=X // 只跑名为 BenchmarkX 的性能测试
// go test -bench=测试名 -benchtime 1000x // 执行次数指定1000 ，benchtime不写的话默认：自动跑到 ≥1 秒，次数不固定
// go test -bench=. -benchtime=2s // 执行时间指定2秒

// 性能测试输出的每一列含义如下：

// BenchmarkStringSprintf-8：测试函数名，-8 表示 Go 运行时使用的 CPU 核心数为 8。
// 139744：该测试在基准测试期间执行的总次数（迭代次数）。
// 8653 ns/op：每次操作（op = operation）的平均耗时，单位为纳秒（nanoseconds）。数值越小，性能越好。
// 11332 B/op：每次操作平均分配的内存字节数（Bytes per operation）。
// 198 allocs/op：每次操作平均发生的内存分配次数（allocations per operation）。
