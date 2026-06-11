package srv

import (
	"fmt"
	gpb "mxshop/api/goods/v1"
	"mxshop/app/goods/srv/config"
	v12 "mxshop/app/goods/srv/internal/controller/v1"
	db2 "mxshop/app/goods/srv/internal/data/v1/db"
	"mxshop/app/goods/srv/internal/data_search/v1/es"
	v1 "mxshop/app/goods/srv/internal/service/v1"
	"mxshop/gmicro/core/trace"
	"mxshop/gmicro/server/rpcserver"

	"mxshop/pkg/log"
)

// goods 商品微服务的启动入口函数 NewGoodsRPCServer
// 职责：
// 初始化链路追踪（OpenTelemetry）
// 初始化 MySQL 数据层、ES 检索层工厂
// 组装 service 业务层、controller RPC 接口层
// 创建 gRPC RPC 服务、注册路由
// 返回 RPC Server，外部直接 Start() 即可启动服务
func NewGoodsRPCServer(cfg *config.Config) (*rpcserver.Server, error) {
	//初始化open-telemetry的exporter
	trace.InitAgent(trace.Options{
		cfg.Telemetry.Name,
		cfg.Telemetry.Endpoint,
		cfg.Telemetry.Sampler,
		cfg.Telemetry.Batcher,
	})

	//有点繁琐，wire， ioc-golang
	// 当前是手动依赖注入，层级一多组装代码变冗长。
	// 服务变复杂（新增缓存、MQ、中间件）后，NewGoodsRPCServer 会越来越长。
	// 优化方向：引入 Google Wire 做自动依赖注入，自动生成组装代码；使用 Wire / ioc-golang 自动依赖注入，减少手动组装代码。
	dataFactory, err := db2.GetDBFactoryOr(cfg.MySQLOptions)
	if err != nil {
		log.Fatal(err.Error())
	}

	//构建，繁琐 - 工厂模式
	searchFactory, err := es.GetSearchFactoryOr(cfg.EsOptions)
	if err != nil {
		log.Fatal(err.Error())
	}

	srvFactory := v1.NewService(dataFactory, searchFactory)
	goodsServer := v12.NewGoodsServer(srvFactory)
	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	grpcServer := rpcserver.NewServer(rpcserver.WithAddress(rpcAddr))

	gpb.RegisterGoodsServer(grpcServer.Server, goodsServer)

	//r := gin.Default()
	//upb.RegisterUserServerHTTPServer(userver, r)
	//r.Run(":8075")
	return grpcServer, nil
}
