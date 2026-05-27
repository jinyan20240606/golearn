package main

// Alibaba Sentinel（阿里巴巴限流 / 熔断 / 降级框架） 的 Go 版本
// 作用：给资源设置 QPS 限流，超过阈值直接拒绝 or 匀速排队

import (
	"fmt"
	"log"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
)

func main() {
	//先初始化sentinel
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("初始化sentinel 异常: %v", err)
	}

	//配置限流规则
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test",     // 资源名（你要限流的接口/方法名）
			TokenCalculateStrategy: flow.Direct,     // 直接按阈值计算
			ControlBehavior:        flow.Throttling, //匀速通过（排队，不拒绝）
			Threshold:              100,             //QPS 阈值 = 100，每秒正常扛 100 次请求，1ms处理一个
			StatIntervalInMs:       1000,            // 统计窗口 1 秒
		},

		{
			Resource:               "some-test2",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject, //直接拒绝
			Threshold:              10,          // QPS 阈值 = 10
			StatIntervalInMs:       1000,
		},
	})

	if err != nil {
		log.Fatalf("加载规则失败: %v", err)
	}

	// 模拟 12 次请求（最关键）

	for i := 0; i < 12; i++ {
		// 申请进入一个被限流保护的资源如 “some-test”。通过：b=nil，被限流 / 熔断：b≠nil
		// 就是 真正生产环境的代码，不是测试代码，限流请求就这么用
		e, b := sentinel.Entry("some-test2", sentinel.WithTrafficType(base.Inbound)) // WithTrafficType选项参数，用于设置流量类型，用来给本次请求附加信息
		// base.Inbound ，表示：入站流量
		// Inbound = 别人调用我（服务端接口、gRPC 服务端、HTTP 服务端）
		// Outbound = 我调用别人（客户端调用其他服务）
		if b != nil {
			fmt.Println("限流了")
		} else {
			fmt.Println("检查通过")
			// 检查通过 → 执行业务，业务结束必须调用：e.Exit()
			// 告诉 Sentinel 我执行完了，更新统计数据，不写 Exit () 会导致限流统计不准！
			e.Exit()
		}
		time.Sleep(11 * time.Millisecond)
	}
}
