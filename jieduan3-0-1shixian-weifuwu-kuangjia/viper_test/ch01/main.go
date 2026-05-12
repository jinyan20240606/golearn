package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	ServiceName string `mapstructure:"name"`
	Port        int    `mapstructure:"port"`
}

func main() {
	v := viper.New()
	// 文件的路径如何设置
	v.SetConfigFile("viper_test/ch01/config.yaml")
	// 用v的方法 读进来
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	// 读进来之后，viper会把yaml文件中的数据解析成一个map[string]interface{}类型的值，存储在viper内部的结构体中，
	// 我们可以通过viper提供的Unmarshal反序列化方法把这个map[string]interface{}类型的值转换成我们定义的ServerConfig结构体类型的值
	serverConfig := ServerConfig{}
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	fmt.Println(serverConfig)
	// 或者不反序列化到结构体，直接用viper的Get方法获取值，v.Get("name")返回的是interface{}类型的值
	fmt.Printf("%V", v.Get("name"))
}
