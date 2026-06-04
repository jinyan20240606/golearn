package registry

import "context"

// 服务注册接口
// gmicro 不把注册中心写死成 consul，它只规定：
// 你给我一个注册器，我在启动和停止时调用它。
type Registrar interface {
	//注册
	Register(ctx context.Context, service *ServiceInstance) error
	//注销
	Deregister(ctx context.Context, service *ServiceInstance) error
}

// 服务发现接口
type Discovery interface {
	//获取服务实例：通过serviceName去发现的
	GetService(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
	//创建服务监听器：返回的是Watcher接口方法
	Watch(ctx context.Context, serviceName string) (Watcher, error)
}

type Watcher interface {
	//获取服务实例, next在下面的情况下会返回服务
	//1. 第一次监听时，如果服务实例列表不为空，则返回服务实例列表
	//2. 如果服务实例发生变化，则返回服务实例列表
	//3. 如果上面两种情况都不满足，则会阻塞到context deadline或者cancel
	Next() ([]*ServiceInstance, error)
	//主动放弃监听
	Stop() error
}

// 注册到注册中心的服务实例
type ServiceInstance struct {
	//注册到注册中心的服务id
	ID string `json:"id"`

	//服务名称
	Name string `json:"name"`

	//服务版本
	Version string `json:"version"`

	//服务元数据
	Metadata map[string]string `json:"metadata"`

	// endpoints是个数组，代表当前服务实例的地址信息
	// 因为一个实例可能不止一种入口：
	// grpc://10.0.0.12:9000 专门标识给grpc服务用
	// http://10.0.0.12:8000 专门标识给http服务用
	Endpoints []string `json:"endpoints"`
}
