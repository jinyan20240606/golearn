package handler

// 名称冲突，需要加前缀

const HelloServiceName = "handler.HelloService"

type HelloService struct {
}

func (p *HelloService) Hello(request string, reply *string) error {
	*reply = "hello:" + request // *reply 指针，指针变量 存的是地址，*指针 才是取值 / 赋值
	return nil
}
