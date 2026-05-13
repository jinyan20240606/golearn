package forms

type SendSmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"` //手机号码格式有规范可寻， 自定义validator
	Type   uint   `form:"type" json:"type" binding:"required,oneof=1 2"`  // 因为可能有很多种类的短信，1代表注册，2代表验证吗登录
	//1. 如注册发送短信验证码和动态验证码登录发送验证码 这就是不同类的短信
}
