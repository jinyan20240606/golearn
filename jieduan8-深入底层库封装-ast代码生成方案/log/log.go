package log

// 基于zap实现自定义logger封装
// 负责：
// 定义 Logger 接口
// 用 zap 实现日志功能
// 接收 Options 配置 → 创建真正的 logger

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zapcore.Field

// 首先定义个自己的logger的接口，暴露给外部使用，方便后续替换底层日志库
type Logger interface {
	// 带 C：支持传入 context.Context（拿 trace_id）
	//带 f：支持格式化字符串
	//带 W：支持键值对结构化日志
	Debug(msg string)
	DebugC(context context.Context, msg string)
	Debugf(format string, args ...interface{})
	DebugfC(context context.Context, msg string)
	DebugW(msg string, keysAndValues ...interface{})
	DebugWC(context context.Context, msg string, keysAndValues ...interface{})
	// 你的业务代码 只依赖这个 Logger 接口
	//不直接依赖 zap.Logger、logrus、zerolog
	//以后想把 zap 换成别的日志库，只需要写一个新的结构体实现这个接口，业务代码完全不动
}

type zapLogger struct {
	zapLogger *zap.Logger
}

var _ Logger = &zapLogger{}

type otherLogger struct {
}

func (z *zapLogger) Debug(msg string) {
	z.zapLogger.Debug(msg)
}

func (z *zapLogger) DebugC(context context.Context, msg string) {
	//TODO implement me
	panic("implement me")
}

func (z *zapLogger) Debugf(format string, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (z *zapLogger) DebugfC(context context.Context, msg string) {
	//TODO implement me
	panic("implement me")
}

func (z *zapLogger) DebugW(msg string, keysAndValues ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (z *zapLogger) DebugWC(context context.Context, msg string, keysAndValues ...interface{}) {
	//TODO implement me
	panic("implement me")
}

var (
	defaultLogger = New(NewOptions())
	mu            sync.Mutex
)

func Debug(msg string) {
	defaultLogger.Debug(msg)
}

// New 根据配置选项 opts，创建并返回一个封装好的 zapLogger 实例
// opts：日志配置，如果传 nil 则使用默认配置
func New(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}

	//实例化zap
	// 2. 把字符串日志级别（如 "info"）转换成 zap 内部的枚举级别，字符串映射成内部变量
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}
	// 3. 构建 zap 原生配置（最小化版本，为了演示，只设置级别字段）
	// Level 是原子级别，支持运行时动态修改日志级别
	loggerConfig := zap.Config{
		Level: zap.NewAtomicLevelAt(zapLevel),
	}
	// 4. 真正创建 zap 底层 logger
	// zap.AddStacktrace：只有 Panic 级别才打印堆栈，减少日志体积
	l, err := loggerConfig.Build(zap.AddStacktrace(zapcore.PanicLevel))
	if err != nil {
		panic(err)
	}
	// 5. 封装成我们自己的 zapLogger 结构体并返回
	// l.Named(opts.Name)：给日志器起一个名字，会出现在日志中
	logger := &zapLogger{
		zapLogger: l.Named(opts.Name),
	}
	return logger
}

// 入口方法：实例化全局logger
func Init(opt *Options) {
	// 看起来没有问题， 并发问题：如果多个goroutine同时来调用,本地可以用锁
	// 想解决并发问题：不要使用sync.Once,使用锁

	// 1. 加互斥锁 mu.Lock()
	// 目的：解决并发安全
	// 如果多个 goroutine 同时调用 Init()
	// 不加锁会出现 数据竞争（data race）
	// 加锁保证同一时间只有一个 goroutine 能修改全局变量
	// 2. 不用 sync.Once
	// 原因：支持动态更新日志
	// sync.Once 只能执行 一次，执行完永远不能改
	// 加锁 mu → 可以反复调用 Init () 动态更新日志配置
	// 运行时改日志级别
	// 运行时切换输出
	// 运行时热更新

	mu.Lock()
	defer mu.Unlock()
	defaultLogger = New(opt)
}
