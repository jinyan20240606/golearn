package utils

import (
	"net"
)

// 微服务启动时，自动获取当前操作系统空闲的端口号（随机可用端口），避免端口冲突！
//
// 返回可用端口号和错误
func GetFreePort() (int, error) {
	// 端口 0 是特殊端口
	// 意思：让操作系统自动分配一个空闲端口
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	// 尝试在 localhost:0 上监听
	// 操作系统会自动找一个当前没被占用的端口
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	// 函数结束后关闭监听
	// 不占用端口，只是临时拿一下端口号
	defer l.Close()
	// 拿到系统分配的真实端口，返回给你
	// l.Addr()返回的装地址类型的接口，需要对接口使用断言拿到具体的地址类型值，即不知道是 TCP 还是 UDP，把 Addr 接口 转换成 具体的 TCP 地址结构
	return l.Addr().(*net.TCPAddr).Port, nil
}
