package restserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/penglongli/gin-metrics/ginmetrics"

	mws "mxshop/gmicro/server/restserver/middlewares"
	"mxshop/gmicro/server/restserver/pprof"
	"mxshop/gmicro/server/restserver/validation"
	"mxshop/pkg/errors"
	"mxshop/pkg/log"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
)

type JwtInfo struct {
	// defaults to "JWT"
	Realm string
	// defaults to empty
	Key string
	// defaults to 7 days
	Timeout time.Duration
	// defaults to 7 days
	MaxRefresh time.Duration
}

// 把gin.Engine封装成Server
type Server struct {
	*gin.Engine

	//端口号， 默认值 8080
	port int

	//开发模式， 默认值 debug
	mode string

	//是否开启健康检查接口， 默认开启， 如果开启会自动添加 /health 接口
	healthz bool

	//是否开启pprof接口， 默认开启， 如果开启会自动添加 /debug/pprof 接口
	enableProfiling bool

	//是否开启metrics接口， 默认开启， 如果开启会自动添加 /metrics 接口
	enableMetrics bool

	//中间件: 字符串，并没有支持用户传进函数进来，而是字符串，我们在这里只内置gin成熟的中间件即可
	middlewares []string

	//jwt配置信息
	jwt *JwtInfo

	//翻译器, 默认值 zh
	transName string
	trans     ut.Translator

	server *http.Server

	serviceName string
}

// 入口初始化方法，返回一个Server结构体
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		port:            8080,
		mode:            "debug",
		healthz:         true,
		enableProfiling: true,
		jwt: &JwtInfo{
			"JWT",
			"mwGDMGtSpdwXaiihF5WnEgRajSFpdZj8",
			7 * 24 * time.Hour,
			7 * 24 * time.Hour,
		},
		Engine:      gin.Default(),
		transName:   "zh",
		serviceName: "gmicro",
	}

	for _, o := range opts {
		o(srv)
	}

	srv.Use(mws.TracingHandler(srv.serviceName))
	for _, m := range srv.middlewares {
		mw, ok := mws.Middlewares[m]
		if !ok {
			log.Warnf("can not find middleware: %s", m)
			continue
			//panic(errors.Errorf("can not find middleware: %s", m))
		}

		log.Infof("intall middleware: %s", m)
		srv.Use(mw)
	}

	return srv
}

func (s *Server) Translator() ut.Translator {
	return s.trans
}

// start rest server
func (s *Server) Start(ctx context.Context) error {
	//设置开发模式，打印路由信息
	// 	Gin 有三种模式：
	// debug（开发）
	// release（生产）
	// test（测试）
	if s.mode != gin.DebugMode && s.mode != gin.ReleaseMode && s.mode != gin.TestMode {
		return errors.New("mode must be one of debug/release/test")
	}

	//设置开发模式，打印路由信息
	gin.SetMode(s.mode)
	// Gin 框架的【开发模式 + 自定义路由日志打印】，作用：让 Gin 启动时，用你们自己的日志格式，漂亮打印所有 API 路由
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		// %-6s = 左对齐、固定 6 位宽度字符串
		log.Infof("%-6s %-s --> %s(%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	//TODO 初始化翻译器
	err := s.initTrans(s.transName)
	if err != nil {
		log.Errorf("initTrans error %s", err.Error())
		return err
	}

	//注册mobile验证码
	validation.RegisterMobile(s.trans)

	//根据配置初始化pprof路由
	if s.enableProfiling {
		pprof.Register(s.Engine)
	}

	if s.enableMetrics {
		// get global Monitor object
		m := ginmetrics.GetMonitor()
		// +optional set metric path, default /debug/metrics
		m.SetMetricPath("/metrics")
		// +optional set slow time, default 5s
		// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
		// used to p95, p99
		m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
		m.Use(s)
	}

	log.Infof("rest server is running on port: %d", s.port)
	address := fmt.Sprintf(":%d", s.port)
	s.server = &http.Server{
		Addr:    address,
		Handler: s.Engine,
	}
	// 就是 清空信任列表，不信任任何代理。直接从 TCP 连接拿真实 IP，无法伪造
	// 因为你们是微服务，直接对外提供 API，没有前置代理。为了安全，不信任任何代理。
	_ = s.SetTrustedProxies(nil)
	// 启动 HTTP 服务，只有发生【真正启动失败 / 异常错误】才返回错误；
	// 【主动优雅关闭服务如shutdown】的正常错误，直接忽略。
	// 	启动服务时要忽略这 3 种都叫 优雅关闭时的错误，它门算正常错误可以忽略
	//   - Ctrl+C
	//   - kill <pid>
	//     - kill -9 的算强杀
	//   - 代码主动 Shutdown()
	//   - 它们全部返回：http.ErrServerClosed
	//   - 你们代码里都会忽略这个错误，不打印错误日志，平滑退出。
	// 普通的 s.Run () 的实例没有优雅退出方法，所以改用s.server.ListenAndServe()
	if err = s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Infof("rest server is stopping")
	// 只有这个server对象才有优雅退出功能即Shutdown方法
	// 普通的s.run的实例没有这个方法，所以改用s.server.ListenAndServe()
	// 接收一个 ctx：最长等待多久（超时时间）
	if err := s.server.Shutdown(ctx); err != nil {
		log.Errorf("rest server shutdown error: %s", err.Error())
		return err
	}
	log.Info("rest server stopped")
	return nil
}
