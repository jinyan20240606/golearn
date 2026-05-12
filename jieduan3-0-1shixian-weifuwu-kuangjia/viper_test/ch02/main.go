package main

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

//如何将线上和线下的配置文件隔离
//不用改任何代码而且线上和线上的配置文件能隔离开

type MysqlConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	ServiceName string      `mapstructure:"name"`
	MysqlInfo   MysqlConfig `mapstructure:"mysql"`
}

// 读取环境变量
func GetEnvInfo(env string) bool {
	// 作用：让 viper 自动读取系统环境变量
	viper.AutomaticEnv()
	return viper.GetBool(env)
	// 刚才设置的环境变量 想要生效 我们必须得重启goland
	// 环境变量是进程启动时加载的Goland 启动后，不会自动刷新环境变量
}

func main() {
	// 根据部署代码时，自动读取系统环境变量，自动读取不同配置
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("viper_test/ch02/%s-pro.yaml", configFilePrefix)
	// 如果是debug环境：读取debug环境变量相关的配置文件
	if debug {
		configFileName = fmt.Sprintf("viper_test/ch02/%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()
	//文件的路径如何设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	serverConfig := ServerConfig{}
	// 反序列化到serverConfig结构体
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	fmt.Println(serverConfig)
	fmt.Printf("%V", v.Get("name"))

	//viper的功能 - 动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file channed： ", e.Name)
		_ = v.ReadInConfig()
		// 反序列化绑定到结构体
		_ = v.Unmarshal(&serverConfig)
		// 打印结构体
		fmt.Println(serverConfig)
	})

	time.Sleep(time.Second * 300)
}
