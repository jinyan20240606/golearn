package validation

// 一句话总结
// 这段代码作用：
// 给 Gin 框架增加一个 mobile 手机号校验规则 + 中文错误提示
// 使用方式：在结构体 tag 里写 binding:"mobile"
import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"

	// validator.v10：Go 最流行的校验库
	"github.com/go-playground/validator/v10"
)

// 使用方式gin的标签校验机制：type User struct {    Mobile string `json:"mobile" binding:"required,mobile"`}
// 然后 Gin 参数绑定会自动校验：
func RegisterMobile(translator ut.Translator) {
	// binding.Validator：Gin 内置的校验器
	// validator.v10：Go 最流行的校验库
	// 把 Gin 的校验器转换成 *validator.Validate
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 1. 注册校验规则：tag 名叫 mobile
		_ = v.RegisterValidation("mobile", ValidateMobile)
		// 2. 注册中文翻译信息
		_ = v.RegisterTranslation("mobile", translator, func(ut ut.Translator) error {
			// 注册错误提示语
			// ut：翻译器（用来输出中文错误信息）
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
			// 翻译函数：把错误翻译成中文
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}
}

func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	//使用正则表达式判断是否合法
	ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
	if !ok {
		return false
	}
	return true
}
