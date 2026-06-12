//go:build wireinject
// +build wireinject

// 语法：//go:build 标签名
// 含义：声明当前文件仅在编译条件包含 wireinject 时才被编译 / 解析。---- 不会触发IDE的与wire_gen.go重名编译语法校验提示
// 场景：wire 工具执行时，会自动带上 wireinject 标签，所以能识别并解析这个文件；
// 普通 go 命令默认不带该标签，文件被跳过。
// # 编译当前包，启用 wireinject 标签
// go build -tags="wireinject"

// # 直接运行
// go run -tags="wireinject" .

// 旧版构建标签（兼容 Go 1.16 及更早版本）
// Go 历史语法，作用和上一行完全一致。
// 为什么写两行？
// 为了向前、向后全版本兼容：
// 新版 Go 识别 //go:build
// 旧版 Go 识别 // +build
// 两行并存是 Go 社区标准写法。

package main

// 导入 wire 核心包，仅用于声明 wire.Build，运行时不依赖该包
import "github.com/google/wire"

// initEvent 是 Wire 注入器：对外暴露的组件初始化入口，最终返回目标类型 Event
// 函数名完全自定义，只要符合规范即可
func initEvent() Event {
	// 参数：传入所有提供者函数（Provider），这个每个Provider函数必须符合的规则： 每个函数的返回值都是其他函数的参数；
	// 功能：Wire 工具编译期静态分析所有 Provider 的依赖关系，自动梳理调用顺序、校验依赖完整性、检查循环依赖。
	panic(wire.Build(NewEvent, NewGreeter, NewMessage))
}
