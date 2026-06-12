package main

import (
	// 库存服务 proto 生成的结构体
	proto "GoStart/api/inventory/v1"
	"fmt"

	// DTM gRPC 客户端
	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/gin-gonic/gin"

	// 生成唯一订单号
	"github.com/lithammer/shortuuid/v3"
)

// 通过 DTM gRPC 客户端 创建 SAGA 事务：
// 正向操作：调用库存服务 Inventory/Sell（扣减库存）
// 补偿操作：调用库存服务 Inventory/Reback（归还 / 回补库存）
// 提交 SAGA 事务给 DTM Server，由 DTM 负责全局调度、失败补偿、重试
func main() {
	r := gin.Default()
	r.GET("start", func(c *gin.Context) {
		// 1. 生成全局唯一订单号，同时作为 SAGA 的 gid(全局事务ID)
		orderSn := shortuuid.New()
		// 2. 组装库存扣减请求：商品ID=421，数量=2
		req := &proto.SellInfo{
			GoodsInfo: []*proto.GoodsInvInfo{
				{
					GoodsId: 421,
					Num:     2,
				},
			},
			OrderSn: orderSn,
		}
		// 3. DTM Server 地址（独立部署的DTM服务监听地址）
		dmtServer := "127.0.0.1:36790"
		// 4. 库存gRPC服务地址：使用服务发现(discovery)，对应注册中心
		qsBusi := "discovery:///mxshop-inventory-srv"
		// qsBusi := "127.0.0.1:8019" // 默认直连的话直接写库存微服务的端口号连grpc，上行是链接consul的服务发现的方式
		fmt.Println(orderSn)
		// 5. 创建 SAGA gRPC 实例
		// 参数：dtm地址 + 全局事务ID(gid=orderSn)
		saga := dtmgrpc.NewSagaGrpc(dmtServer, orderSn).
			// 添加一个SAGA步骤：正向接口、补偿接口、请求参数
			// 正向：扣库存，// 补偿：还库存，// 传给正向接口的参数
			// 3参为对应的请求参数
			Add(qsBusi+"/Inventory/Sell", qsBusi+"/Inventory/Reback", req)

		// 6. 提交SAGA事务，交给DTM调度执行
		err := saga.Submit()
		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
		}
		c.JSON(200, gin.H{"message": "ok"})
	})

	r.Run(":8089")
}
