package code

// 错误码规则：6位 AABBCC
// AA: 服务编号   (10=通用服务)
// BB: 模块编号
// CC: 具体错误

// -----------------------------------------
// 基础常量：服务编号 + 模块编号（用来组合错误码）
// -----------------------------------------
const (
	// AA 服务段
	ServiceCommon = 10 * 10000 // 10 通用服务 → 100000 基准

	// BB 模块段 (在 ServiceCommon 下)
	ModBase     = 0 * 100 // 00 基础模块
	ModDatabase = 1 * 100 // 01 数据库模块
	ModAuth     = 2 * 100 // 02 认证模块
	ModCrypto   = 3 * 100 // 03 编解码/加解密模块
)

// ==========================
// 10 00 xx 基础错误
// ==========================
const (
	// ErrSuccess - 200: OK.
	ErrSuccess = ServiceCommon + ModBase + iota // const组提供的iota默认从0开始，自增

	// ErrUnknown - 500: 未知错误.
	ErrUnknown

	// ErrParamInvalid - 400: 参数无效.
	ErrParamInvalid

	// ErrValidation - 400: 数据校验失败.
	ErrValidation

	// ErrBind - 400: 请求体解析失败.
	ErrBind

	// ErrTooManyRequests - 429: 请求频繁.
	ErrTooManyRequests
)

// ==========================
// 10 01 xx 数据库错误
// ==========================
const (
	// ErrDatabaseError - 500: 数据库错误.
	ErrDatabaseError = ServiceCommon + ModDatabase + iota

	// ErrDatabaseRecordNotFound - 404: 记录不存在.
	ErrDatabaseRecordNotFound

	// ErrDatabaseDuplicate - 409: 数据重复.
	ErrDatabaseDuplicate

	// ErrDatabaseExec - 500: SQL执行失败.
	ErrDatabaseExec

	// ErrDatabaseTransaction - 500: 事务失败.
	ErrDatabaseTransaction
)

// ==========================
// 10 02 xx 认证/Token 错误
// ==========================
const (
	// ErrTokenInvalid - 401: Token无效.
	ErrTokenInvalid = ServiceCommon + ModAuth + iota

	// ErrTokenTimeout - 401: Token过期.
	ErrTokenTimeout

	// ErrTokenGenerate - 500: Token生成失败.
	ErrTokenGenerate

	// ErrPermissionDenied - 403: 权限不足.
	ErrPermissionDenied

	// ErrPasswordInvalid - 400: 密码错误.
	ErrPasswordInvalid
)

// ==========================
// 10 03 xx 编解码/加解密
// ==========================
const (
	// ErrEncrypt - 500: 加密失败.
	ErrEncrypt = ServiceCommon + ModCrypto + iota

	// ErrDecrypt - 500: 解密失败.
	ErrDecrypt

	// ErrSignInvalid - 400: 签名无效.
	ErrSignInvalid

	// ErrEncode - 500: 编码失败.
	ErrEncode

	// ErrDecode - 500: 解码失败.
	ErrDecode
)
