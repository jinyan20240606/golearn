package srv

import (
	"mxshop/app/pkg/options"
	"mxshop/app/user/srv/config"
	gapp "mxshop/gmicro/app"
	"mxshop/gmicro/server/rpcserver"
	"mxshop/pkg/app"
	"mxshop/pkg/log"

	"github.com/google/wire"
	"github.com/hashicorp/consul/api"

	"mxshop/gmicro/registry"
	"mxshop/gmicro/registry/consul"
)

var ProviderSet = wire.NewSet(NewUserApp, NewRegistrar, NewUserRPCServer, NewNacosDataSource)

func NewApp(basename string) *app.App {
	cfg := config.New()
	appl := app.NewApp("user",
		"mxshop",
		app.WithOptions(cfg),
		app.WithRunFunc(run(cfg)),
		//app.WithNoConfig(), //设置不读取配置文件
	)
	return appl
}

func NewRegistrar(registry *options.RegistryOptions) registry.Registrar {
	c := api.DefaultConfig()
	c.Address = registry.Address
	c.Scheme = registry.Scheme
	cli, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}
	// consul客户端实例传进consul方法
	r := consul.New(cli, consul.WithHealthCheck(true))
	return r
}

// NewUserApp
// 作用：初始化并创建 User 服务的应用实例（组装所有依赖：日志、服务、注册中心）
// 入参：日志配置、服务注册器、服务配置、gRPC服务器
// 返回值：创建好的应用实例、错误
func NewUserApp(logOpts *log.Options, register registry.Registrar,
	serverOpts *options.ServerOptions, rpcServer *rpcserver.Server) (*gapp.App, error) {
	//初始化log
	// 1. 初始化日志系统（根据配置加载日志）
	log.Init(logOpts)
	// 延迟刷新日志：防御性编程
	// 很多日志框架为了性能，默认会把日志先存在内存缓冲区里。如果在初始化阶段程序就崩溃退出了，那些还没来得及刷新到磁盘的报错日志策略就会丢失，导致排查问题时毫无头绪。
	defer log.Flush() // 启动正常后，后续服务运行中的日志也有自动flush的策略
	// 2. 组装并创建应用实例 gapp.App
	// 作用：把 服务名、RPC服务、注册中心 全部绑定到一个应用实例里统一管理
	return gapp.New(
		gapp.WithName(serverOpts.Name), // 设置服务名称（用于注册中心、监控）
		gapp.WithRPCServer(rpcServer),  // 绑定gRPC服务
		gapp.WithRegistrar(register),   // 绑定服务注册器（启动时自动注册到注册中心）
	), nil
}

func run(cfg *config.Config) app.RunFunc {
	return func(baseName string) error {
		// 这是wire重构之前的用法，现在使用wire重构成直接使用initApp
		// userApp, err := NewUserApp(cfg.Nacos, cfg.Log, cfg.Server, cfg.Registry, cfg.Telemetry, cfg.MySQLOptions)
		userApp, err := initApp(cfg.Nacos, cfg.Log, cfg.Server, cfg.Registry, cfg.Telemetry, cfg.MySQLOptions)
		if err != nil {
			return err
		}

		//启动
		if err := userApp.Run(); err != nil {
			log.Errorf("run user app error: %s", err)
			return err
		}
		return nil
	}
}
