package main

import (
	"encoding/json"
	"fmt"
	"io" // Go 标准输入输出工具包
	"net/http"
)

type ResponseData struct {
	// Code int    `json:"code"`
	// 想要用json包生成JSON 字符串，必须首字母大写且给必须字段加上 Tag
	Data int `json:"你好"` // 这叫 Tag（标签），是 Go 特色语法！作用：告诉 json 序列化 / 反序列化时，字段对应 JSON 里的名字
	// 即把 json 里的 data 字段 映射到 ResponseData 结构体的 Data 字段
}

func Add(a, b int) int {
	fmt.Println("add")
	fmt.Println("hello world")
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/%s?a=%d&b=%d", "add", a, b))
	if err != nil {
		fmt.Println("请求失败:", err)
		return 0
	}
	defer resp.Body.Close() // 必须关！

	// 读取返回 body： HTTP 响应里的「响应体数据」全部读取出来，变成一段字节数组（[] byte）
	body, _ := io.ReadAll(resp.Body) // 把流里的所有数据一次性读完，读完后返回：[] byte（字节数组）
	// 字节数组可以表示任何数据包括中文，按照UTF-8编码，一个中文字符占3个字节等
	// 字节数组类型：数组里的每一个数字就代表一个字符 如：
	// 123 → {
	// 34  → "
	// 100 → d
	// 97  → a
	// 116 → t
	// 97  → a
	// 34  → "
	// 58  → :
	// 52  → 4
	// 125 → }
	fmt.Println(body, string(body), "=========") // [123 34 100 97 116 97 34 58 52 125] {"data":4} =========

	// 反序列化 JSON → 结构体
	var rspData ResponseData
	err = json.Unmarshal(body, &rspData) // body：JSON 原始数据,&rspData：必须传指针！ 才能把值填进去
	// json.Unmarshal 只接收 []byte类型 ，不能识别 string
	fmt.Println(rspData, "====1=====", err) // {4} ====1===== <nil>
	// 返回结果
	return rspData.Data // 返回结果 数字4
}

// client端
func main() {

	fmt.Println(Add(1, 3))

}
