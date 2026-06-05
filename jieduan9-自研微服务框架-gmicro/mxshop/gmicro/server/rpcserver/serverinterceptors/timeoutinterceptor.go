package serverinterceptors

// Go 语言编写的 gRPC 一元请求超时拦截器（Unary Timeout Interceptor）
// 直接照搬 go-zero 框架的生产级实现
// 核心作用：为每个 gRPC 一元调用强制设置超时，防止慢请求耗尽服务资源
import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryTimeoutInterceptor
// 传入超时时间，返回一个 gRPC 服务端一元拦截器
func UnaryTimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	// 返回 gRPC 标准的拦截器匿名函数
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		// --------------------------
		// 1. 给当前请求上下文设置超时
		// 超时时间一到，ctx.Done() 通道会被关闭
		// --------------------------
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel() // 确保函数退出时释放上下文资源，防止泄漏

		// 存放业务返回结果
		var resp interface{}
		// 存放业务错误
		var err error
		// 锁：保护 resp 和 err，防止多协程并发读写
		var lock sync.Mutex
		// 通知：业务逻辑正常执行完成
		done := make(chan struct{})
		// 捕获 panic：带缓冲=1，防止协程泄漏（关键）
		panicChan := make(chan interface{}, 1)

		// --------------------------
		// 2. 开协程 执行业务逻辑（handler）
		// 必须异步，否则无法实现超时打断
		// --------------------------
		go func() {
			// 捕获业务代码里的 panic（如空指针、数组越界）
			defer func() {
				if p := recover(); p != nil {
					// 把 panic 信息 + 堆栈一起发出去，方便排查
					panicChan <- fmt.Sprintf("%+v\n\n%s", p, strings.TrimSpace(string(debug.Stack())))
				}
			}()

			// 加锁：安全写入 resp 和 err -------- 只要涉及resp是全局变量，在多个协程中读写，必须加锁
			lock.Lock()
			defer lock.Unlock()

			// --------------------------
			// 真正执行业务方法（Controller 层）
			// --------------------------
			resp, err = handler(ctx, req)
			// 业务执行完毕，通知 done
			close(done)
		}()

		// --------------------------
		// 3. 核心：select 等待 3 种goroutine结果
		// 哪个先到就走哪个分支
		// --------------------------
		select {
		// 情况1：业务代码 panic了
		case p := <-panicChan:
			panic(p) // 把错误抛出去（外层还有崩溃拦截器处理）

		// 情况2：业务正常执行完成
		case <-done:
			lock.Lock()
			defer lock.Unlock()
			return resp, err // 读：直接返回业务结果

		// 情况3：超时了 或 上下文被取消
		case <-ctx.Done():
			err := ctx.Err()
			// 把标准 context 错误 转换成 gRPC 标准错误
			if err == context.Canceled {
				//我们之前说过我们把error统一了， grpc的error我们也可以统一, 自己完成
				err = status.Error(codes.Canceled, err.Error())
			} else if err == context.DeadlineExceeded {
				err = status.Error(codes.DeadlineExceeded, err.Error())
			}
			return nil, err
		}
	}
}
