package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

const (
	traceName = "mxshop-otel"
)

// 子函数A：创建一个子Span
func funcA(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	tr := otel.Tracer("traceName")                       // 获取Tracer
	_, span := tr.Start(ctx, "func-a")                   // 基于父ctx创建子Span：func-a
	span.SetAttributes(attribute.String("name", "funA")) // 给子Span设置标签

	type _LogStruct struct {
		CurrentTime time.Time `json:"current_time"`
		PassWho     string    `json:"pass_who"`
		Name        string    `json:"name"`
	}

	logTest := _LogStruct{
		CurrentTime: time.Now(),
		PassWho:     "bobby",
		Name:        "func-a",
	}

	b, _ := json.Marshal(logTest)
	span.SetAttributes(attribute.Key("这是测试日志的key").String(string(b)))
	time.Sleep(time.Second)
	span.End()
}

func funcB(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	tr := otel.Tracer("traceName")
	_, span := tr.Start(ctx, "func-b")
	fmt.Println("trace:", span.SpanContext().TraceID(), span.SpanContext().SpanID())
	time.Sleep(time.Second)
	span.End()
}

func main() {
	// ==================== 1. 初始化Jaeger导出器（直连，不经过Collector）=====================
	url := "http://127.0.0.1:14268/api/traces"
	jexp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		panic(err)
	}

	// ==================== 2. 创建TracerProvider（全局追踪管理器）=====================
	tp := trace.NewTracerProvider(
		trace.WithBatcher(jexp),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("mxshop-user"),
				attribute.String("environment", "dev"),
				attribute.Int("ID", 1),
			),
		),
	)
	otel.SetTracerProvider(tp) // 设置为全局TracerProvider
	// ==================== 3. 程序退出时优雅关闭，保证数据发送完毕 =====================
	ctx, cancel := context.WithCancel(context.Background())
	defer func(ctx context.Context) {
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}(ctx)
	// ===================== 创建一个 Tracer =====================
	tr := otel.Tracer(traceName)
	// ==================== 4. 创建主Span =====================
	spanCtx, span := tr.Start(ctx, "func-main")
	// ==================== 5. 并发调用两个函数，传递ctx保持链路 =====================
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go funcA(spanCtx, wg)
	go funcB(spanCtx, wg)

	span.AddEvent("this is an event")
	time.Sleep(time.Second)
	wg.Wait()
	span.End()
}
