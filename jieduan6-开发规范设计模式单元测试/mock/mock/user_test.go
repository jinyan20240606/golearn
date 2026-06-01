package mock

import (
	"GoStart/mock"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetUserByMobile(t *testing.T) {
	//mock 准备工作
	// ====================== 1. mock 核心流程 ======================
	// 1. 创建 mock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 测试结束自动验证预期
	// 2. 创建 mock 对象（由 mockgen 自动生成）
	mockUserData := NewMockUserData(ctrl)
	// 3. 设置 mock 预期（重点：打桩 Stub）
	// 当有人调用 GetUserByMobile 方法，参数是（任意上下文，手机号 18）时，必须返回 用户bobby18 和 无错误。
	mockUserData.EXPECT().GetUserByMobile(gomock.Any(), "18").Return(mock.User{ // 预期入参：任意ctx + 手机号18
		NickName: "bobby18",
	}, nil) // 预期返回：用户bobby18 + 无错误

	// ====================== 2. 业务调用 ======================
	// 注入 mock 到 UserServer（依赖注入）
	userServer := mock.UserServer{
		Db: mockUserData, // Db 必须是大写导出字段
	}
	// 执行业务方法
	user, err := userServer.GetUserByMobile(context.Background(), "18")

	/**
	大致逻辑：你要对业务方法进行单元测试，先通过mockUserData.EXPECT这个方法进行数据库接口的数据mock，然后你用这个mock的数据进行真正测试
	业务方法 → 调用 mock 的 GetUserByMobile
		→ 返回我们预设的 bobby18
		→ 业务逻辑判断：如果是 bobby18 就改成 bobby17
	*/

	// ====================== 3. 断言结果 ======================
	// 断言1：无错误
	if err != nil {
		t.Errorf("error: %v", err)
	}
	// 断言2：返回的昵称与预期一致（必须和上面Return一致！）
	if user.NickName != "bobby17" {
		t.Errorf("error: %v", err)
	}

	t.Log("测试通过！")

}

func TestGetUserByMobileFail(t *testing.T) {
	//mock 准备工作
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserData := NewMockUserData(ctrl)
	mockUserData.EXPECT().GetUserByMobile(gomock.Any(), "19").Return(mock.User{
		NickName: "bobby19",
	}, nil)

	//实际调用过程
	userServer := mock.UserServer{
		Db: mockUserData,
	}
	user, err := userServer.GetUserByMobile(context.Background(), "19")

	//判断正确与否
	if err != nil {
		t.Errorf("error: %v", err)
	}

	if user.NickName != "bobby19" {
		t.Errorf("error: %v", err)
	}

}
