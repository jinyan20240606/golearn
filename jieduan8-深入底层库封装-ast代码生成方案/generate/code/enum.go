//go:generate stringer -type ErrCode -linecomment

// 上面就是意思就是调用stringer工具生成代码，所以我们需要单独安装这个可执行文件在go系统路径下，以便她能找到
// -type ErrCode：是 stringer工具的 唯一必选参数，意思是到底要给哪个类型生成 String () 方法。你写的枚举类型名字叫 ErrCode，工具只会给 ErrCode 生成代码
// // ---- stringer 只认：type 名字 = 整数类型（int/uint 系列），其他不支持
// -linecomment（最关键！）：让 stringer 读取【行尾注释】当作字符串返回，不加这个参数：只会输出 ErrCode(110001)显示数字，加了：直接返回注释文案
// 当你打印 / 输出 / 格式化这个枚举值ErrCode时，linecomment生成代码时会自动加个ErrCode的自定义String方法，拦截打印时的输出逻辑，输出对应的中文注释
// Go 自动调用 String () 方法
// 返回你写的注释：ok或参数错误或超时
// go generate执行后：自动生成 errcode_string.go，以后打印错误码，直接输出中文注释，超级方便
package code

type ErrCode int64 //错误码

// 错误码
const (
	ERR_CODE_OK             ErrCode = 0 //ok
	ERR_CODE_INVALID_PARAMS ErrCode = 1 //参数错误
	ERR_CODE_TIMEOUT        ErrCode = 2 //超时
)
