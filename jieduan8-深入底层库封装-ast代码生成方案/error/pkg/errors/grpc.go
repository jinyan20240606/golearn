package errors

import (
	"encoding/json"
	"errors"
	"fmt"

	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// grpcStatus 用于 JSON 序列化/反序列化 withCode 信息到 gRPC message
type grpcStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ToGrpcError 将带有 withCode 的 error 转换为 gRPC error。
// ------ 一般用在grpc服务端中返回自定义错误消息时使用
//
// ⚠️ 为什么不能直接 status.Error(codes.Code(perr.code), perr.err.Error()) ？
//   - perr.code 是我们自定义的6位业务错误码，如 110200
//   - gRPC codes.Code 只有 0~16 共 17 个合法枚举值（OK=0, NotFound=5, Internal=13...）
//   - 直接把 110200 强转为 codes.Code，gRPC 会把它当作 Unknown 处理，
//     客户端用 status.FromError 拿到的 Code() 不可预期
//   - 所以必须通过 httpStatusToGrpcCode 将 HTTP 状态码（如404）映射为标准 gRPC code（如 NotFound=5）
//   - 业务错误码则序列化到 message 字段中，通过 JSON 传输
//
// 原理：将业务错误码和错误消息序列化为 JSON，放入 gRPC status 的 message 字段，
// 同时根据注册的 HTTP 状态码映射到合适的 gRPC codes。
//
// 使用示例：
//
//	e := errors.WithCode(code.ErrUserNotFound, "user not found")
//	return nil, errors.ToGrpcError(e)
func ToGrpcError(err error) error {
	if err == nil {
		return nil
	}

	// 尝试解析为 withCode error
	var w *withCode
	// 必须用As深层断言是否含withCode类型，如果断言成功，直接赋值到w变量中，w就不是nil了
	if !errors.As(err, &w) {
		// 非 withCode error，直接作为 Unknown gRPC error
		return status.Error(grpccodes.Unknown, err.Error())
	}

	// 如果是自定义错误类型：就获取注册的 Coder 信息
	coder := ParseCoder(err)

	// 构造 grpcStatus，将业务 code 和 message 序列化到 JSON
	gs := grpcStatus{
		Code:    coder.Code(),
		Message: w.err.Error(),
	}

	data, _ := json.Marshal(gs)

	// 将 HTTP 状态码映射到 gRPC status code
	grpcCode := httpStatusToGrpcCode(coder.HTTPStatus())

	return status.Error(grpcCode, string(data))
}

// FromGrpcError 将 gRPC error 还原为带有 withCode 的 error。
// 原理：从 gRPC status 的 message 字段中反序列化出业务错误码，
// 然后重建 withCode error。
//
// 使用示例（客户端必须先空导入 code 包触发注册）：
//
//	import _ "mxshop/app/pkg/code"
//
//	s := errors.FromGrpcError(err)
//	coder := errors.ParseCoder(s)
//	fmt.Println(coder.Code(), coder.HTTPStatus(), coder.String())
func FromGrpcError(err error) error {
	if err == nil {
		return nil
	}

	// 解析 gRPC status
	st, ok := status.FromError(err)
	if !ok {
		// 非 gRPC error，直接返回原始 error
		return err
	}

	// 尝试从 message 中反序列化业务错误信息
	var gs grpcStatus
	if jsonErr := json.Unmarshal([]byte(st.Message()), &gs); jsonErr != nil {
		// message 不是 JSON 格式，是普通 gRPC error，用 unknownCoder 的 code
		return &withCode{
			err:   fmt.Errorf(st.Message()),
			code:  unknownCoder.Code(),
			stack: callers(),
		}
	}

	// 还原为 withCode error，code 字段就是业务错误码
	return &withCode{
		err:   fmt.Errorf(gs.Message),
		code:  gs.Code,
		stack: callers(),
	}
}

// httpStatusToGrpcCode 将 HTTP 状态码映射到 gRPC status code
func httpStatusToGrpcCode(httpStatus int) grpccodes.Code {
	switch httpStatus {
	case 200:
		return grpccodes.OK
	case 400:
		return grpccodes.InvalidArgument
	case 401:
		return grpccodes.Unauthenticated
	case 403:
		return grpccodes.PermissionDenied
	case 404:
		return grpccodes.NotFound
	case 409:
		return grpccodes.AlreadyExists
	case 429:
		return grpccodes.ResourceExhausted
	case 500:
		return grpccodes.Internal
	case 501:
		return grpccodes.Unimplemented
	case 503:
		return grpccodes.Unavailable
	case 504:
		return grpccodes.DeadlineExceeded
	default:
		return grpccodes.Unknown
	}
}
