package message

import (
	"context"
	"mxshop-api/userop-web/api"
	"mxshop-api/userop-web/forms"
	"mxshop-api/userop-web/global"
	"mxshop-api/userop-web/models"
	"mxshop-api/userop-web/proto"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// List 留言列表接口
// 作用：获取用户留言列表（普通用户只能看自己的，管理员可以看全部）
func List(ctx *gin.Context) {
	// 1. 初始化 proto 请求结构体（用于调用 gRPC 服务）
	request := &proto.MessageRequest{}
	// 2. 从 Gin Context 中获取登录用户的 userId（JWT 中间件设置的）
	userId, _ := ctx.Get("userId")
	// 3. 从 Context 中获取 JWT 解析后的用户权限信息
	claims, _ := ctx.Get("claims")
	// 4. 类型断言，转换成自定义的 CustomClaims（里面有 AuthorityId 角色ID）
	model := claims.(*models.CustomClaims)
	// 5. 判断角色：AuthorityId == 1 表示是普通用户
	// 普通用户只能查看【自己的留言】，所以要把自己的 userId 传给 gRPC
	if model.AuthorityId == 1 {
		request.UserId = int32(userId.(uint))
	}
	// 如果是管理员，不设置 UserId → gRPC 会返回【所有留言】

	// 6. 调用 gRPC 服务端的 MessageList 方法，获取留言列表数据
	rsp, err := global.MessageClient.MessageList(context.Background(), request)
	if err != nil {
		zap.S().Errorw("获取留言失败")
		// 将 gRPC 错误转换成 HTTP 错误返回给前端
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	// 7. 构造返回给前端的 map（包含 total 和 data）
	reMap := map[string]interface{}{
		"total": rsp.Total,
	}
	// 8. 定义一个切片，存放最终要返回的留言数据
	result := make([]interface{}, 0)
	// 9. 遍历 gRPC 返回的留言数据，逐个组装成前端需要的格式
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["user_id"] = value.UserId
		reMap["type"] = value.MessageType
		reMap["subject"] = value.Subject
		reMap["message"] = value.Message
		reMap["file"] = value.File
		// 加入结果切片
		result = append(result, reMap)
	}
	// 10. 把组装好的列表数据放进返回 map 的 data 字段
	reMap["data"] = result

	ctx.JSON(http.StatusOK, reMap)
}

func New(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")

	messageForm := forms.MessageForm{}
	if err := ctx.ShouldBindJSON(&messageForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	rsp, err := global.MessageClient.CreateMessage(context.Background(), &proto.MessageRequest{
		UserId:      int32(userId.(uint)),
		MessageType: messageForm.MessageType,
		Subject:     messageForm.Subject,
		Message:     messageForm.Message,
		File:        messageForm.File,
	})

	if err != nil {
		zap.S().Errorw("添加留言失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}
