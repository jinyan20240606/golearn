package main

import (
	// bytes 包 = 用来操作「字节数组」的工具包
	// 里面最最常用的，就是 bytes.Buffer
	// 作用：一块可写、可读、可拼接的内存缓冲区
	"bytes"
	// 模版库
	"html/template" // Go 标准库「文本模板」工具包
	"strings"
)

// 1. 定义模板字符串（就是你要生成的 Go 代码骨架）
// 双引号里的所有内容 = 模板
var tpl = `
type {{.Name}}HttpServer struct {
	server {{$.Name}}Server

	router gin.IRouter
}

func Register{{.Name}}HttpServer(server {{.Name}}Server, router gin.IRouter) {
	//我现在想用gin.Default,如果开发中我想使用qit
	g := &{{.Name}}HttpServer{server: server, router: router}
	g.RegisterService()
}


{{ range .Methods }}
func (g *{{ $.Name }}HttpServer) {{ .HandlerName }}(c *gin.Context) {
	var in {{ .Request }}
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := g.server.{{ .Name }}(c, &in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}
{{ end }}

func (g *{{$.Name}}HttpServer) RegisterService() {
{{ range .Methods }}
	g.router.Handle("{{ .Method }}", "{{ .Path }}", g.{{ .HandlerName }})
{{ end }}
}
`

// 2. 定义数据结构：给模板填的数据
type serviceDesc struct {
	Name    string   // 服务名 Greeter
	Methods []method // 方法列表
}

type method struct {
	Name    string // 方法名 SayHello
	Request string // 请求参数
	Reply   string // 返回值

	//http rule
	Path   string // HTTP 路径
	Method string // POST/GET
	Body   string
}

func (m *method) HandlerName() string {
	return m.Name + "_0"
}

func main() {
	//模板
	// 3. 创建一个缓冲区，用来存放最终生成好的代码，因为它tmpl.Execute方法只接受一个 能不断写入的接口（io.Writer）
	// ❌ string（不行）
	// 只能读
	// 不能写
	// 不能追加
	// 只能整体替换
	// ✅ bytes.Buffer（可以）
	// 可以写
	// 可以追加
	// 可以拼接
	// 可以不断往里填内容
	buf := new(bytes.Buffer)
	// 4. 创建一个名叫 http 的模板，把 tpl 字符串解析编译成可执行模板
	// New("http")：给模板随便起个名字，叫 http
	// 把字符串模板，编译成 Go 能执行的模板
	// Parse 做的事：
	// 读取模板字符串
	// 识别 {{}}
	// 检查语法是否正确
	// 变成可执行的模板结构
	tmpl, err := template.New("http").Parse(strings.TrimSpace(tpl)) // TrimSpace去掉模板字符串首尾多余的空格、换行
	if err != nil {
		panic(err)
	}
	// 5. 构造要填进模板的数据！
	// 这部分数据，实际开发中是从 proto 文件里解析出来的
	s := serviceDesc{
		Name: "Greeter", // 服务名 → 会替换模板里的 {{.Name}}

		Methods: []method{
			{
				Name:    "SayHello", // 方法名
				Request: "HelloRequest",
				Reply:   "HelloReply",
				Path:    "/v1/sayhello", // HTTP 路径
				Method:  "POST",         // 请求方法
				Body:    "*",
			},
		},
	}
	// 6. 执行模板：把数据 s 填进模板，输出到 buf
	err = tmpl.Execute(buf, s)
	if err != nil {
		panic(err)
	}
	// 7. 打印最终生成的代码！
	println(buf.String())
}
