package trace

import (
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"

	"mxshop/pkg/log"
)

/*
初始化不同的export的设置
*/

const (
	kindJaeger = "jaeger"
	kindZipkin = "zipkin"
)

var (
	//set ,struct 空结构体不占内存， zerobase
	agents = make(map[string]struct{})
	lock   sync.Mutex
)

func InitAgent(o Options) {
	lock.Lock()
	defer lock.Unlock()

	_, ok := agents[o.Endpoint]
	if ok {
		return
	}
	err := startAgent(o)
	if err != nil {
		return
	}
	agents[o.Endpoint] = struct{}{}
}

// startAgent 启动 OpenTelemetry 追踪代理（用于分布式链路追踪）
func startAgent(o Options) error {
	// 声明 链路追踪的导出器（把 spans 数据发给 Jaeger/Zipkin）
	var sexp trace.SpanExporter
	// 声明错误变量
	var err error
	// 构建 TracerProvider 的配置选项
	opts := []trace.TracerProviderOption{
		// 采样器：基于父级采样 + 按比例采样（o.Sampler 是采样率，如 0.1=10%）
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(o.Sampler))),
		// 设置服务名（在追踪平台里显示的服务名称）
		trace.WithResource(resource.NewSchemaless(semconv.ServiceNameKey.String(o.Name))),
	}
	// 如果配置了上报地址（Endpoint），才创建导出器
	if len(o.Endpoint) > 0 {
		// 根据不同的 Batcher 类型创建对应的导出器
		switch o.Batcher {
		case kindJaeger: // 使用 Jaeger 做链路存储
			sexp, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(o.Endpoint)))
			if err != nil {
				return err
			}
		case kindZipkin:
			sexp, err = zipkin.New(o.Endpoint)
			if err != nil {
				return err
			}
		}
		// 把导出器加入 TracerProvider 选项
		opts = append(opts, trace.WithBatcher(sexp))
	}
	// 创建 TracerProvider（核心：追踪管理器）
	tp := trace.NewTracerProvider(opts...)
	// 设置全局 TracerProvider（整个程序都用这个追踪）
	otel.SetTracerProvider(tp)
	// 设置全局传播器：跨服务传递 traceId、spanId（支持 W3C TraceContext +  baggage）
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	// 设置全局错误处理器：追踪内部出错时打印日志，不影响主业务
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Errorf("[otel] error: %v", err)
	}))
	return nil
}
