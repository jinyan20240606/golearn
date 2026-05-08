package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	// validator的翻译器
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator

// 登录的时候需要验证的结构体
type LoginForm struct { // 作为验证器的结构体// binding:"required" 代表该字段是必填的
	User     string `json:"user" binding:"required,min=3,max=10"` // 需要接收的json类型的参数，对于form类型参数就会默认校验失败
	Password string `json:"password" binding:"required"`
}

// 注册的时候需要验证的结构体
type SignUpForm struct {
	Age      uint8  `json:"age" binding:"gte=1,lte=130"`
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	// eqfield=Password re_password 必须等于 Password 字段的值
	RePassword string `json:"re_password" binding:"required,eqfield=Password"` //支持跨字段验证
}

// 去掉错误提示中的结构体名称
func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

// 初始化 validator 多语言翻译器，让校验错误变成中文提示！
func InitTrans(locale string) (err error) {
	// 使用gin的bingding方法：获取 Gin 内部的 validator 实例，准备做自定义修改
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个方法：让错误提示里的字段名，使用你 json 标签里的名字(因为它默认用的是结构体里的字段名)
		// reflect.StructField为反射类型
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			// Go JSON 库官方支持逗号，二参为额外配置选项如json:"字段名,omitempty"：omitempty：字段为空时，不返回给前端，-：忽略这个字段，不序列化、不返回
			if name == "-" {
				return ""
			}
			return name
		})
		// 创建语言包（中文 + 英文）
		zhT := zh.New() //中文翻译器
		enT := en.New() //英文翻译器
		// 创建多语言管理工具：第一个参数是备用的语言环境，后面的参数是应该支持的语言环境
		uni := ut.New(enT, zhT, enT)
		// 1-获取当前语言的翻译器
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s)", locale)
		}
		// 2-注册默认翻译器（把 validator 的错误信息 → 注册成中文 / 英文提示
		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(v, trans)
		case "zh":
			zh_translations.RegisterDefaultTranslations(v, trans)
		default:
			en_translations.RegisterDefaultTranslations(v, trans)
		}
		return
	}

	return
}

func main() {
	//代码侵入性很强 中间件
	if err := InitTrans("zh"); err != nil {
		fmt.Println("初始化翻译器错误")
		return
	}
	router := gin.Default()
	// 登录接口
	router.POST("/loginJSON", func(c *gin.Context) {

		var loginForm LoginForm
		if err := c.ShouldBind(&loginForm); err != nil {
			errs, ok := err.(validator.ValidationErrors)
			if !ok {
				c.JSON(http.StatusOK, gin.H{
					"msg": err.Error(), // 这个错误将是带堆栈的英文错误信息
				})
			}
			c.JSON(http.StatusBadRequest, gin.H{
				// 返回具体错误时，使用翻译功能
				"error": removeTopStruct(errs.Translate(trans)),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg": "登录成功",
		})
	})

	// 注册接口，验证规则更复杂一些
	router.POST("/signup", func(c *gin.Context) {
		var signUpFrom SignUpForm
		if err := c.ShouldBind(&signUpFrom); err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg": "注册成功",
		})
	})

	_ = router.Run(":8083")
}
