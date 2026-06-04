package app

import (
	"context"
	"net/url"
	"syscall"
	"time"

	"github.com/google/uuid"
	// 同时跑一组 goroutine，并统一收集错误、联动取消。它可以看成是 sync.WaitGroup 的增强版
	"golang.org/x/sync/errgroup"

	"mxshop/gmicro/registry"
	gs "mxshop/gmicro/server"
	"mxshop/pkg/log"
	"os"
	"os/signal"
	"sync"
)

// App 是微服务的应用启动器，统一管理当前服务实例信息、服务注册与优雅启停流程。
// 对外通常只暴露创建入口和运行入口，其余方法作为内部生命周期管理细节。
type App struct {
	// opts 保存函数选项模式注入的配置。
	opts options

	// lk 用于保护 instance 的并发读写，避免 Run 和 Stop 多协程情况下并发读写相同变量时出现竞态。
	lk sync.Mutex
	// instance 表示当前服务在注册中心中的实例信息。
	instance *registry.ServiceInstance

	// cancel 用于取消整个应用运行上下文，触发所有 server 停止。
	cancel func()
}

// New 创建一个 App，并初始化默认配置。
// 默认配置参数行为包括：
// 1. 监听常见退出信号；
// 2. 设置注册中心超时和停止超时；
// 3. 自动生成服务实例 ID；
// 4. 再使用用户传入的 Option 覆盖默认值。
func New(opts ...Option) *App {
	o := options{
		// 默认监听常见退出信号
		sigs:             []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		registrarTimeout: 10 * time.Second,
		stopTimeout:      10 * time.Second,
	}

	// 为当前进程生成一个默认实例 ID，便于服务注册中心区分不同实例。
	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}

	// 应用用户传入的配置项，覆盖默认配置。
	for _, opt := range opts {
		opt(&o)
	}

	return &App{
		opts: o,
	}
}

// Run 启动当前微服务实例的应用生命周期。
// 核心流程：
// 1. 构建服务实例信息；
// 2. 启动传进来的restServer和rpcServer： 如启动grpc的服务和web服务；
// 3. 全部启动后将他们注册到注册中心；
// 4. 监听退出信号；
// 5. 当任一 server 异常或收到退出信号时，触发整体停止。
func (a *App) Run() error {
	// 先构建注册到注册中心需要的当前实例信息。
	instance, err := a.buildInstance()
	if err != nil {
		return err
	}

	// 保存实例信息，供 Stop 阶段反注册使用。--- 需要加锁
	a.lk.Lock()
	a.instance = instance
	a.lk.Unlock()

	//if a.opts.rpcServer != nil {
	//	// 启动rpc服务， 如果我想要给这个rpc服务设置port 我们想要给这个rpc服务register我们自定义的interceptor
	//	a.opts.rpcServer.Serve()
	//}

	//重点， 写的很简单， http服务要启动
	//if a.opts.rpcServer != nil {
	//	err := a.opts.rpcServer.Start()
	//	if err != nil {
	//		return err
	//	}
	//}

	//现在启动了两个server，一个是restserver，一个是rpcserver
	/*
		这两个server是否必须同时启动成功？
		如果有一个启动失败，那么我们就要停止另外一个server
		如果启动了多个， 如果其中一个启动失败，其他的应该被取消
			如果剩余的server的状态：
				1. 还没有开始调用start
					stop
				2. start进行中
					调用进行中的cancel
				3. start已经完成
					调用stop
		如果我们的服务启动了然后这个时候用户立马进行了访问
	*/

	var servers []gs.Server
	if a.opts.restServer != nil {
		servers = append(servers, a.opts.restServer)
	}
	if a.opts.rpcServer != nil {
		servers = append(servers, a.opts.rpcServer)
	}

	// 构建应用级上下文：
	// 任意一个 goroutine 返回错误，errgroup 会取消 ctx，
	// 其他 server 会收到停止信号并进入 Stop 流程。
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	eg, ctx := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}

	for _, srv := range servers {
		srv := srv

		// 每个 server 都对应一个“停止协程”：
		// 一旦应用 ctx 被取消，就调用 server 的 Stop 完成优雅停机。
		eg.Go(func() error {
			<-ctx.Done()
			sctx, cancel := context.WithTimeout(context.Background(), a.opts.stopTimeout)
			defer cancel()
			return srv.Stop(sctx)
		})

		// 启动 server。
		// 这里配合 WaitGroup 的目的，是确保所有 Start goroutine 都已真正被调度出去，
		// 然后再继续执行后续的服务注册逻辑。
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			log.Info("start rest server")
			return srv.Start(ctx)
		})
	}

	// 等待所有启动协程进入运行状态后，再进行注册。
	wg.Wait()

	// 所有 server 启动流程已发起后，将当前实例注册到注册中心。
	// 这样外部服务发现系统就能感知到该服务实例。
	if a.opts.registrar != nil {
		rctx, rcancel := context.WithTimeout(context.Background(), a.opts.registrarTimeout)
		defer rcancel()
		err := a.opts.registrar.Register(rctx, instance)
		if err != nil {
			log.Errorf("register service error: %s", err)
			return err
		}
	}

	// 监听系统退出信号，实现优雅退出。
	// 收到信号后调用 Stop，先反注册，再取消 ctx，触发全部 server 停止。
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	// 把一个新的 goroutine 加入 errgroup 统一管理
	eg.Go(func() error {
		// 		这里启动一个并发任务
		// 这个任务将来返回的 error，会被 eg.Wait() 收集
		// 如果返回非 nil error，会影响整个应用退出流程
		select { // select 是 Go 里专门配合 channel 使用的并发语法，作用是：同时等待多个 channel 操作，谁先就绪就执行谁
		// 监听应用级 context 的取消信号
		case <-ctx.Done():
			return ctx.Err()
		// 监听系统信号 channel c
		case <-c:
			// 收到退出信号后，执行应用停止逻辑
			return a.Stop()
		}
	})

	// 阻塞等待所有 goroutine 结束。
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

// Stop 停止整个应用。
// 执行顺序是：
// 1. 从注册中心反注册，避免外部继续请求该实例；
// 2. 取消应用上下文；
// 3. 由各个 server 的停止协程调用具体 Stop 完成收尾。
func (a *App) Stop() error {
	// 由于这里是放在一个独立的 goroutine 里执行的--所有就有并发竞争问题，所以需要加锁保护 instance 变量，避免与 Run 中的并发读写冲突。
	a.lk.Lock()
	instance := a.instance
	a.lk.Unlock()

	log.Info("start deregister service")
	if a.opts.registrar != nil && instance != nil {
		rctx, rcancel := context.WithTimeout(context.Background(), a.opts.stopTimeout)
		defer rcancel()
		if err := a.opts.registrar.Deregister(rctx, instance); err != nil {
			log.Errorf("deregister service error: %s", err)
			return err
		}
	}

	// 取消应用上下文，通知所有依赖 ctx 的 goroutine 开始退出。
	if a.cancel != nil {
		a.cancel()
	}

	return nil
}

// buildInstance： 收集当前服务实例的信息，组装成 registry.ServiceInstance，后面交给注册中心去注册。。
func (a *App) buildInstance() (*registry.ServiceInstance, error) {
	endpoints := make([]string, 0)
	// endpoints是个数组，代表当前服务实例的地址信息
	// 因为一个实例可能不止一种入口：
	// grpc://10.0.0.12:9000 专门标识给grpc服务用
	// http://10.0.0.12:8000 专门标识给http服务用
	for _, e := range a.opts.endpoints {
		endpoints = append(endpoints, e.String())
	}

	// 如果配置了 RPC 服务，则优先从 RPC server 自身提取 endpoint。
	// 若 server 没有返回完整 Endpoint，则退化为根据 Address 拼接一个 grpc:// 地址。
	if a.opts.rpcServer != nil {
		if a.opts.rpcServer.Endpoint() != nil {
			endpoints = append(endpoints, a.opts.rpcServer.Endpoint().String())
		} else {
			u := &url.URL{
				Scheme: "grpc",
				Host:   a.opts.rpcServer.Address(),
			}
			endpoints = append(endpoints, u.String())
		}
	}

	return &registry.ServiceInstance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		Endpoints: endpoints,
	}, nil
}
