package main

import (
	proto "GoStart/api/inventory/v1"
	"fmt"
	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v3"
)

func main() {
	r := gin.Default()
	r.GET("start", func(c *gin.Context) {
		orderSn := shortuuid.New()
		req := &proto.SellInfo{
			GoodsInfo: []*proto.GoodsInvInfo{
				{
					GoodsId: 421,
					Num:     2,
				},
			},
			OrderSn: orderSn,
		}
		dmtServer := "127.0.0.1:36790"
		qsBusi := "discovery:///mxshop-inventory-srv"
		fmt.Println(orderSn)
		saga := dtmgrpc.NewSagaGrpc(dmtServer, orderSn).
			Add(qsBusi+"/Inventory/Sell", qsBusi+"/Inventory/Reback", req)
		err := saga.Submit()
		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
		}
		c.JSON(200, gin.H{"message": "ok"})
	})

	r.Run(":8089")
}
