package service

import (
	"context"
	"errors"
	"testing"
	"你的项目路径/model" // 替换为实际model包路径

	"github.com/stretchr/testify/assert"
)

// ====================== 1. 定义接口（契约层）======================
// UserData 数据操作接口
// 作用：定义数据层规范，业务层只依赖接口，不依赖具体实现
type UserData interface {
	GetUserByMobile(ctx context.Context, mobile string) (model.User, error)
}

// ====================== 2. 业务服务层（依赖接口）======================
// UserServer 业务逻辑服务
// 依赖：UserData接口（外部注入，不自己创建）
type UserServer struct {
	db UserData // 依赖注入核心字段
}

// NewUserServer 构造函数：统一注入依赖（标准写法）
// 依赖注入（DI）
// 控制反转：对象的依赖由外部创建传入，而非内部生成
// 优点：随时替换实现，无需修改业务代码

func NewUserServer(db UserData) *UserServer {
	return &UserServer{db: db}
}

// GetUserByMobile 业务方法
func (s *UserServer) GetUserByMobile(ctx context.Context, mobile string) (model.User, error) {
	// 业务逻辑：调用数据层接口获取数据
	return s.db.GetUserByMobile(ctx, mobile)
}

// ====================== 3. 真实数据实现（如MySQL）======================
// MySQLUserData 真实数据库实现
// 实现UserData接口，对接真实MySQL
type MySQLUserData struct {
	// 可存放数据库连接等字段
	// db *gorm.DB
}

// GetUserByMobile 实现UserData接口
func (m *MySQLUserData) GetUserByMobile(ctx context.Context, mobile string) (model.User, error) {
	// 真实SQL查询逻辑
	var user model.User
	// err := m.db.WithContext(ctx).Where("mobile = ?", mobile).First(&user).Error
	// return user, err
	return user, nil
}

// ====================== 4. Mock实现（单元测试专用）======================
// MockUserData 空Mock：直接返回空值，快速测试流程
type MockUserData struct{}

func (mud *MockUserData) GetUserByMobile(ctx context.Context, mobile string) (model.User, error) {
	return model.User{}, nil
}

// FakeUserData 带内存数据的Mock：自定义测试数据
type FakeUserData struct {
	data map[string]model.User // 内存存储测试数据
}

// GetUserByMobile 实现UserData接口
func (fud *FakeUserData) GetUserByMobile(ctx context.Context, mobile string) (model.User, error) {
	user, ok := fud.data[mobile]
	if !ok {
		return model.User{}, errors.New("用户不存在")
	}
	return user, nil
}

// 测试：使用FakeUserData模拟数据
func TestUserServer_GetUserByMobile(t *testing.T) {
	// 1. 构造测试数据
	testData := map[string]model.User{
		"13800138000": {ID: 1, Mobile: "13800138000", Name: "测试用户"},
	}

	// 2. 注入FakeUserData----依赖注入
	server := NewUserServer(&FakeUserData{data: testData})

	// 3. 执行测试
	user, err := server.GetUserByMobile(context.Background(), "13800138000")

	// 4. 断言结果
	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "测试用户", user.Name)
}
