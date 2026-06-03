// 包名：log 自定义日志封装包
package log

// 这个文件作用：定义日志配置项 Options 以及支持相关的命令行参数解析到options配置项中、支持校验、构建zap实例等功能

import (
	"fmt"
	"strings"

	"encoding/json"

	"github.com/spf13/pflag"  // 命令行参数解析：增强版的标准库 flag
	"go.uber.org/zap"         // 底层zap高性能日志库
	"go.uber.org/zap/zapcore" // zap核心配置
)

// 常量定义：定义一些命令行flag的参数名，将来在 AddFlags 中映射
const (
	flagLevel             = "log.level"              // 日志级别
	flagDisableCaller     = "log.disable-caller"     // 关闭调用者信息
	flagDisableStacktrace = "log.disable-stacktrace" // 关闭堆栈追踪
	flagFormat            = "log.format"             // 日志格式 console/json
	flagEnableColor       = "log.enable-color"       // 开启颜色
	flagOutputPaths       = "log.output-paths"       // 日志输出路径
	flagErrorOutputPaths  = "log.error-output-paths" // 错误日志输出路径
	flagDevelopment       = "log.development"        // 是否开发模式
	flagName              = "log.name"               // 日志器名称

	consoleFormat = "console" // 控制台友好格式
	jsonFormat    = "json"    // JSON格式，生产环境用
)

// Options 日志所有配置项（结构体映射配置文件+命令行）
// 作用：统一管理所有日志参数
type Options struct {
	OutputPaths       []string `json:"output-paths"       mapstructure:"output-paths"`       // 输出位置：stdout/文件
	ErrorOutputPaths  []string `json:"error-output-paths" mapstructure:"error-output-paths"` // 错误输出
	Level             string   `json:"level"              mapstructure:"level"`              // 日志级别 debug/info/warn/error
	Format            string   `json:"format"             mapstructure:"format"`             // 格式 console/json
	DisableCaller     bool     `json:"disable-caller"     mapstructure:"disable-caller"`     // 关闭打印文件名+行号
	DisableStacktrace bool     `json:"disable-stacktrace" mapstructure:"disable-stacktrace"` // 关闭错误堆栈
	EnableColor       bool     `json:"enable-color"       mapstructure:"enable-color"`       // 控制台彩色输出
	Development       bool     `json:"development"        mapstructure:"development"`        // 开发模式
	Name              string   `json:"name"               mapstructure:"name"`               // 日志器名称
	EnableTraceID     bool     `json:"enable-trace-id"     mapstructure:"enable-trace-id"`   // 开启trace_id（链路追踪）
	EnableTraceStack  bool     `json:"enable-trace-stack" mapstructure:"enable-trace-stack"` // 开启追踪堆栈
}

// 入口方法：NewOptions 创建一个带**默认值**的日志配置
// 作用：保证开箱即用，不用手动填一堆默认值
func NewOptions() *Options {
	return &Options{
		Level:             zapcore.InfoLevel.String(), // 默认级别：Info
		DisableCaller:     false,                      // 默认显示调用行号
		DisableStacktrace: false,                      // 默认显示堆栈
		Format:            consoleFormat,              // 默认控制台格式
		EnableColor:       false,                      // 默认关闭颜色
		Development:       false,                      // 默认非开发模式
		OutputPaths:       []string{"stdout"},         // 默认输出到控制台（云原生标准）
		ErrorOutputPaths:  []string{"stderr"},         // 错误输出到stderr
	}
}

// Validate 校验配置是否合法
// 作用：防止用户乱填日志级别、格式，导致程序启动失败
func (o *Options) Validate() []error {
	var errs []error

	// 校验日志级别是否合法
	var zapLevel zapcore.Level
	// UnmarshalText 就是：把一段文本（字符串），解析成结构体 / 枚举变量的方法。
	// 和 json.Unmarshal 不同，UnmarshalText 是针对单个字段的解析，而不是整个结构体。
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		errs = append(errs, err)
	}

	// 校验日志格式只支持 console/json
	format := strings.ToLower(o.Format)
	if format != consoleFormat && format != jsonFormat {
		errs = append(errs, fmt.Errorf("not a valid log format: %q", o.Format))
	}

	return errs
}

// AddFlags ：从命令行启动参数中读取配置，赋值到options结构体中
// 作用：可以通过 命令行中--log.level=debug 动态修改日志配置
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Level, flagLevel, o.Level, "Minimum log output `LEVEL`.")
	// fs.StringVar(
	// 	&o.Level,          // 1. 要把值存到哪里（指针）
	// 	flagLevel,         // 2. 命令行参数叫什么名字：--log.level
	// 	o.Level,           // 3. 默认值是什么（不传参就用这个）
	// 	"说明文字"         // 4. --help 时显示的提示
	// 	)
	fs.BoolVar(&o.DisableCaller, flagDisableCaller, o.DisableCaller, "Disable output of caller information in the log.")
	fs.BoolVar(&o.DisableStacktrace, flagDisableStacktrace,
		o.DisableStacktrace, "Disable the log to record a stack trace for all messages at or above panic level.")
	fs.StringVar(&o.Format, flagFormat, o.Format, "Log output `FORMAT`, support plain or json format.")
	fs.BoolVar(&o.EnableColor, flagEnableColor, o.EnableColor, "Enable output ansi colors in plain format logs.")
	fs.StringSliceVar(&o.OutputPaths, flagOutputPaths, o.OutputPaths, "Output paths of log.")
	fs.StringSliceVar(&o.ErrorOutputPaths, flagErrorOutputPaths, o.ErrorOutputPaths, "Error output paths of log.")
	fs.BoolVar(
		&o.Development,
		flagDevelopment,
		o.Development,
		"Development puts the logger in development mode, which changes "+
			"the behavior of DPanicLevel and takes stacktraces more liberally.",
	)
	fs.StringVar(&o.Name, flagName, o.Name, "The name of the logger.")
}

// String 把配置序列化为JSON，方便打印查看
func (o *Options) String() string {
	data, _ := json.Marshal(o)
	return string(data)
}

// Build 最核心函数：
// 根据Options构建一个zap日志实例，并设置为全局logger
func (o *Options) Build() error {
	var zapLevel zapcore.Level
	// 解析日志级别，非法则默认info
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	// 日志级别显示：是否带颜色
	encodeLevel := zapcore.CapitalLevelEncoder
	if o.Format == consoleFormat && o.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// 组装zap原生配置
	zc := &zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel), // 日志级别
		Development:       o.Development,                  // 开发模式
		DisableCaller:     o.DisableCaller,                // 是否关闭调用者
		DisableStacktrace: o.DisableStacktrace,            // 是否关闭堆栈
		Sampling: &zap.SamplingConfig{ // 日志采样，防止日志刷屏
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: o.Format, // 输出格式 console / json
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",    // 日志内容key
			LevelKey:       "level",      // 级别key
			TimeKey:        "timestamp",  // 时间key
			NameKey:        "logger",     // 日志器名key
			CallerKey:      "caller",     // 调用文件+行号key
			StacktraceKey:  "stacktrace", // 堆栈key
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    encodeLevel,                 // 级别显示格式（彩色/大写）
			EncodeTime:     timeEncoder,                 // 时间格式化（外部定义）
			EncodeDuration: milliSecondsDurationEncoder, // 耗时格式化
			EncodeCaller:   zapcore.ShortCallerEncoder,  // 短文件名/行号
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths:      o.OutputPaths,      // 输出位置
		ErrorOutputPaths: o.ErrorOutputPaths, // 错误输出位置
	}

	// 真正创建zap logger
	logger, err := zc.Build(zap.AddStacktrace(zapcore.PanicLevel))
	if err != nil {
		return err
	}

	// 替换zap全局logger
	zap.RedirectStdLog(logger.Named(o.Name))
	zap.ReplaceGlobals(logger)

	return nil
}

// ===================== 下面是日志扩展选项（适配opentelemetry/trace） =====================

// Option 函数式选项模式，灵活配置Logger
type Option func(l *Logger)

// WithMinLevel 设置日志最小级别
func WithMinLevel(lvl zapcore.Level) Option {
	return func(l *Logger) {
		l.minLevel = lvl
	}
}

// WithErrorStatusLevel 当日志 >= 此级别时，自动把span标记为错误
// 作用：集成opentelemetry链路追踪
func WithErrorStatusLevel(lvl zapcore.Level) Option {
	return func(l *Logger) {
		l.errorStatusLevel = lvl
	}
}

// WithCaller 开启/关闭调用行号打印
func WithCaller(on bool) Option {
	return func(l *Logger) {
		l.caller = on
	}
}

// WithStackTrace 开启/关闭堆栈追踪
func WithStackTrace(on bool) Option {
	return func(l *Logger) {
		l.stackTrace = on
	}
}

// WithTraceIDField 日志中自动输出 trace_id
// 作用：全链路追踪必备
func WithTraceIDField(on bool) Option {
	return func(l *Logger) {
		l.withTraceID = on
	}
}
