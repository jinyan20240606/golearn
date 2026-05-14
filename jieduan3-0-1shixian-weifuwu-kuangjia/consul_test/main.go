package main

import (
	"fmt"
	// consul的go客户端
	"github.com/hashicorp/consul/api"
)

// 注册服务
func Register(address string, port int, name string, tags []string, id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.1.103:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		// 这块本地测试也不能写localhost，必须写具体ip地址
		HTTP:                           "http://192.168.1.102:8021/health",
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = name
	registration.ID = id
	registration.Port = port
	registration.Tags = tags
	registration.Address = address
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	client.Agent().ServiceDeregister()
	if err != nil {
		panic(err)
	}
	return nil
}

// 服务发现：获取所有服务
func AllServices() {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.1.103:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().Services()
	if err != nil {
		panic(err)
	}
	for key, _ := range data {
		fmt.Println(key)
	}
}

// 服务发现：按服务名过滤
func FilterSerivice() {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.1.103:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	// 调用 client.Agent().Services() 打印出一个服务完整结构，就能看到所有可过滤字段
	// 下面过滤方法是用字段名key-val过滤

	data, err := client.Agent().ServicesWithFilter(`Service == "user-web"`)
	if err != nil {
		panic(err)
	}
	for key, _ := range data {
		fmt.Println(key)
	}
}

func main() {
	_ = Register("192.168.1.102", 8021, "user-web", []string{"mxshop", "bobby"}, "user-web")
	AllServices()
	FilterSerivice()
	fmt.Println(fmt.Sprintf(`Service == "%s"`, "user-srv"))
}
