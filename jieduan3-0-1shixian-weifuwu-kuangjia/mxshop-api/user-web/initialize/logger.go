package initialize

import "go.uber.org/zap"

func InitLogger() {
	logger, _ := zap.NewDevelopment()
	// 替换全局的logger
	zap.ReplaceGlobals(logger)
}
