package user

import (
	v1 "mxshop/api/user/v1"
	srv1 "mxshop/app/user/srv/service/v1"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewUserServer)

// 内部使用的server结构体
type userServer struct {
	v1.UnimplementedUserServer
	srv srv1.UserSrv
}

//func (us *userServer) mustEmbedUnimplementedUserServer() {
//	//TODO implement me
//	panic("implement me")
//}

// java中的ioc，控制翻转 ioc = injection of control
// 内部函数级别调用时没必要使用ioc（因为会带来一定的复杂度），一般在代码分层，第三方服务， rpc， redis等可以用ioc
// 对其他模块只暴露实例化方法，传入具体srv1.UserSrv实现，返回v1.UserServer接口实例，一般在初始化模块中调用，在controller其他文件中使用
func NewUserServer(srv srv1.UserSrv) v1.UserServer {
	return &userServer{srv: srv}
}

// 编译期断言也明确说明了这一点：var _ v1.UserServer = &userServer{}
// 这句表示 userServer 必须完整实现 v1.UserServer，否则编译不过。
var _ v1.UserServer = &userServer{}
