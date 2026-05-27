package middlewares

import (
	"fmt"
	"mxshop-api/goods-web/global"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// 这个写法错误
func Trace() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// ===================== 【致命错误】=====================
		// 错误：每次请求都重新初始化 Jaeger、新建 Tracer
		// 正确：Tracer 应该全局初始化一次，不是每次请求都新建！
		// ======================================================
		cfg := jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans:           true,
				LocalAgentHostPort: fmt.Sprintf("%s:%d", global.ServerConfig.JaegerInfo.Host, global.ServerConfig.JaegerInfo.Port),
			},
			ServiceName: global.ServerConfig.JaegerInfo.Name,
		}
		// 创建一个全局tracer
		// ===================== 【错误】每次请求都 NewTracer =====================
		// 后果：fd暴涨、连接泄漏、性能极低、链路混乱
		tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
		if err != nil {
			panic(err)
		}
		// 每次请求都设置全局Tracer，完全错误！
		opentracing.SetGlobalTracer(tracer)
		// 每次请求都 defer closer.Close()，会导致立即关闭，链路直接断了
		defer closer.Close()
		// ===================== 启动根Span =====================
		// 以当前请求路径作为Span名称
		startSpan := tracer.StartSpan(ctx.Request.URL.Path)
		defer startSpan.Finish()
		// 将 tracer 和 父Span 存入 Gin Context
		ctx.Set("tracer", tracer) // 后续在子span中使用时，直接从 Gin Context 中获取即可，用tracer创建关联span父子关系
		ctx.Set("parentSpan", startSpan)
		// 执行后续中间件 / 接口逻辑
		ctx.Next()
	}
}
