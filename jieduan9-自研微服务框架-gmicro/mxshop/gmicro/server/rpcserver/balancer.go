package rpcserver

// 这段代码 = 把你们自己写的 【selector 负载均衡算法】
// 完美包装成 gRPC 官方要求的 Picker
// 让 gRPC 原生框架可以直接使用
import (
	"mxshop/gmicro/registry"

	"mxshop/gmicro/server/rpcserver/selector"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/metadata"
)

// 给 gRPC 内部注册一个名字叫 selector 的负载均衡器。
const (
	balancerName = "selector"
)

var (
	_ base.PickerBuilder = &balancerBuilder{}
	_ balancer.Picker    = &balancerPicker{}
)

// 告诉 gRPC：
// 我有一个自定义负载均衡器，名字叫 selector
// 以后你遇到 loadBalancingPolicy": "selector"
// 就用我这个！
// 这就是 注册自定义算法
func InitBuilder() {
	b := base.NewBalancerBuilder(
		balancerName,
		&balancerBuilder{
			// 设置具体算法策略实例
			builder: selector.GlobalSelector(),
		},
		base.Config{HealthCheck: true},
	)
	balancer.Register(b)
}

// 它的作用只有一个：
// 把服务发现拿到的节点列表 → 交给你们自己的 selector 算法
type balancerBuilder struct {
	builder selector.Builder
}

// 把节点变成你们的 Node
// 作用：
// gRPC 每次节点列表变化（新增/下线/健康状态变化）
// 就会调用这个 Build 方法，创建一个新的 Picker（选择器）
func (b *balancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	// ---------------------- 第1行：判断有没有可用连接 ----------------------
	// ReadySCs = 已经准备好的、可以正常使用的 gRPC 子连接
	// 如果一个可用连接都没有
	if len(info.ReadySCs) == 0 {
		// 返回一个错误选择器：告诉 gRPC 现在没有可用节点
		// Block the RPC until a new picker is available via UpdateState().
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	// ---------------------- 第2行：创建节点列表 ----------------------
	// 把 gRPC 底层连接 → 转换成你们业务层的 selector.Node
	nodes := make([]selector.Node, 0, len(info.ReadySCs))
	// ---------------------- 第3行：遍历所有可用连接 ----------------------
	// conn = gRPC 底层连接（SubConn）
	// info = 这个连接对应的地址信息
	for conn, info := range info.ReadySCs {
		// ---------------------- 第4行：取出服务实例信息 ----------------------
		// 从地址属性里取出 注册中心（Consul）拿到的完整服务实例
		// 包含：服务名、版本、权重、机房、地址等
		ins, _ := info.Address.Attributes.Value("rawServiceInstance").(*registry.ServiceInstance)
		// ---------------------- 第5行：组装成你们自己的节点 ----------------------
		// 重点：
		// grpcNode 组合了两个东西：
		// 1. selector.Node（业务层：地址、权重、实例信息）
		// 2. subConn（gRPC底层真实TCP连接）
		nodes = append(nodes, &grpcNode{
			Node:    selector.NewNode("grpc", info.Address.Addr, ins),
			subConn: conn,
		})
	}
	// ---------------------- 第6行：创建 Picker（真正做负载均衡） ----------------------
	// b.builder 是你们设置的：random / round / hash 算法
	// b.builder.Build() = 构建一个选择器实例（随机/轮询/一致性哈希）
	p := &balancerPicker{
		selector: b.builder.Build(),
	}
	// ---------------------- 第7行：把节点列表喂给负载均衡器 ----------------------
	// 让你们的 selector 知道：现在有哪些节点可以选
	p.selector.Apply(nodes)
	// ---------------------- 第8行：返回给 gRPC ----------------------
	// 返回这个 Picker
	// 之后每次请求，gRPC 都会调用 p.Pick()
	return p
}

// balancerPicker：真正做选择的地方（负载均衡本体）
type balancerPicker struct {
	// selector/default_selector.go 定义的
	selector selector.Selector
}

// Pick pick instances.
func (p *balancerPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	// 【真正负载均衡在这里！】--- 就是你们自己写的 selector.Select 选择算法！
	n, done, err := p.selector.Select(info.Ctx)
	if err != nil {
		return balancer.PickResult{}, err
	}

	return balancer.PickResult{
		SubConn: n.(*grpcNode).subConn,
		Done: func(di balancer.DoneInfo) {
			done(info.Ctx, selector.DoneInfo{
				Err:           di.Err,
				BytesSent:     di.BytesSent,
				BytesReceived: di.BytesReceived,
				// 让你们的 selector 可以拿到请求完成后的响应头信息
				ReplyMD: Trailer(di.Trailer),
			})
		},
	}, nil
}

// 在 Pick() 的 Done 回调里：
// 让你们的 selector 可以拿到请求完成后的响应头信息
// Trailer is a grpc trailder MD.
type Trailer metadata.MD

// Get get a grpc trailer value.
func (t Trailer) Get(k string) string {
	v := metadata.MD(t).Get(k)
	if len(v) > 0 {
		return v[0]
	}
	return ""
}

// 把你们的业务节点 + gRPC 底层连接绑在一起
type grpcNode struct {
	selector.Node
	subConn balancer.SubConn
}
