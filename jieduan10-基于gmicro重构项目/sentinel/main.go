package main

// 集成 Alibaba Sentinel-Go：Go 语言流量控制、熔断降级、系统防护组件；
// 对接 Nacos 配置中心：从 Nacos 动态拉取流控规则，规则变更实时生效；
// 启动 10 个常驻协程模拟持续请求流量；
// 通过原子计数器统计 总请求数、通过数、被拦截数，每秒打印监控指标；
// 规则完全托管在 Nacos，无需重启服务即可更新限流策略。
import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	// 阿里 Sentinel Go 客户端
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/ext/datasource"
	"github.com/alibaba/sentinel-golang/pkg/datasource/nacos"

	// Nacos 官方 Go SDK
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
)

type Counter struct {
	pass  *int64
	block *int64
	total *int64
}

func main() {
	//流量计数器,为了流控打印日志更直观,和集成nacos数据源无关。
	counter := Counter{
		pass:  new(int64),
		block: new(int64),
		total: new(int64),
	}

	//nacos server地址
	// Nacos 服务端地址
	sc := []constant.ServerConfig{
		{
			ContextPath: "/nacos",
			Port:        8848,
			IpAddr:      "39.107.30.137",
		},
	}

	//nacos client 相关参数配置,具体配置可参考github.com/nacos-group/nacos-sdk-go
	// Nacos 客户端配置
	cc := constant.ClientConfig{
		NamespaceId: "public",
		TimeoutMs:   5000,
	}
	// 创建 Nacos 配置客户端
	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}

	// ================================== 新知识点 ===================
	// 注册流控规则解析器 + 处理器
	h := datasource.NewFlowRulesHandler(datasource.FlowRuleJsonArrayParser)
	// 创建NacosDataSource数据源 ---- 重要！
	// 前面创建好的 Nacos 配置客户端，负责和 Nacos 服务端通信、拉取配置。
	// "sentinel-go" ---- 对应 Nacos 配置的 Group Name（分组名）。
	// "flow"：对应 Nacos 配置的 Data Id（配置 ID）。组合规则：Group + DataId 唯一定位一条 Nacos 配置
	// h：上一步创建的流控规则处理器。
	nds, err := nacos.NewNacosDataSource(client, "sentinel-go", "flow", h)
	if err != nil {
		panic(err)
	}
	// 	初始化工作：这是触发动作的关键一步，调用后执行完整流程：
	// 首次拉取：立即根据 Group+DataId 从 Nacos 拉取当前配置；
	// 解析 + 加载：走前面的 parser + handler，把规则灌入 Sentinel；
	// 开启长轮询监听：后台常驻协程，持续监听 Nacos 配置变更；
	// 热更新：一旦 Nacos 配置修改发布，自动重复「拉取 → 解析 → 加载」，无需重启服务。
	// 不调用 Initialize()：数据源只是一个空对象，不会拉取任何规则，限流完全不生效
	err = nds.Initialize()
	if err != nil {
		panic(err)
	}

	// =============== 下面就是测试代码，跟上面核心逻辑无关 =========
	// 定时打印看效果
	go timerTask(&counter)

	//模拟流量
	// 开启 10 个常驻协程，无限循环模拟持续请求，不断经过 Sentinel 做限流判断。
	// 启动chanel的目的：阻止用空通道阻塞主线程，防止程序直接退出，这种缺点：不支持主动退出，用waitgroup也可以实现
	ch := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			for {
				// 1. 总请求数 +1
				atomic.AddInt64(counter.total, 1)
				// 2. Sentinel 埋点：进入资源 "test" 链路
				a, b := sentinel.Entry("test", sentinel.WithTrafficType(base.Inbound))
				if b != nil {
					// 被限流/拦截
					atomic.AddInt64(counter.block, 1)
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
				} else {
					// 正常放行
					atomic.AddInt64(counter.pass, 1)
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
					// 必须调用 Exit 释放令牌、上报指标
					a.Exit()
				}
			}
		}()
	}
	<-ch // // 永久阻塞，没人发数据就一直等
}

// 定时统计协程 timerTask
// 每秒执行一次；
// 用 atomic.LoadInt64 原子读取计数；
// 计算每秒增量：总请求、通过、拦截，并打印日志；
// 直观观察限流效果。
func timerTask(counter *Counter) {
	var (
		// 总请求、通过、拦截，
		oldTotal, oldPass, oldBlock int64
	)

	for {
		time.Sleep(time.Second) // 每一秒
		globalTotal := atomic.LoadInt64(counter.total)

		oneSecondTotal := globalTotal - oldTotal
		oldTotal = globalTotal

		globalPass := atomic.LoadInt64(counter.pass)
		oneSecondPass := globalPass - oldPass
		oldPass = globalPass

		globalBlock := atomic.LoadInt64(counter.block)
		oneSecondBlock := globalBlock - oldBlock
		oldBlock = globalBlock

		fmt.Println("total:", oneSecondTotal, "pass:", oneSecondPass, "block:", oneSecondBlock)
	}
}
