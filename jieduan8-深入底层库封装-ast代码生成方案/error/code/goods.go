package code

// 错误码规则：6位 AABBCC
// AA: 服务编号
// BB: 模块编号
// CC: 具体错误

// -----------------------------------------
// 服务 & 模块 基础定义
// -----------------------------------------
const (
	ServiceGoods = 11 * 10000 // 11 商品服务

	// 商品服务下的模块
	ModGoodsBase = 1 * 100 // 01 商品基础模块
)

// ==========================
// 11 01 xx 商品模块
// ==========================
const (
	// ErrGoodsNotFound - 404: 商品不存在.
	ErrGoodsNotFound = ServiceGoods + ModGoodsBase + iota

	// ErrGoodsSoldOut - 400: 商品已下架.
	ErrGoodsSoldOut

	// ErrGoodsCreateFail - 500: 商品创建失败.
	ErrGoodsCreateFail

	// ErrGoodsStockNotEnough - 400: 商品库存不足.
	ErrGoodsStockNotEnough
)

// ==========================
// 11 02 xx 用户模块
// ==========================
const (
	ModUser = 2 * 100 // 02 用户模块
)

const (
	// ErrUserNotFound - 404: 用户不存在.
	ErrUserNotFound = ServiceGoods + ModUser + iota

	// ErrUserAlreadyExist - 409: 用户已存在.
	ErrUserAlreadyExist

	// ErrUserLoginFail - 401: 登录失败.
	ErrUserLoginFail
)
