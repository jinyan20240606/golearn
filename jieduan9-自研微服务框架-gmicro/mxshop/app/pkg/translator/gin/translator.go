package gin

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

// 接收校验失败的错误 → 翻译成中文 → 友好返回给前端
// err error：校验失败的错误，trans：翻译器（转中文用）
func HandleValidatorError(c *gin.Context, err error, trans ut.Translator) {
	errs, ok := err.(validator.ValidationErrors)
	// 如果不是校验错误，直接返回原始错误
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		// 去掉结构体名字，让错误更干净
		// 英文错误 翻译成中文
		"error": removeTopStruct(errs.Translate(trans)),
	})
	return
}
