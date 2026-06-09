package main

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"           // 作用：OpenTelemetry 核心入口（总开关 / 全局 API）
	"go.opentelemetry.io/otel/attribute" // 给 Trace/Span 打标签（键值对）

	// jaeger exporter：直连 Jaeger，不经过 Collector（你现在用的）
	"go.opentelemetry.io/otel/exporters/jaeger" // 导出器：发给 Jaeger

	// 定义当前服务的公共资源信息
	// 	描述服务本身的信息：
	// 服务名
	// 环境（dev/prod）
	// 实例 ID、IP、版本号
	// 所有 Span 都会自动带上这些信息
	// 相当于：服务的身份证
	"go.opentelemetry.io/otel/sdk/resource"

	"go.opentelemetry.io/otel/sdk/trace" // Trace 的具体实现（SDK 核心）

	// 	官方标准属性名（规范字段名）
	// 提供官方统一规范字段，避免大家乱写
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0" // 标准属性
)

func main() {
	// ===================== 1. 配置 Jaeger 地址 =====================
	// Jaeger 服务地址
	url := "http://127.0.0.1:14268/api/traces"

	// 创建 Jaeger 导出器：后面把 Trace 发给 Jaeger
	jexp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		panic(err)
	}

	// ===================== 2. 创建 TracerProvider（全局唯一） =====================
	tp := trace.NewTracerProvider(
		// 批量发送（性能更好）
		trace.WithBatcher(jexp),

		// 配置服务基本信息
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				// 这个ServiceNameKey设置后，jaeger页面搜索筛选项中默认服务名就以这个命名
				semconv.ServiceNameKey.String("mxshop-user"), // 🔥 服务名：最重要！Jaeger 里靠这个找服务
				// 属于全局的，整个应用的信息和属性
				attribute.String("environment", "dev"), // 环境
				attribute.Int("ID", 1),                 // 自定义属性
			),
		),
	)

	// 把创建的 tp 设置成全局 tracer
	otel.SetTracerProvider(tp)

	// ===================== 3. 程序退出时安全关闭 =====================
	ctx, cancel := context.WithCancel(context.Background())
	defer func(ctx context.Context) {
		// 退出前等待 5 秒，把剩余的 span 发完
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		// 优雅关闭时，一般都加个超时上下文，防止它卡住退出不成功
		if err := tp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}(ctx)

	// ===================== 4. 创建一个 Tracer =====================
	tr := otel.Tracer("mxshop-otel")

	// ===================== 5. 开始一个 Span（一条链路的一段） =====================
	_, span := tr.Start(ctx, "func-main") // span 名字叫 func-main

	// ===================== 6. 给 Span级别 设置自定义属性标签，这是更常用 =====================
	var attrs []attribute.KeyValue
	attrs = append(attrs, attribute.String("key1", "value1"))
	attrs = append(attrs, attribute.Bool("key2", false))
	attrs = append(attrs, attribute.Int("key3", 123))
	attrs = append(attrs, attribute.StringSlice("key4", []string{"value4-1", "value4-2"}))

	// 把标签设置到 span 里
	span.SetAttributes(attrs...)

	// ===================== 7. 添加一个事件（记录一个动作） =====================
	span.AddEvent("this is an event")

	// 模拟业务执行耗时 1 秒
	time.Sleep(time.Second)

	// ===================== 8. 结束 Span =====================
	span.End()
}
