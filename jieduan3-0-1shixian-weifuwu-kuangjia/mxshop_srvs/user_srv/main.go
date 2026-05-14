package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	// 专门用于健康检查的grpc接口
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"mxshop_srvs/user_srv/global"
	"mxshop_srvs/user_srv/handler"
	"mxshop_srvs/user_srv/initialize"
	"mxshop_srvs/user_srv/proto"
	"mxshop_srvs/user_srv/utils"

	// consul 的 go客户端
	"github.com/hashicorp/consul/api"
)

func main() {
	// go语言内置的flag包，来解析命令行参数，解析用户传入的ip和端口号启动
	IP := flag.String("ip", "0.0.0.0", "ip地址") // 默认值为0.0.0.0，第三个参数："ip地址"，帮助说明文字用 ./program -help 就可以看到
	Port := flag.Int("port", 0, "端口号")         // 默认值为0，表示随机端口号

	//初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	// 打印配置信息
	zap.S().Info(global.ServerConfig)
	// 触发解析命令行参数解析方法----将结果注入到变量中
	flag.Parse()
	zap.S().Info("ip: ", *IP)
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	zap.S().Info("port: ", *Port)

	// 【grpc-server启动】1. 实例化一个grpc的server
	server := grpc.NewServer()
	// 【grpc-server启动】2. 注册服务
	proto.RegisterUserServer(server, &handler.UserServer{})
	// 【grpc-server启动】3. 监听端口
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	// 注册服务健康检查
	// 把 gRPC 官方提供的健康检查服务 注册到 gRPC 服务器
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的consul用的检查对象
	check := &api.AgentServiceCheck{
		// 必须使用GRPC字段.内部健康检查时用
		GRPC:                           fmt.Sprintf("192.168.0.103:%d", *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	registration.ID = serviceID
	registration.Port = *Port
	registration.Tags = []string{"imooc", "bobby", "user", "srv"}
	// 外部服务发现使用
	registration.Address = "192.168.0.103"
	registration.Check = check
	//1. 如何启动两个服务
	//2. 即使我能够通过终端启动两个服务，但是注册到consul中的时候也会被覆盖

	// 注册这个grpc服务到consul中
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	go func() {
		// 【grpc-server启动】4. 启动服务
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}
