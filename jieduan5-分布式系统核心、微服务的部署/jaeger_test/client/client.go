package main

import (
	"context"
	"fmt"

	// 开源的 gRPC + OpenTracing 集成库（自动帮你处理Inject/Extract/父子Span）
	"OldPackageTest/jaeger_test/otgrpc"
	// OpenTracing 标准接口
	"github.com/opentracing/opentracing-go"
	// Jaeger 客户端
	"github.com/uber/jaeger-client-go"
	// Jaeger 配置
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"

	"OldPackageTest/grpc_test/proto"
)

func main() {
	// ===================== 1. 初始化 Jaeger 配置 =====================
	// 配置采样、上报地址、服务名
	cfg := jaegercfg.Configuration{
		// 采样配置：全部采样（1=全采，0=不采）
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		// 报告器配置：发送到 jaeger-agent
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "192.168.0.104:6831",
		},
		ServiceName: "mxshop", // 当前服务名称（在Jaeger UI显示）
	}
	// ===================== 2. 创建 Tracer =====================
	// 根据配置生成 Tracer（链路追踪管理器）
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	// 设置全局 Tracer，整个程序都能用
	opentracing.SetGlobalTracer(tracer)
	// 程序退出时关闭，把剩余的链路数据刷到 Jaeger
	defer closer.Close()

	// ===================== 3. 创建 gRPC 连接，并加入 链路追踪拦截器 =====================
	// 重点：
	// otgrpc.OpenTracingClientInterceptor(tracer)
	// 这是 gRPC 客户端拦截器，**自动帮你完成：
	// 1. 创建Span
	// 2. 建立父子关系
	// 3. Inject 把追踪信息（trace-id、span-id）放入gRPC Header请求头，发给服务端
	// // ---- 具体服务端怎么使用这些元数据，还得去服务端改造集成追踪系统
	// 4. 当每个grpc请求响应完成后，自动计时、自动结束、自动上报
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure(),
	// ✅ 核心：gRPC 自动链路追踪拦截器
	// 它是 gRPC 客户端的一个 “中间件”，在你发起 gRPC 请求的前后，自动完成分布式追踪的所有脏活累活：创建 span、注入 trace、计时、结束、上报。
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer()))
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// ===================== 4. 调用 gRPC 接口 =====================
	c := proto.NewGreeterClient(conn)
	// 调用时，拦截器会**自动创建Span**，并传递Trace信息
	// 你完全不用手动写 StartSpan / ChildOf / Inject
	r, err := c.SayHello(context.Background(), &proto.HelloRequest{Name: "bobby"})
	if err != nil {
		panic(err)
	}
	fmt.Println(r.Message)
}
