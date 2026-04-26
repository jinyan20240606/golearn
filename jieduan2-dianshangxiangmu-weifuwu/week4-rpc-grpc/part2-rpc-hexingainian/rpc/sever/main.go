package main

import (
	"encoding/json" // Go 官方的 JSON 工具包
	"fmt"
	"net/http"
	"strconv"
)

// 使用内置的http  server 库

func main() {
	// http://localhost:8080/add?a=1&b=2
	// 创建一个http server
	// 返回的格式化：json {"data":3}
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		// 响应写入器：http.ResponseWriter
		// *http.Request：请求对象，代表结构体的指针类型，接收指针，共用同一份数据，省内存、效率高，默认是不加*即值传递修改不会影响原数据
		_ = r.ParseForm() // 解析参数
		fmt.Println("url-path:", r.URL.Path)
		a, _ := strconv.Atoi(r.FormValue("a"))
		b, _ := strconv.Atoi(r.FormValue("b"))
		w.Header().Set("Content-Type", "application/json")
		// Marshal：意思是序列化 → 把 Go 数据map类型变成 JSON 格式
		// 反序列化方法：json.Unmarshal
		jData, err := json.Marshal(map[string]int{"你好": a + b})
		if err != nil {
			http.Error(w, `{"error":"json encode failed"}`, http.StatusInternalServerError)
			return
		}
		w.Write(jData)
	})

	// 必须加 :，否则监听会失败。
	http.ListenAndServe(":8080", nil)

}
