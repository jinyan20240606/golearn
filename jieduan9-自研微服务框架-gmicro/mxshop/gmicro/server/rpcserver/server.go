package rpcserver

import (
	"context"

	// 自定义服务端拦截器实现的包
	srvintc "mxshop/gmicro/server/rpcserver/serverinterceptors"
	"mxshop/pkg/host"
	"mxshop/pkg/log"
	"net"
	"net/url"
	"time"

	// gRPC 全链路追踪（Trace）官方拦截器
	// 	自动生成 TraceID / SpanID
	// 自动在 gRPC 请求里传递链路上下文
	// 把调用数据上报给 Jaeger / Zipkin / 阿里云链路追踪
	// 框架必须有，不能删
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	// 自动注册 health 接口 (gRPC，http)
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	// 作用：gRPC 反射服务
	// 功能：
	// 让外部工具（grpcurl、grpcui、postman）不用 proto 文件也能访问 gRPC
	// 方便测试、调试
	// 生产环境可开可关
	"google.golang.org/grpc/reflection"

	// 拷贝的kratos的metadata功能
	apimd "mxshop/api/metadata"
)

type ServerOption func(o *Server)

type Server struct {
	*grpc.Server

	// 监听地址：ip:port格式，默认 :0（随机端口）
	address    string
	unaryInts  []grpc.UnaryServerInterceptor
	streamInts []grpc.StreamServerInterceptor
	grpcOpts   []grpc.ServerOption
	// 开启监听端口服务的实例：如lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))的lis值
	lis     net.Listener
	timeout time.Duration

	health   *health.Server
	metadata *apimd.Server
	// 这个地址会被注册到 Consul /etcd，让其他服务能发现并调用。
	endpoint *url.URL

	enableMetrics bool
}

func (s *Server) Endpoint() *url.URL {
	return s.endpoint
}

func (s *Server) Address() string {
	return s.address
}

// 入口函数：NewServer 创建一个Server，支持函数选项模式
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		address: ":0",
		health:  health.NewServer(),
		//timeout: 1 * time.Second, // 这块不能设默认值，要动态传入
	}

	for _, o := range opts {
		o(srv)
	}

	//TODO 我们现在希望用户不设置拦截器的情况下，我们会自动默认加上一些必须的拦截器， crash，tracing
	unaryInts := []grpc.UnaryServerInterceptor{
		srvintc.UnaryCrashInterceptor,     // 必选：panic 崩溃恢复
		otelgrpc.UnaryServerInterceptor(), // 必选：opentelemetry 全链路追踪
	}
	// 如果开启监控，自动加入 Prometheus 监控拦截器
	if srv.enableMetrics {
		unaryInts = append(unaryInts, srvintc.UnaryPrometheusInterceptor)
	}
	// 如果配置了超时，自动加入超时控制拦截器
	if srv.timeout > 0 {
		unaryInts = append(unaryInts, srvintc.UnaryTimeoutInterceptor(srv.timeout))
	}
	// 最后：把【用户自己定义的拦截器】追加到后面
	if len(srv.unaryInts) > 0 {
		unaryInts = append(unaryInts, srv.unaryInts...)
	}

	//把我们传入的拦截器转换成grpc的ServerOption，把所有拦截器打包成 grpc 的链式拦截器
	grpcOpts := []grpc.ServerOption{grpc.ChainUnaryInterceptor(unaryInts...)}

	//把用户自己传入的grpc.ServerOption放在一起
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}

	// grpc的服务初始化方法
	srv.Server = grpc.NewServer(grpcOpts...)

	//注册metadata的Server
	srv.metadata = apimd.NewServer(srv.Server)

	//自动解析address，提取出可用的address
	err := srv.listenAndEndpoint()
	if err != nil {
		panic(err)
	}

	// 给 gRPC 服务注册 3 个官方 / 框架内置系统服务
	// 1、注册 gRPC 官方标准健康检查服务
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	// 2、注册服务元数据服务：可以查询当前服务信息
	apimd.RegisterMetadataServer(srv.Server, srv.metadata)
	// 3、开启 gRPC 反射服务（最关键、最常用）：支持用户直接通过 grpc 的一个http接口不需要proto文件查看当前支持的所有的 rpc 服务
	reflection.Register(srv.Server)

	return srv
}

func WithAddress(address string) ServerOption {
	return func(s *Server) {
		s.address = address
	}
}

func WithMetrics(metric bool) ServerOption {
	return func(s *Server) {
		s.enableMetrics = metric
	}
}

func WithTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func WithLis(lis net.Listener) ServerOption {
	return func(s *Server) {
		s.lis = lis
	}
}

func WithUnaryInterceptor(in ...grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.unaryInts = in
	}
}

func WithStreamInterceptor(in ...grpc.StreamServerInterceptor) ServerOption {
	return func(s *Server) {
		s.streamInts = in
	}
}

func WithOptions(opts ...grpc.ServerOption) ServerOption {
	return func(s *Server) {
		s.grpcOpts = opts
	}
}

// 完成ip和端口的提取
// 启动 gRPC 监听 + 生成服务注册用的最终地址（IP:Port）
func (s *Server) listenAndEndpoint() error {
	// 如果还没监听，就根据配置的 address（如 :8081、0.0.0.0:8082）启动 TCP 监听
	// 监听成功后，把 net.Listener 保存起来，供后面提取出监听服务启动后真实的 ip:port，供gRPC 使用
	if s.lis == nil {
		lis, err := net.Listen("tcp", s.address)
		if err != nil {
			return err
		}
		s.lis = lis
	}
	// 获取最终地址
	// 如果是 0.0.0.0 → 获取本机真实网卡 IP
	//如果是 :port → 自动补全 IP
	//返回 真实可被访问的 ip:port
	addr, err := host.Extract(s.address, s.lis)
	if err != nil {
		_ = s.lis.Close()
		return err
	}
	// 生成 endpoint：grpc://192.168.1.103:8081
	// 这个地址会被注册到 Consul /etcd，让其他服务能发现并调用。
	s.endpoint = &url.URL{Scheme: "grpc", Host: addr}
	return nil
}

// 启动grpc的服务
func (s *Server) Start(ctx context.Context) error {
	log.Infof("[grpc] server listening on: %s", s.lis.Addr().String())
	s.health.Resume()
	return s.Serve(s.lis)
}

func (s *Server) Stop(ctx context.Context) error {
	//设置服务的状态为not_serving，防止接收新的请求过来
	s.health.Shutdown()
	s.GracefulStop()
	log.Infof("[grpc] server stopped")
	return nil
}
