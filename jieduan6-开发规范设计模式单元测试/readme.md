# 阶段6 开发规范设计模式单元测试
## 21周 开发规范和go基础扩展
### 1章 开发规范
#### 1-1 后续学习前的思维介绍

之前算是广度，接下来是重点深度能力，学习工程化思维和技术深度
#### 1-2 课程要用到的基本开发工具说明

1. go在1.18后提供了2个重要的功能：模糊测试和范型，我们后面基于最新版本学习，也会讲新语法

#### 1-3 项目开发有哪些规范要遵守
5. go·微服务代码目录规范
   1. a.微服务项目和单体服务的目录不同点
   2. b.微服务应该如何管理目录
   3. 6.代码发布规范
6. a.go项目的发布步骤
   1. i.静态代码扫描。
   2. i.代码自动格式化
   3. i.代码自动运行单元测试
   4. iv.go vet检查竞态
   5. v.自动编译
   6. vi.镜像上传

#### 1-4 项目开发流程
略
#### 1-5&6 git代码分支管理&commit规范
#### 1-7 go的代码规范

> uber开源的代码规范:https://github.com/xxjwxc/uber_go_guide_cn
- 代码规范一下全部接受不容易，经常看，养成习惯即可
- 规范不代表权威，某个规范自己思考一下，不一定就正确，要结合自己的需求具体情况具体分析
- 简单给大家介绍几个
  - 零值Mutex是有效的
  - erors比较重要，后面有专门的章节讲解
  - 这里有一些我们后面会有专门的章节讲解，所以建议大家学习完课程以后再来看一下这里的规范
- 下吗是uber的代码规范示例：
```go
// 结构体实现接口的代码规范
type Handler struct {}
// 提前用这个赋值语法，判断是否满足接口
// 解释：我要把一个 *Handler 类型的 nil 指针，赋值给 http.Handler 接口，看下是否报错
var _ http.Handler = (*Handler)(nil) // *Handler是指针类型，后面的括号：是类型强转语法，将nil转化类型为*Handler，跟int(64)一样
func (h *Handler) ServeHTTP(w http.ResponseWriter,r http.Request){}
```

#### 1-8 go目录规范

很多目录规范是随着某个框架而确定的，并不是语言本身可以决定目录规范，比如python中的django目录，java的spring目录规范，但是go目前还没有出现spring一样一统天下的框架，所以目录规范也并不统一，但是在某种程度上还是有大家的共识的，
> 我们以uber的目录规范来做一下说明。参考: https://github.com/golang-standards/project-layout/blob/master/README_zh.md

1. Go 项目根目录下不要自己建 /src，那是 Java 习惯，在 Go 里会乱掉路径和工具链。
   1. Go 原本工作区就有`$GOPATH/src/`,项目里再搞一个 src就不合适
#### 1-9 微服务该采用monorepo还是多个repo？
#### 1-10 微服务的项目目录规范
```js
项目根目录
├── api          # 专门放所有服务对外暴露的接口,接口定义：proto 文件、API 文档、swagger
├── app          # 所有微服务的核心业务代码（每个服务一个文件夹）
│   ├── order    # 订单服务
│   ├── user     # 用户服务
│   └── pkg      # 服务内部公共包
├── build        # 构建脚本：Dockerfile、构建配置
├── cmd          # 项目启动入口（main 函数所在），每个服务一个启动入口
├── configs      # 配置文件：yaml、toml、json
├── docs         # 文档：设计文档、接口文档
├── internal     # 【全局私有包】整个项目能用，不是开源给外部用的，外部不能引用 import
├── init         # 系统初始化脚本：systemd 脚本
├── pkg          # 全局公共库，公共的（外部可 import）（所有服务都能使用），工具、通用组件，希望开源 / 给别人用
├── scripts      # 运维脚本：启动、停止、部署、打包
├── test         # 测试文件：压力测试、集成测试
├── third_party  # 第三方依赖、私有库，不通过 go mod 管理的私有包
├── tools        # 工具类：代码生成、脚本工具
├── examples     # 示例代码
└── go.mod        # Go Module 依赖管理
```

#### 1-11 govet进行代码检测

- go内部官方提供的静态代码分析工具是`go vet ./main.go` go vet 专门检查 编辑器 / 编译器 检查不出来 的错误
- 静态：静态分析（go vet）
  - 不运行代码
  - 只扫描代码文本
  - 编译前 / 编译中检查
  - 发现潜在风险
  - 速度极快
#### 1-12 golangci-lint进行代码检测
golangci-lint = Go 语言的静态代码检查集成工具，比go vet 更加强大，类似于eslint

- 特点：
  - 不运行代码 → 纯静态分析
  - 集成了几十种 linter
    - go vet
    - golint
    - staticcheck
    - errcheck
    - unused
    - 等等...
  - 速度快、配置简单、可自定义规则
  - 企业级 Go 项目标配
- 看下如何集成到编辑器本地环境和项目中
  - 可以在编辑器中的save保存文件时的钩子中设置运行golangci-lint，到时候自己搜
  - 本地项目中创建.golangci.yml文件，编辑器钩子和项目中运行时都基于项目级配置运行
### 2章 go基础知识扩展

主要讲下go开发中容易犯的错 和一些新特性

#### 2-1 map初始化容易犯的错

```go
// 错误代码
var course map[string]string
course["name"]="go体系课"


// var course map[string]string 只声明了 map，值是 nil（空指针）
// 对 nil map 赋值 直接 panic：assignment to entry in nil map
// slice 可以不用 make 直接用（自动扩容），但 map 必须手动 make 初始化
// 1. 必须 make 初始化 map，slice不用初始化
course := make(map[string]string)
// 2为什么？
// map 本质是个指针，只要经过 make 初始化，就不再是空指针！，它会变成一个指向真实哈希表内存的有效指针
// slice 本质是结构体，零值合法，append 能自动扩容。
```
#### 2-2 常见错误：结构体的空指针
```go
type User struct {
    Name string
}

func main() {
    // 只声明指针，没初始化 → u = nil----- 直接变成空指针，使用就报错
    var u *User

    // 错误！空指针访问字段 → 直接崩溃！
    u.Name = "张三"
}

// u 是 nil 指针，没有指向任何内存
// 你让它去找 Name 字段 → 找不到 → 程序直接挂！
 // 正确 2：字面量初始化
    u2 := &User{}
```

- 空指针不能点，一点就崩溃！
- map 要 make，结构体指针要 new！
- 只有 slice 特殊，nil 也能 append！
#### 2-3 常见错误-使用对循环迭代器变量的引用

```go
// 错误场景1:
package main
import "fmt"
func main() {
	var out []*int
	// for循环的临时变量会复用
	for i := 0; i < 3; i++ {
		out = append(out, &i) // 把 i 的地址存进去
	}
	for _, value := range out {
		fmt.Println(*value) // 打印出来全是 3！
	}
}
// 运行结果: 3 3 3 ---- 严重的错误！！！
// 正确写法：
package main
import "fmt"
func main() {
	var out []*int
	for i := 0; i < 3; i++ {
		// 关键：创建一个新临时变量，每次都分配新地址
		x := i
		out = append(out, &x)
	}
	for _, value := range out {
		fmt.Println(*value)
	}
}

// 原因：for 循环里的变量 i 只有一个，地址永远不变，每次循环只是修改它的值！
// 你存的永远是同一个地址，最后 i=3 退出循环，所以打印全是 3。

// 错误场景2: 如果直接在协程里用，所有协程都会读到最后一个值！
for _, id := range goodsID {
	// 错误！直接在闭包里用循环变量 id
	go func() {
		fmt.Println("正在查询商品:", id)
	}()

    // 正确解法1:传参形式
	go func(id uint64) {
		fmt.Println(id)
	}(id)


    // 正确解法2：临时变量形式
	// 关键：循环内新建临时变量！
	newID := id  
	go func() {
		// 直接用 newID，完全安全！
		fmt.Println(newID)
	}()
}

```

#### 2-4&5 什么是范型
go的1.18版本开始，go支持泛型。

跟ts一样，go中的范型传参语法是用[]中括号包裹的，ts是<>包裹的

- 代码详见`jieduan6-开发规范设计模式单元测试/generics/main.go`
```go
package main

import "fmt"

func Add[T int | float64](a, b T) T {
	return a + b
}

func main() {
	// 自动推导出 T=int，
	fmt.Println(Add(1, 2)) 
    // 也可以指定参数类型
    fmt.Println(Add[int](1, 2))

	// 自动推导出 T=float64
	fmt.Println(Add(1.1, 2.2)) 
}


// 这是范型出现之前的老版麻烦写法
func IAdd(a, b interface{}) interface{} {
	switch a.(type) {
	case int:
		return a.(int) + b.(int)
	case int32:
		return a.(int32) + b.(int32)
	case float32:
		return a.(float32) + b.(float32)
	case float64:
		return a.(float64) + b.(float64)
	}
	return nil
}

func main() {
	fmt.Println(IAdd(1, 2))
	fmt.Println(IAdd(1.1, 2.2))
}
```

#### 2-6 范型的常见用法

范型不仅用在函数中，也可以用在其他类型中

- 代码详见`jieduan6-开发规范设计模式单元测试/generics/ch02`
#### 2-7 范型的错误用法

- 代码详见`jieduan6-开发规范设计模式单元测试/generics/ch02`中的错误用法常见约束

## 22周 设计模式和单元测试