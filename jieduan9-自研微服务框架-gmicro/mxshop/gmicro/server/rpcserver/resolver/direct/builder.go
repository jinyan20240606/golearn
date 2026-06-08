package direct

// gRPC 直连发现器（direct resolver）

import (
	"strings"

	"google.golang.org/grpc/resolver"
)

// 包初始化时：自动把 direct 解析器注册到 gRPC 全局
// init 函数里注册，全局生效,不能传参数（比如你的注册中心客户端）
// 适合 direct:// 这种简单的, 所有 grpc.Dial 都能用 direct://
func init() {
	// 把 direct 这个解析器，注册到 gRPC 内部的 全局 map 里！
	resolver.Register(NewBuilder())
}

// directBuilder：实现 gRPC resolver.Builder 接口
type directBuilder struct{}

// NewBuilder
// 创建一个 direct 解析器构造器
// 客户端使用格式：direct:///127.0.0.1:8080,127.0.0.1:8081
func NewBuilder() *directBuilder {
	return &directBuilder{}
}

// Build
// 【核心方法】
// 解析 gRPC 目标地址，把 URL 里的 IP 列表拆出来，传给 gRPC 客户端
func (d *directBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// 从 URL Path 中取出地址，按逗号分割成多个地址
	// 例如 target.URL.Path = "/127.0.0.1:9000,127.0.0.1:9001"
	addrs := make([]resolver.Address, 0)
	for _, addr := range strings.Split(strings.TrimPrefix(target.URL.Path, "/"), ",") {
		addrs = append(addrs, resolver.Address{Addr: addr})
	}

	//grpc建立连接的逻辑都是这里UpdateState
	// 【关键】把解析好的地址列表更新给 gRPC 客户端
	// gRPC 内部靠这个建立连接
	err := cc.UpdateState(resolver.State{Addresses: addrs})
	if err != nil {
		return nil, err
	}
	// 返回一个空 resolver（直连不需要监听变化）
	return newDirectResolver(), nil
}

// Scheme
// 协议名称：客户端使用 direct:// 开头
func (d *directBuilder) Scheme() string {
	return "direct"
}

// 编译期检查：确保 directBuilder 实现了 resolver.Builder 接口
var _ resolver.Builder = &directBuilder{}
