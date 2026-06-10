package main

import (
	"github.com/gin-gonic/gin"
	// Prometheus 官方 SDK 核心包
	// 作用：提供 Counter、Gauge、Histogram 等指标类型，所有指标的定义、类型、结构都来自这里
	"time"

	"github.com/prometheus/client_golang/prometheus"
	// 自动注册自定义的指标
	//作用：不用手动把指标注册到 Prometheus 注册表，自动帮你注册好
	"github.com/prometheus/client_golang/prometheus/promauto"
	// 提供 /metrics 接口
	// 作用：把内存里的指标转换成 Prometheus 能识别的文本格式
	// 就是 Prometheus 拉取的 /metrics 接口
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	// 无限循环，每 2 秒给指标 +1
	for {
		ops.Inc()

		time.Sleep(2 * time.Second)
	}
}

var (
	// 	自定义一个 Prometheus 指标：
	// 类型：Counter（只增不减）
	// 名字：mxshop_test（Prometheus 里看到的指标名）
	// Help：说明文字
	ops = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mxshop_test",
		Help: "just for test",
	})
)

func main() {
	//	启动一个后台协程
	//让指标自动增加，不阻塞 Web 服务
	go recordMetrics()
	r := gin.Default()
	// gin.WrapH → 把标准 http handler 包装成 gin handler
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.Run(":8050")
}
