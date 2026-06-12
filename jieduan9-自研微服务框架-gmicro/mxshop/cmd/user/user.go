package main

import (
	"math/rand"
	"mxshop/app/user/srv"
	"os"
	"runtime"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	// 	Cobra 设计上支持两种启动方式：
	// 常规方式：程序编译为二进制，在 终端执行，Cobra 自动解析命令行参数；
	// 内嵌调用：代码里手动构造 cobra.Command，直接调用 Run()/Execute()，不依赖终端、不解析进程启动参数。
	// 启动服务
	srv.NewApp("user-server").Run()
}
