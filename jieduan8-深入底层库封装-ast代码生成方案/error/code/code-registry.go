package code

import (
	"net/http"

	"mxshop/pkg/errors"
)

// 错误码规则：6位 AABBCC
// AA: 服务编号
// BB: 模块编号
// CC: 具体错误

// coder 实现 errors.Coder 接口，用于向 pkg/errors 注册错误码
type coder struct {
	code int    // 业务错误码（6位）
	http int    // 对应 HTTP 状态码
	ext  string // 面向用户的错误描述
	ref  string // 参考文档链接
}

func (c coder) Code() int         { return c.code }
func (c coder) HTTPStatus() int   { return c.http }
func (c coder) String() string    { return c.ext }
func (c coder) Reference() string { return c.ref }

// register 是注册错误码的快捷方法
func register(code, httpStatus int, message, ref string) {
	errors.MustRegister(coder{
		code: code,
		http: httpStatus,
		ext:  message,
		ref:  ref,
	})
}

// init 在包初始化时将所有错误码注册到 pkg/errors 的全局注册表
// 使用方只需 _ "mxshop/app/pkg/code" 空导入即可完成注册
func init() {
	// ==========================
	// 10 00 xx 基础错误
	// ==========================
	register(ErrSuccess, http.StatusOK, "OK", "")
	register(ErrUnknown, http.StatusInternalServerError, "未知错误", "")
	register(ErrParamInvalid, http.StatusBadRequest, "参数无效", "")
	register(ErrValidation, http.StatusBadRequest, "数据校验失败", "")
	register(ErrBind, http.StatusBadRequest, "请求体解析失败", "")
	register(ErrTooManyRequests, http.StatusTooManyRequests, "请求过于频繁，请稍后重试", "")

	// ==========================
	// 10 01 xx 数据库错误
	// ==========================
	register(ErrDatabaseError, http.StatusInternalServerError, "数据库错误", "")
	register(ErrDatabaseRecordNotFound, http.StatusNotFound, "记录不存在", "")
	register(ErrDatabaseDuplicate, http.StatusConflict, "数据重复", "")
	register(ErrDatabaseExec, http.StatusInternalServerError, "SQL执行失败", "")
	register(ErrDatabaseTransaction, http.StatusInternalServerError, "事务失败", "")

	// ==========================
	// 10 02 xx 认证/Token 错误
	// ==========================
	register(ErrTokenInvalid, http.StatusUnauthorized, "Token无效", "")
	register(ErrTokenTimeout, http.StatusUnauthorized, "Token已过期", "")
	register(ErrTokenGenerate, http.StatusInternalServerError, "Token生成失败", "")
	register(ErrPermissionDenied, http.StatusForbidden, "权限不足", "")
	register(ErrPasswordInvalid, http.StatusBadRequest, "密码错误", "")

	// ==========================
	// 10 03 xx 编解码/加解密
	// ==========================
	register(ErrEncrypt, http.StatusInternalServerError, "加密失败", "")
	register(ErrDecrypt, http.StatusInternalServerError, "解密失败", "")
	register(ErrSignInvalid, http.StatusBadRequest, "签名无效", "")
	register(ErrEncode, http.StatusInternalServerError, "编码失败", "")
	register(ErrDecode, http.StatusInternalServerError, "解码失败", "")

	// ==========================
	// 11 01 xx 商品模块错误 (goods.go 中定义的常量)
	// ==========================
	register(ErrGoodsNotFound, http.StatusNotFound, "商品不存在", "")
	register(ErrGoodsSoldOut, http.StatusBadRequest, "商品已下架", "")
	register(ErrGoodsCreateFail, http.StatusInternalServerError, "商品创建失败", "")
	register(ErrGoodsStockNotEnough, http.StatusBadRequest, "商品库存不足", "")

	// ==========================
	// 11 02 xx 用户模块错误（示例：RPC server 中使用的 ErrUserNotFound）
	// ==========================
	register(ErrUserNotFound, http.StatusNotFound, "用户不存在", "")
	register(ErrUserAlreadyExist, http.StatusConflict, "用户已存在", "")
	register(ErrUserLoginFail, http.StatusUnauthorized, "登录失败", "")
}
