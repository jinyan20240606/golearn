package rpcserver

import (
	"context"
	"mxshop/gmicro/server/rpcserver/clientinterceptors"
	"mxshop/gmicro/server/rpcserver/resolver/discovery"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	grpcinsecure "google.golang.org/grpc/credentials/insecure"

	"mxshop/gmicro/registry"
	"mxshop/pkg/log"
	"time"

	"google.golang.org/grpc"
)

type ClientOption func(o *clientOptions)

type clientOptions struct {
	endpoint string
	timeout  time.Duration
	// discovery接口服务发现的方法
	discovery     registry.Discovery
	unaryInts     []grpc.UnaryClientInterceptor
	streamInts    []grpc.StreamClientInterceptor
	rpcOpts       []grpc.DialOption
	balancerName  string
	log           log.LogHelper
	enableTracing bool
	enableMetrics bool
}

func WithEnableTracing(enable bool) ClientOption {
	return func(o *clientOptions) {
		o.enableTracing = enable
	}
}

// 设置地址
func WithEndpoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

// 设置超时时间
func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}

// 设置服务发现
func WithDiscovery(d registry.Discovery) ClientOption {
	return func(o *clientOptions) {
		o.discovery = d
	}
}

// 设置拦截器
func WithClientUnaryInterceptor(in ...grpc.UnaryClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.unaryInts = in
	}
}

// 设置stream拦截器
func WithClientStreamInterceptor(in ...grpc.StreamClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.streamInts = in
	}
}

// 设置grpc的dial选项
func WithClientOptions(opts ...grpc.DialOption) ClientOption {
	return func(o *clientOptions) {
		o.rpcOpts = opts
	}
}

// 设置负载均衡器
func WithBalancerName(name string) ClientOption {
	return func(o *clientOptions) {
		o.balancerName = name
	}
}

// 创建一个不安全的grpc拨号方法
func DialInsecure(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, true, opts...)
}

// 创建一个安全的grpc拨号方法
func Dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, false, opts...)
}

func dial(ctx context.Context, insecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
	options := clientOptions{
		timeout:       2000 * time.Millisecond,
		balancerName:  "round_robin",
		enableTracing: true,
	}

	for _, o := range opts {
		o(&options)
	}

	//TODO 客户端默认拦截器
	ints := []grpc.UnaryClientInterceptor{
		clientinterceptors.TimeoutInterceptor(options.timeout),
	}
	if options.enableTracing {
		ints = append(ints, otelgrpc.UnaryClientInterceptor())
	}

	if options.enableMetrics {
		ints = append(ints, clientinterceptors.PrometheusInterceptor())
	}

	streamInts := []grpc.StreamClientInterceptor{}

	if len(options.unaryInts) > 0 {
		ints = append(ints, options.unaryInts...)
	}
	if len(options.streamInts) > 0 {
		streamInts = append(streamInts, options.streamInts...)
	}

	grpcOpts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "` + options.balancerName + `"}`),
		grpc.WithChainUnaryInterceptor(ints...),
		grpc.WithChainStreamInterceptor(streamInts...),
	}

	//TODO 服务发现的选项
	if options.discovery != nil {
		// 给【本次 Dial 客户端】单独注册一个 resolver 解析器
		// 		只给当前这一个 client 连接使用
		// 可以动态传入注册中心客户端
		// 不污染全局
		grpcOpts = append(grpcOpts, grpc.WithResolvers(
			discovery.NewBuilder(
				options.discovery,
				// 我要连接的是【不安全的 gRPC 服务（没有 TLS）】，不要尝试加密
				discovery.WithInsecure(insecure),
			),
		))
	}

	// 	这两个 insecure的逻辑 长得一样，作用完全不同
	// 一个给 服务发现 用
	// 一个给 gRPC 连接 用

	// 是否安全
	if insecure {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(grpcinsecure.NewCredentials()))
	}

	if len(options.rpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.rpcOpts...)
	}

	return grpc.DialContext(ctx, options.endpoint, grpcOpts...)
}
