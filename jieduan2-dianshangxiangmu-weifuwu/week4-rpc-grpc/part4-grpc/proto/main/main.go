package main

import (
	"encoding/json"
	helloworld "golearn/part4-grpc/proto"

	"fmt"

	"github.com/golang/protobuf/proto"
)

type Hello struct {
	Name   string   `json:"name"`
	Age    int      `json:"age"`
	Course []string `json:"course"`
}

func main() {

	req := helloworld.HelloWorldRequest{
		Name:    "小王",
		Age:     18,
		Courses: []string{"go", "gin", "微服务"},
	}

	// 使用proto.Marshal序列化下，rsp为[]byte字节切片类型
	rsp, _ := proto.Marshal(&req)

	fmt.Println(string(rsp)) // 输出是 乱码的小王字符串
	fmt.Println(len(rsp))    // 输出是 7

	jsonStruct := Hello{"小王", 18, []string{"go", "gin", "微服务"}}
	jsonRsp, _ := json.Marshal(jsonStruct)
	fmt.Println(string(jsonRsp)) // 输出是 {"Name":"小王"}
	fmt.Println(len(jsonRsp))    // 输出 16

	// 类比于json的Marshal 和 Unmarshal

	// 对比看：明显是proto的压缩比更高，但是不易读

	// 虽然直接打印序列化的值看不懂，但是借助proto反序列化还是能看懂，只有proto自己能看懂
	newReq := helloworld.HelloWorldRequest{}

	_ = proto.Unmarshal(rsp, &newReq)
	fmt.Println(newReq.Name) // 输出是 小王
	fmt.Println(newReq.Age)  // 输出是 18
	fmt.Println(newReq.Courses)

	/**
		输出结果：
	 * 小王gogin       微服务
		30
		{"name":"小王","age":18,"course":["go","gin","微服务"]}
		60
		小王
		18
		[go gin 微服务]
	*/
}
