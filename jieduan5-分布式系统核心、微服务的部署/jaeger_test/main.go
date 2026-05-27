package main

import (
	"time"
	// OpenTracing 标准接口
	"github.com/opentracing/opentracing-go"
	// Jaeger 客户端实现
	"github.com/uber/jaeger-client-go"
	// Jaeger 配置包
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func main() {
	// ===================== 1. 配置 Jaeger =====================
	cfg := jaegercfg.Configuration{
		// 采样
		Sampler: &jaegercfg.SamplerConfig{
			// 采样配置：决定哪些请求需要被追踪
			Type:  jaeger.SamplerTypeConst, // 固定采样策略（全采样/不采样）--- 生产环境必须采样！，全采样会导致性能下降 + Jaeger 存储爆掉
			Param: 1,                       // 1=全采样，0=不采样
		},
		// 报告器配置：数据发送到哪里
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true, // 是否打印 span 日志
			LocalAgentHostPort: "192.168.0.104:6831",
		},
		ServiceName: "mxshop",
	}
	// ===================== 2. 初始化 Tracer =====================
	// 根据配置创建 Tracer（追踪管理器）
	// closer：用于程序退出前关闭上报，把剩余数据发送出去
	// tracer：符合 OpenTracing 标准的追踪器
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	// ===================== 3. 设置全局 Tracer =====================
	// 将创建好的 tracer 设置为全局 tracer，方便项目任意地方使用
	opentracing.SetGlobalTracer(tracer)
	// 程序退出时，关闭 tracer，确保所有 span 都发送到 Jaeger
	defer closer.Close()
	// ===================== 4. 创建并启动一个 Span =====================
	// 启动一个名为 "go-grpc-web" 的 Span（代表一段操作/调用）
	span := opentracing.StartSpan("go-grpc-web")
	// 模拟业务耗时
	time.Sleep(time.Second)
	// ===================== 5. 结束 Span =====================
	// 结束 Span，计算耗时，并异步上报给 Jaeger，
	// 必须手动调用 span.Finish()很重要！！！！
	defer span.Finish()

	// 有多个defer时，defer的执行顺序是后进先出
}
