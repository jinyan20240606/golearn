package main

import (
	"flag"
	"fmt"

	// Go 官方提供的插件生成器框架
	"google.golang.org/protobuf/compiler/protogen"
	// 定义了插件通信的数据结构（Protocol Buffer Stubs）
	"google.golang.org/protobuf/types/pluginpb"
	// 自己编写的业务逻辑包
	"GoStart/tools/generator"
)

var release = "v1.0.0" // 假设的版本号变量
// 定义一个布尔类型的命令行参数
// 返回的是一个 *bool 指针
var showVersion = flag.Bool("version", false, "print the version and exit")

func main() {
	// protoc 在调用插件时，可能会通过命令行参数传递一些配置，如：protoc --go_out=. --go-grpc_out=. helloworld.proto base.proto
	flag.Parse() // // 解析命令行参数，配合flag.Bool一起使用

	// 注意：因为 showVersion 是指针，所以需要用 * 来解引用获取它的实际值
	if *showVersion {
		fmt.Printf("protoc-gen-go-errors %v\n", release)
	}

	var flags flag.FlagSet // 声明了一个新的 Flag 集合，目前还没有值，用于接收 protogen.Options 内部处理的参数

	// 	protogen 库的核心启动逻辑。
	// ParamFunc: flags.Set：告诉 protogen，如果有额外的参数，就使用 flags 这个集合来接收。
	// .Run(...)：启动插件。
	protogen.Options{
		ParamFunc: flags.Set, // 解析出来的插件参数最终都会被存入你定义的 flags 变量里
	}.Run(func(gen *protogen.Plugin) error {
		// 声明插件支持的特性
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		// gen.Files：获取 protoc 传进来的所有 .proto 文件列表，比如命令行里面我放了2个proto文件需要解析，那么gen.Files就会有2个文件元素
		for _, f := range gen.Files {
			// 检查该文件是否被标记为“需要生成”（有时候一个 import 的文件只是用来引用类型，不需要单独生成代码）。
			if !f.Generate {
				continue
			}
			generator.GenerateFile(gen, f)
		}
		return nil
	})
}
