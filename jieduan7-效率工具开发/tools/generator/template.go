package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"strings"
)

// go1.16 之后，Go 官方引入了一个新的特性：embed 包，go:embed 允许你在编译时将指定template.go.tpl文件的内容嵌入到 Go 程序中，成为一个字符串或字节数组。这样，你就可以直接在代码里使用这些文件的内容，而不需要在运行时去读取文件了。
//
//go:embed template.go.tpl
var tpl string

// rpc GetDemoName(*Req, *Resp)
type method struct {
	Name    string // GetDemoName
	Num     int    // 一个 rpc 方法可以对应多个http请求
	Request string // *Req
	Reply   string // *Resp

	// http rule
	Path         string
	PathParams   []string
	Method       string
	Body         string
	ResponseBody string
}

func (m *method) HandlerName() string {
	return fmt.Sprintf("%s_%d", m.Name, m.Num)
}

// HasPathParams 是否包含路由参数
func (m *method) HasPathParams() bool {
	paths := strings.Split(m.Path, "/")
	for _, p := range paths {
		if len(p) > 0 && (p[0] == '{' && p[len(p)-1] == '}' || p[0] == ':') {
			return true
		}
	}

	return false
}

// initPathParams 转换参数路由 {xx} --> :xx
func (m *method) initPathParams() {
	paths := strings.Split(m.Path, "/")
	for i, p := range paths {
		if p != "" && (p[0] == '{' && p[len(p)-1] == '}' || p[0] == ':') {
			paths[i] = ":" + p[1:len(p)-1]
			m.PathParams = append(m.PathParams, paths[i][1:])
		}
	}

	m.Path = strings.Join(paths, "/")
}

type service struct {
	Name     string
	FullName string

	Methods   []*method
	MethodSet map[string]*method
}

func (s *service) execute() string {
	if s.MethodSet == nil {
		s.MethodSet = make(map[string]*method, len(s.Methods))

		for _, m := range s.Methods {
			m := m // TODO ?
			s.MethodSet[m.Name] = m
		}
	}

	buf := new(bytes.Buffer)
	tmpl, err := template.New("http").Parse(strings.TrimSpace(tpl))
	if err != nil {
		panic(err)
	}

	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}

	return buf.String()
}

func (s *service) ServiceName() string {
	return s.Name + "Server"
}

func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// GoCamelCase camel-cases a protobuf name for use as a Go identifier.
//
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
// copy from protobuf/internal/strings.go
// 它是一个「字符串格式转换工具」
// 作用：把 proto 里的下划线命名 → 转换成 Go 语言的 大驼峰命名
// 遇到下划线 _ 删掉
// 下划线后面的字母变大写
// 第一个字母一定大写
// 小写字母串保持不变
func (s *service) GoCamelCase(str string) string {
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	var b []byte
	// 遍历字符串的每一个字符
	for i := 0; i < len(str); i++ {
		c := str[i] // 当前字符
		switch {
		//			情况 1：遇到 . 且后面跟着小写字母 → 跳过这个点
		case c == '.' && i+1 < len(str) && isASCIILower(str[i+1]):
			// Skip over '.' in ".{{lowercase}}".
		case c == '.':
			b = append(b, '_') // convert '.' to '_'
		case c == '_' && (i == 0 || str[i-1] == '.'):
			// Convert initial '_' to ensure we start with a capital letter.
			// Do the same for '_' after '.' to match historic behavior.
			b = append(b, 'X') // convert '_' to 'X'
		case c == '_' && i+1 < len(str) && isASCIILower(str[i+1]):
			// Skip over '_' in "_{{lowercase}}".
		case isASCIIDigit(c):
			b = append(b, c)
		default:
			// Assume we have a letter now - if not, it's a bogus identifier.
			// The next word is a sequence of characters that must start upper case.
			if isASCIILower(c) {
				// 小写字母减去 32，就变成了对应的大写字母！
				c -= 'a' - 'A' // convert lowercase to uppercase
			}
			b = append(b, c)

			// Accept lower case sequence that follows.
			for ; i+1 < len(str) && isASCIILower(str[i+1]); i++ {
				b = append(b, str[i+1])
			}
		}
	}
	return string(b)
}
