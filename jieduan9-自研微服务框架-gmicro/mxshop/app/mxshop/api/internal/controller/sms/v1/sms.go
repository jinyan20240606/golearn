package v1

import (
	"mxshop/app/mxshop/api/internal/service"
	v1 "mxshop/app/mxshop/api/internal/service/sms/v1"
	"mxshop/app/pkg/code"
	gin2 "mxshop/app/pkg/translator/gin"
	"mxshop/pkg/common/core"
	"mxshop/pkg/errors"
	"mxshop/pkg/storage"
	"time"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
)

type SendSmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"` //手机号码格式有规范可寻， 自定义validator
	Type   uint   `form:"type" json:"type" binding:"required,oneof=1 2"`
	//1. 注册发送短信验证码和动态验证码登录发送验证码
}

type SmsController struct {
	sf    service.ServiceFactory
	trans ut.Translator
}

func NewSmsController(sf service.ServiceFactory, trans ut.Translator) *SmsController {
	return &SmsController{sf, trans}
}

// 这是 Gin 接口控制器方法，实现手机验证码发送完整流程：
// 接收前端手机号参数并校验
// 生成 6 位随机短信验证码
// 调用短信服务商接口下发短信
// 把 手机号-验证码 存入 Redis，设置 5 分钟过期
// 统一返回响应
func (sc *SmsController) SendSms(c *gin.Context) {
	sendSmsForm := SendSmsForm{}
	if err := c.ShouldBind(&sendSmsForm); err != nil {
		gin2.HandleValidatorError(c, err, sc.trans)
	}

	smsCode := v1.GenerateSmsCode(6)
	err := sc.sf.Sms().SendSms(c, sendSmsForm.Mobile, "SMS_181850725", "{\"code\":"+smsCode+"}")
	if err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrSmsSend, err.Error()), nil)
		return
	}

	//将验证码保存起来 - redis
	rstore := storage.RedisCluster{}
	// 	SetKey：
	// Key：手机号（以手机号作为唯一键）
	// Value：验证码
	// 过期时间：5 分钟，超时自动失效

	err = rstore.SetKey(c, sendSmsForm.Mobile, smsCode, 5*time.Minute)
	if err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrSmsSend, err.Error()), nil)
		return
	}

	core.WriteResponse(c, nil, nil)
}
