package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/alibaba/sentinel-golang/util"
)

// 自定义状态监听器：用来打印熔断器 关闭/打开/半开 状态变化
type stateChangeTestListener struct {
}

// 这3个方法是可选的，只是给你看下状态转移的日志查看
// 切换到 关闭状态（正常）
func (s *stateChangeTestListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.steategy: %+v, From %s to Closed, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

// 切换到 打开状态（熔断，全部拦截）
func (s *stateChangeTestListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	fmt.Printf("rule.steategy: %+v, From %s to Open, snapshot: %d, time: %d\n", rule.Strategy, prev.String(), snapshot, util.CurrentTimeMillis())
}

// 切换到 半开状态（探测恢复）
func (s *stateChangeTestListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.steategy: %+v, From %s to Half-Open, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func main() {
	// 全局统计
	total := 0      // 总请求数
	totalPass := 0  // 通过数
	totalBlock := 0 // 被熔断拦截数
	totalErr := 0   // 业务报错数

	// 1. 初始化 Sentinel
	conf := config.NewDefaultConfig()
	// 配置里主要加了个自定义日志的方法
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger() // 日志输出到控制台
	err := sentinel.InitWithConfig(conf)
	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan struct{})
	// 2. 注册状态监听器 → 打印熔断状态变化
	circuitbreaker.RegisterStateChangeListeners(&stateChangeTestListener{})
	// 3. 加载【错误数熔断规则】⭐⭐⭐核心
	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		// Statistic time span=10s, recoveryTimeout=3s, maxErrorCount=50
		{
			Resource:         "abc",                     // 资源名
			Strategy:         circuitbreaker.ErrorCount, // 熔断策略：错误数策略
			RetryTimeoutMs:   3000,                      // 熔断3秒后 → 进入半开探测
			MinRequestAmount: 10,                        // 静默请求数：至少10个请求，才开始判断熔断
			StatIntervalMs:   5000,                      // 统计窗口：5秒内统计错误数
			Threshold:        50,                        // 错误数阈值：5秒内错误≥50 → 熔断
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	logging.Info("[CircuitBreaker ErrorCount] Sentinel Go circuit breaking demo is running. You may see the pass/block metric in the metric log.")
	// 4. 协程1：无限请求，模拟 50% 业务错误
	go func() {
		for {
			total++
			e, b := sentinel.Entry("abc")
			if b != nil {
				// g1 blocked
				totalBlock++
				fmt.Println("协程熔断了")
				// 生成一个 0 ~ 极大的随机无符号整数（随机 64 位数字）
				time.Sleep(time.Duration(rand.Uint64()%20) * time.Millisecond)
			} else {
				totalPass++
				// 计算出0到19的整数 制造50%概率
				if rand.Uint64()%20 > 9 {
					totalErr++
					// 主动抛出错误
					sentinel.TraceError(e, errors.New("biz error"))
				}
				// g1 passed
				time.Sleep(time.Duration(rand.Uint64()%20+10) * time.Millisecond)
				e.Exit()
			}
		}
	}()
	// 5. 协程2：无限请求，模拟正常请求（无错误）
	go func() {
		for {
			total++
			e, b := sentinel.Entry("abc")
			if b != nil {
				// g2 blocked
				totalBlock++
				time.Sleep(time.Duration(rand.Uint64()%20) * time.Millisecond)
			} else {
				// g2 passed
				totalPass++
				time.Sleep(time.Duration(rand.Uint64()%80) * time.Millisecond)
				e.Exit()
			}
		}
	}()
	// 6. 每秒打印错误总数
	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Println(totalErr)
		}
	}()
	// 阻塞主程序，不退出
	<-ch
}
