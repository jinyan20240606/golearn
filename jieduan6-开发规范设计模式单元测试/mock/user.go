package mock

import "context"

/*
gomock 测试
安装：
	go get -u github.com/golang/mock/gomock

	go install github.com/golang/mock/mockgen
*/

type User struct {
	Mobile   string
	Password string
	NickName string
}

type UserServer struct {
	Db UserData
}

// 我们主要单元测试是想测的这个业务逻辑，mockgen生成的主要是里面依赖的DB数据接口类型的具体实现
// 你要对GetUserByMobile业务方法进行单元测试，先通过mockUserData.EXPECT这个方法进行 us.Db.GetUserByMobile(ctx, mobile)响应数据Mock
// ，然后你用这个mock的数据user进行真正的后续业务逻辑测试
func (us *UserServer) GetUserByMobile(ctx context.Context, mobile string) (User, error) {
	// DB 操作：用 mockgen 生成的 Mock 模拟（不连真实数据库）
	user, err := us.Db.GetUserByMobile(ctx, mobile)
	if err != nil {
		return User{}, err
	}
	if user.NickName == "bobby18" {
		user.NickName = "bobby17"
	}
	return user, nil
}

// mockgen方法主要自动生成这个接口鸭子类型DB方法，能进行调用响应mock结果
type UserData interface {
	GetUserByMobile(ctx context.Context, mobile string) (User, error)
}
