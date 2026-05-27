package main

import (
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func main() {
	// 1. 初始化 jaeger 配置
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1, // 全采样
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "192.168.0.104:6831",
		},
		ServiceName: "mxshop-span-demo",
	}

	// 2. 创建tracer
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	// ==========================================
	// 【第一步：创建 父 Span】
	// ==========================================
	parentSpan := opentracing.StartSpan("main_process")
	defer parentSpan.Finish() // 结束父span
	time.Sleep(100 * time.Millisecond)

	// ==========================================
	// 【第二步：创建 子 Span —— 关键！】
	// opentracing.ChildOf(parentSpan.Context())
	// ==========================================
	childSpan := opentracing.StartSpan(
		"call_user_service",
		opentracing.ChildOf(parentSpan.Context()), // 指定父
	)
	time.Sleep(200 * time.Millisecond)
	childSpan.Finish()

	// ==========================================
	// 【第三步：再创建一个孙子】
	// ==========================================
	childSpan2 := opentracing.StartSpan(
		"call_order_service",
		opentracing.ChildOf(childSpan.Context()),
	)
	time.Sleep(150 * time.Millisecond)
	childSpan2.Finish()

	// 注意细节：多个span时，直接调用childSpan2.Finish()，不用加defer
}
