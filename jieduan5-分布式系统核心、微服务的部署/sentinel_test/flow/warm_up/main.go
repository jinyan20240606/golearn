package main

import (
	"fmt"
	"log"
	"math/rand"
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
	// 全局统计变量（多协程共享）
	var globalTotal int       // 所有协程 总共发起的请求数
	var passTotal int         // 所有协程 限流通过的请求数
	var blockTotal int        // 所有协程 被限流拒绝的请求数
	ch := make(chan struct{}) // 阻塞channel，让主程序不退出

	//配置限流规则，核心：冷启动 WarmUp 限流
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test",
			TokenCalculateStrategy: flow.WarmUp, //冷启动策略
			ControlBehavior:        flow.Reject, //直接拒绝
			Threshold:              1000,        // qps阈值仍然是1000
			WarmUpPeriodSec:        30,          // 预热的时长秒，仅对冷启动策略生效
			// WarmUp 说明：
			// 刚启动时 QPS 很小 → 30秒后 逐渐放开到 1000 QPS
			// 防止服务刚启动、缓存未准备好，被流量直接冲垮
		},
	})

	if err != nil {
		log.Fatalf("加载规则失败: %v", err)
	}
	// 3. 启动 100 个协程，疯狂发送请求（模拟高并发流量）
	//我会在每一秒统计一次，这一秒只能 你通过了多少，总共有多少， block了多少, 每一秒会产生很多的block
	for i := 0; i < 100; i++ {
		go func() {
			for {
				// 无限循环发请求
				globalTotal++
				// Sentinel 核心：检查是否被限流
				e, b := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))
				// 被限流/拒绝了
				if b != nil {
					//fmt.Println("限流了")
					blockTotal++
					// 随机休眠 0~10ms，继续尝试
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
				} else {
					// 限流检查通过，可以执行业务
					passTotal++
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
					e.Exit()
				}
			}
		}()
	}

	// 4. 单独协程：每秒打印一次统计数据（看QPS变化）
	go func() {
		var oldTotal int //过去1s总共有多少个
		var oldPass int  //过去1s总共pass多少个
		var oldBlock int //过去1s总共block多少个
		for {
			oneSecondTotal := globalTotal - oldTotal
			oldTotal = globalTotal

			oneSecondPass := passTotal - oldPass
			oldPass = passTotal

			oneSecondBlock := blockTotal - oldBlock
			oldBlock = blockTotal

			time.Sleep(time.Second)
			fmt.Printf("total:%d, pass:%d, block:%d\n", oneSecondTotal, oneSecondPass, oneSecondBlock)
		}
	}()
	// 阻塞主goroutine，让程序一直运行
	<-ch
}
