package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"                         // OTel 全局核心API
	"go.opentelemetry.io/otel/attribute"               // 给Span设置标签/属性
	"go.opentelemetry.io/otel/exporters/jaeger"        // Jaeger导出器（直连，不走Collector）
	"go.opentelemetry.io/otel/propagation"             // 链路传播：跨服务传递TraceID
	"go.opentelemetry.io/otel/sdk/resource"            // 服务资源（服务名、环境等全局标签）
	"go.opentelemetry.io/otel/sdk/trace"               // Trace SDK核心实现
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0" // OTel标准字段规范

	"github.com/valyala/fasthttp" // HTTP客户端，用于跨服务调用

	"GoStart/log" // 自定义日志包（带Trace上下文关联）
)

const (
	traceName = "mxshop-otel" // Tracer名称
)

var tp *trace.TracerProvider // 全局TracerProvider（管理所有追踪）

// tracerProvider 初始化OTel追踪（连接Jaeger、设置全局Tracer）
func tracerProvider() error {
	// Jaeger采集端地址（直连模式）
	url := "http://127.0.0.1:14268/api/traces"
	// 创建Jaeger导出器
	jexp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		panic(err)
	}

	// 创建全局TracerProvider
	tp = trace.NewTracerProvider(
		trace.WithBatcher(jexp), // 批量发送Span（提升性能）
		trace.WithResource( // 设置服务公共信息
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("mxshop-user"), // 服务名（Jaeger显示用）
				attribute.String("environment", "dev"),       // 环境
				attribute.Int("ID", 1),                       // 自定义ID
			),
		),
	)
	// SetTracerProvider方法：安装 “ tracing 工具箱”
	// 它不干活，不生成 ID，不创建 Span！
	// 它只是：
	// 放好工具
	// 配好 Jaeger 地址
	// 配好服务名
	// 配好全局格式
	// 当第一次调用 tr.Start (ctx, "xxx")，并且 ctx 里没有任何 trace 的时候，才会生成新 traceId
	//已有上游链路时 → 继续沿用旧的 traceId，只生成新的 spanId
	otel.SetTracerProvider(tp) // 设置全局TracerProvider,全局生成 traceId/spanId

	// 设置全局传播器：跨服务传递TraceID（W3C标准）
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return nil
}

// funcA 子函数：创建子Span，打日志、设置标签
func funcA(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done() // 函数结束，WaitGroup-1

	tr := otel.Tracer("traceName")                       // 获取Tracer
	spanCtx, span := tr.Start(ctx, "func-a")             // 创建子Span：func-a
	span.SetAttributes(attribute.String("name", "funA")) // 给Span打标签

	// 自定义日志结构体
	type _LogStruct struct {
		CurrentTime time.Time `json:"current_time"`
		PassWho     string    `json:"pass_who"`
		Name        string    `json:"name"`
	}

	// 构造日志内容
	logTest := _LogStruct{
		CurrentTime: time.Now(),
		PassWho:     "bobby",
		Name:        "func-a",
	}

	b, _ := json.Marshal(logTest) // 转为JSON

	// 带Trace上下文打印日志（日志会自动带上TraceID）
	log.InfofC(spanCtx, "this is funca log: %s", string(b))

	// 将JSON日志作为Span标签存入（可在Jaeger查看）
	span.SetAttributes(attribute.Key("这是测试日志的key").String(string(b)))

	time.Sleep(time.Second) // 模拟业务耗时
	span.End()              // 结束Span
}

// funcB 子函数：创建子Span + 跨服务HTTP调用（自动传递Trace）
func funcB(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done() // 函数结束，WaitGroup-1

	tr := otel.Tracer("traceName")                                                   // 获取Tracer
	spanCtx, span := tr.Start(ctx, "func-b")                                         // 创建子Span：func-b
	fmt.Println("trace:", span.SpanContext().TraceID(), span.SpanContext().SpanID()) // 打印TraceID/SpanID
	time.Sleep(time.Second)

	// 构造HTTP请求
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://127.0.0.1:8090/server") // 调用服务端
	req.Header.SetMethod("GET")

	// ==================== 核心：跨服务自动传递Trace ====================
	// 获取全局传播器
	p := otel.GetTextMapPropagator()
	// 创建Header载体
	headers := make(map[string]string)
	// 将当前Trace上下文注入Header的传播map包裹
	p.Inject(spanCtx, propagation.MapCarrier(headers))

	// 将Trace信息设置到HTTP Header中
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	//req.Header.Set("trace-id", span.SpanContext().TraceID().String())
	//req.Header.Set("span-id", span.SpanContext().SpanID().String())

	// 发送HTTP请求
	fclient := fasthttp.Client{}
	fres := fasthttp.Response{}
	_ = fclient.Do(req, &fres)

	// 带Trace上下文打日志
	log.InfofC(spanCtx, "this is funcB log: %s", "imooc")

	span.End() // 结束Span
}

func main() {
	_ = tracerProvider() // 初始化OTel追踪

	// 创建根上下文
	ctx, cancel := context.WithCancel(context.Background())
	// 程序退出时优雅关闭TracerProvider，确保数据发送完毕
	defer func(ctx context.Context) {
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}(ctx)

	// 创建根Span：func-main
	tr := otel.Tracer(traceName)
	spanCtx, span := tr.Start(ctx, "func-main")

	// 并发启动两个子函数
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go funcA(spanCtx, wg) // 传递父上下文，形成父子链路
	go funcB(spanCtx, wg)

	// 给根Span添加事件
	span.AddEvent("this is an event")
	time.Sleep(time.Second)

	wg.Wait()  // 等待所有协程完成
	span.End() // 结束根Span
}
