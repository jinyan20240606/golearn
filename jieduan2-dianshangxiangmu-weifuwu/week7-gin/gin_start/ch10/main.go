package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

// 学习gin的优雅退出：当我们关闭程序kill或control+c时，应该做的后续处理，避免程序或进程突然异常关闭，导致正在处理的请求被中断，资源没有及时保存等问题
// 监听系统退出信号，让程序可以优雅关闭，而不是暴力中断
// 如微服务中，启动之前或启动之后将当前的服务注册到注册中心，，，我们当前服务停止了之后并没有告知注册中心，注册中心还以为这个服务还在正常运行，继续把请求转发到这个服务上，这时就会出现问题，所以我们需要优雅退出，告诉注册中心这个服务已经停止了，不要再转发请求了

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	go func() {
		router.Run(":8080")
	}()

	// 如果想要接收到系统退出信号 ，用os提供的channel
	// 注意kill命令是能监听到的，但是kill -9 强杀命令-不能监听到
	// 创建一个接收信号的 channel
	quit := make(chan os.Signal)
	// syscall.SIGINT 对应：Ctrl + C，syscall.SIGTERM 对应：系统停止 / 关闭指令，kill PID，服务器重启、容器停止，系统正常关闭程序时发出
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 阻塞，等待接收退出信号
	<-quit
	// 处理后续收尾逻辑
	fmt.Println("关闭server中")
	fmt.Println("注销服务中，清理资源中。。。。")

}
