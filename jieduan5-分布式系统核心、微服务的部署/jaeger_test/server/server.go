package main

import (
	"context"
	"net"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"

	"OldPackageTest/grpc_test/proto"
	"OldPackageTest/jaeger_test/otgrpc" // 必须加这个
)

type Server struct{}

func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply, error) {
	// ======================================
	// 👇 这里就是：从 ctx 里取出当前的 Span
	// 拦截器已经把 Span 放进去了，直接拿就能用
	// ======================================
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		// 可以给当前 Span 加自定义标签（会显示在 Jaeger 里）
		span.SetTag("my.request.name", request.Name)
		span.SetTag("my.business.type", "hello")

		// 可以打日志（会显示在 Jaeger Span 的 Logs 里）
		span.LogKV("event", "执行SayHello成功", "name", request.Name)
	}

	// 你的业务逻辑不变
	return &proto.HelloReply{
		Message: "hello " + request.Name,
	}, nil
}

func main() {
	// ============ 1. 初始化 Jaeger 不变 ============
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "192.168.0.104:6831",
		},
		ServiceName: "mxshop_server",
	}

	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	// ============ 2. 关键：服务端初始化时加入追踪拦截器 ============
	// 这一行自动从 gRPC 元数据提取 Trace 上下文
	// 1. 客户端 Inject → gRPC 头带 TraceID/SpanID 过去
	// 2. 服务端 Extract → 从 gRPC 头取出，得到 parentSpanContext
	// 3. 使用 parent 创建 子Span → 链路连上
	// 4. 子Span放进 ctx → 给业务代码用
	// 5. 业务结束 → 自动结束子Span、上报Jaeger
	g := grpc.NewServer(
		grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)),
	)

	// 注册服务
	proto.RegisterGreeterServer(g, &Server{})

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	// 启动服务
	err = g.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}
