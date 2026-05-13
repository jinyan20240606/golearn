package forms

// 表单相关的结构体单独维护在form目录

type PassWordLoginForm struct {
	// form:"mobile"：接收前端 form 表单提交的参数，Gin 会自动绑定到结构体的 Mobile 字段
	// json:"mobile"：作用：接收前端 JSON 格式提交的数据，Gin 自动绑定到结构体
	// binding:"xxx" （最重要！校验规则）参数自动校验，Gin 内置 validator 库，不满足规则直接报错，不用你手写 if 判断
	// // // // // // // // --- mobile：必须是合法手机号格式（你自定义的校验器）
	Mobile    string `form:"mobile" json:"mobile" binding:"required,mobile"`           //手机号码格式有规范可寻， 自定义validator
	PassWord  string `form:"password" json:"password" binding:"required,min=3,max=20"` // 多个条件时，中间不能有空格
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"`
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`
}

type RegisterForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"` //手机号码格式有规范可寻， 自定义validator
	PassWord string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Code     string `form:"code" json:"code" binding:"required,min=6,max=6"`
}

type UpdateUserForm struct {
	Name     string `form:"name" json:"name" binding:"required,min=3,max=10"`
	Gender   string `form:"gender" json:"gender" binding:"required,oneof=female male"`
	Birthday string `form:"birthday" json:"birthday" binding:"required,datetime=2006-01-02"`
}
