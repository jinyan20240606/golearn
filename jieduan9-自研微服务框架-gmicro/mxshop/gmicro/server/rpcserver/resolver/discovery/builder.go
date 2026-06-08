package discovery

import (
	"context"
	"errors"
	"strings"
	"time"

	"mxshop/gmicro/registry"

	"google.golang.org/grpc/resolver"
)

const name = "discovery"

// Option is builder option.
type Option func(o *builder)

// WithTimeout with timeout option.
func WithTimeout(timeout time.Duration) Option {
	return func(b *builder) {
		b.timeout = timeout
	}
}

// WithInsecure with isSecure option.
func WithInsecure(insecure bool) Option {
	return func(b *builder) {
		b.insecure = insecure
	}
}

type builder struct {
	discoverer registry.Discovery
	timeout    time.Duration
	insecure   bool
}

// NewBuilder creates a builder which is used to factory registry resolvers.
func NewBuilder(d registry.Discovery, opts ...Option) resolver.Builder {
	b := &builder{
		discoverer: d,
		timeout:    time.Second * 10,
		insecure:   false,
	}
	for _, o := range opts {
		o(b)
	}
	return b
}

// 这是 gRPC resolver.Builder 的 Build 方法实现
// 作用：创建一个监听注册中心的解析器
// = 连接 如Consul -> 监听服务变化 -> 把地址喂给 gRPC
func (b *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var (
		err error
		w   registry.Watcher
	)
	done := make(chan struct{}, 1)
	ctx, cancel := context.WithCancel(context.Background())
	// 1. 启动协程异步创建 Watcher（监听服务变化）
	// 调用传参进来的Consul/Nacos发现器discoverer， 中创建一个监听器 Watcher
	// 一旦服务列表变化，立刻通知客户端
	go func() {
		// 这块就是观察者模式，调用一个Watch得到一个观察者监听器即w方法，后续有变化时，立刻通知w方法
		w, err = b.discoverer.Watch(ctx, strings.TrimPrefix(target.URL.Path, "/"))
		close(done)
	}()
	// 2. 如果 3s/5s 内连不上注册中心，直接报错，不卡死。
	select {
	case <-done:
	case <-time.After(b.timeout):
		err = errors.New("discovery create watcher overtime")
	}
	if err != nil {
		cancel()
		return nil, err
	}
	// 3. 创建 discoveryResolver（真正的服务发现解析器）
	r := &discoveryResolver{
		w:        w,
		cc:       cc,
		ctx:      ctx,
		cancel:   cancel,
		insecure: b.insecure,
	}
	// 4. 启动后台协程，持续监听服务变化
	go r.watch()
	return r, nil
}

// Scheme return scheme of discovery
func (*builder) Scheme() string {
	return name
}
