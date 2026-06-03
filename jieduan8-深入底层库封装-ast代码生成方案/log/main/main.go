package main

func main() {
	//log包很重要
	/*
		开发，debug， 故障排查，数据分析， 监控告警，保存现场
		我们需要设计一个优秀的日志包，如果我们要扩展就比较麻烦,常见3种方式：，1. 基于zap封装， 2. 自己实现 3. 改zap的源码
		基于zap封装需要考虑如下问题：
		1. 是否可以替换 后期我们想要替换成另一个日志框架
		2. 我们要考虑扩展性， log打印的时候是否支持打印当前的goroutine的id 是否支持打印当前的context
		3. 我们给大家提供的日志包， 还能支持集成tracing(open-telemetry, metrics, logging),就可以集成jaeger
		4. 是否每个日志打印都能知道这个日志是哪个请求的
		封装日志包很重要！最好是自己封装
	*/

	//gorm， go-redis、我们自己业务代码
	/*
			logger最基本的需求功能
				1. 日志基本的级别 debug、info、warn、error 、fatal、 panic
				2. 打印方式2种：2020-12-02T01:16:18+08:00 INFO example.go:11 std log json (zap)
					- 单行简要信息格式
					- 结构化日志格式（json） 方便大数据分析，一般用在生产环境中
				3. 日志是否支持轮转、单文件不能太大， 压缩，切割
				4. 日志包是否支持hook， gorm
			其他的需求：
				是否支持颜色显示
				是否兼容表中的log
				error打印到error文件，info打印到info文件
				error能否发送到其他的监控软件， 统计一个metrics错误指标
				error是否能支持发送到jaeger

		其他需求：
			高性能
			并发安全
			插件化： 错误告警，发邮件 sentry
			参数控制

		基于zap封装
	*/
}
