package log

// 这个文件作用：定义日志配置项 Options 以及支持相关的命令行参数解析到options配置项中、支持校验等功能

// 负责：
// 定义日志有哪些配置（级别、格式、输出）
// 提供默认配置
// 提供校验配置是否合法
// 提供命令行参数绑定（--log.level）
// 它不创建日志！只存参数！
import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
)

// 常量定义：flag开头的是一些命令行flag的参数名，将来在 AddFlags 中映射
const (
	FORAMT_CONSOLE = "console" // 参数名常量
	FORAMT_JSON    = "json"
	OUTPUT_STD     = "stdout"
	OUTPUT_STD_ERR = "stderr"

	flagLevel = "log.level"
)

type Options struct {
	OutputPaths     []string `json:"output-paths" mapstructure:"output-paths"`
	ErrorOuputPaths []string `json:"error-output-paths" mapstructure:"error-output-paths"`
	Level           string   `json:"level" mapstructure:"level"`
	Format          string   `json:"format" mapstructure:"format"`
	Name            string   `json:"name" mapstructure:"name"`
	// 作用：给 JSON 序列化 / 反序列化用
	// mapstructure:"name" 作用：给 Viper 配置文件解析用（YAML/TOML 等）
}

// 支持函数选项模式来传参options
type Option func(o *Options)

func NewOptions(opts ...Option) *Options {
	options := &Options{
		Level:           zapcore.InfoLevel.String(),
		Format:          FORAMT_CONSOLE,
		OutputPaths:     []string{OUTPUT_STD},
		ErrorOuputPaths: []string{OUTPUT_STD_ERR},
	}

	for _, opt := range opts {
		opt(options)
	}
	return options
}

func WithLevel(level string) Option {
	return func(o *Options) {
		o.Level = level
	}
}

// 就可以自定去定义检测规则
func (o *Options) Validate() []error {
	var errs []error
	format := strings.ToLower(o.Format)
	if format != FORAMT_CONSOLE && format != FORAMT_JSON {
		errs = append(errs, fmt.Errorf("not suppor format %s", o.Format))
	}
	return errs
}

// 可以自己将options具体的列映射到flag的字段上
// AddFlags ：从命令行启动参数中读取配置，赋值到options结构体中
// 作用：可以通过 命令行中--log.level=debug 动态修改日志配置
func (o *Options) AddFlags(fs pflag.FlagSet) {
	// fs.StringVar(
	// 	&o.Level,          // 1. 要把值存到哪里（指针）
	// 	flagLevel,         // 2. 命令行参数叫什么名字：--log.level
	// 	o.Level,           // 3. 默认值是什么（不传参就用这个）
	// 	"说明文字"         // 4. --help 时显示的提示
	// 	)
	fs.StringVar(&o.Level, flagLevel, o.Level, "log level")
}
