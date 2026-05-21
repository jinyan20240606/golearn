package pay

import (
	"context"
	"mxshop-api/order-web/proto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"

	"mxshop-api/order-web/global"
)

// Notify
// 支付宝异步回调接口（支付成功后，支付宝服务器主动请求这个接口）
// 这个接口必须是公网能访问到的地址，不能是本地localhost
// 作用：接收支付结果，更新订单状态 = 支付闭环的最后一步
func Notify(ctx *gin.Context) {
	// ==================== 1. 初始化支付宝客户端 ====================
	// 支付宝回调通知
	// 生产环境和测试环境的网关是不一样的
	client, err := alipay.New(global.ServerConfig.AliPayInfo.AppID, global.ServerConfig.AliPayInfo.PrivateKey, false)
	if err != nil {
		zap.S().Errorw("实例化支付宝失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// ==================== 2. 加载【支付宝公钥】（安全核心） ====================
	// 作用：用支付宝公钥 验证回调签名
	// 确保这个请求真的是支付宝官方发的，不是黑客伪造的
	err = client.LoadAliPayPublicKey((global.ServerConfig.AliPayInfo.AliPublicKey))
	if err != nil {
		zap.S().Errorw("加载支付宝的公钥失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// ==================== 3. 验签（最最核心的安全步骤） ====================
	// 内部自动做了：
	// 1. 解析支付宝传过来的所有参数
	// 2. 用【支付宝公钥】验证签名是否正确
	// 3. 验签失败会直接返回error → 代表是伪造请求
	noti, err := client.GetTradeNotification(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	// ==================== 4. 验签通过 → 更新订单状态 ====================
	// 只有走到这里，才敢相信：支付真的成功了
	// OutTradeNo = 订单号 OrderSn
	// TradeStatus = 支付状态（TRADE_SUCCESS）
	_, err = global.OrderSrvClient.UpdateOrderStatus(context.Background(), &proto.OrderStatus{
		OrderSn: noti.OutTradeNo,
		Status:  string(noti.TradeStatus),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	// ==================== 5. 返回 success 给支付宝 ====================
	// 必须返回 "success" 字符串，支付宝才会停止重复回调
	// 否则支付宝会一直重复请求这个接口
	ctx.String(http.StatusOK, "success")
}
