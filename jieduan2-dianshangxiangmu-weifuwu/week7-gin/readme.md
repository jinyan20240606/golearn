# 第7周 gin快速入门

## 1-1&2 gin的helloworld体验

> 见gin_start/ch01和ch02
> 

是一个http的web服务框架，go语言写的


## 1-3 gin的路由分组

见gin_start/ch03

- Gin 路由匹配顺序：谁写在前面，谁优先！
```go
"/list"
"/:id/xxx"
// /:id 是通配，list 可以被匹配成 id="list" ----- 两个路由规则重叠了


// 解决办法（超级简单，2 种）
// 方法 1：把静态路由写在 动态路由 前面（最常用）-- 这样访问 /list 就会优先匹配列表！
// 方法 2：给动态路由加前缀（最安全、永远不冲突）
goodsGroup.GET("/list", goodsList)
goodsGroup.GET("/detail/:id/:action/add", goodsDetail)
```

## 1-4 获取url中的变量和校验

见gin_start/ch04

### 获取header头值的相关细节
1. go的net/http标准库行为
   
net/http 包在接收请求时，会把所有 header key 统一规范化：

原始：content-type / CONTENT-TYPE / Content-Type

存储：Content-Type（固定格式）

自定义头也一样：

原始：x-user-id: 123

存储：X-User-Id: 123

1. Gin 里两种获取方式（关键区别）
```go
func handler(c *gin.Context) {
    // 1. 直接拿 map：key 是规范后的大小写（首字母大写）
    fmt.Println(c.Request.Header["Content-Type"]) // 能拿到
    fmt.Println(c.Request.Header["content-type"]) // 拿不到！

    // 2. 用 .Get() 或 c.GetHeader()：大小写不敏感
    fmt.Println(c.Request.Header.Get("content-type")) // ✅ 能拿到
    fmt.Println(c.GetHeader("content-type"))            // ✅ 推荐
}
```


## 1-5 获取get和post表单信息

见gin_start/ch05

## 1-6 gin返回json和protobuf

见gin_start/ch06

**可以返回3种格式**
1. c.JSON返回json，json方法的2参也是可以配置的，默认是map[string]interface{}，也可以直接传结构体，gin会自动转换成json格式
2. 也可以用c.ProtoBuf返回protobuf:`c.ProtoBuf(http.StatusOK, &proto.User{})`
   1. 接口返给前端的是原始的protobuf二进制数据，前端是无法直接解析的，除非前端也用protobuf解析，否则就是一串乱码，适合内部服务之间通信
      1. 浏览器 / 前端默认把它当成文本解析,但它其实是二进制压缩格式,所以浏览器会不通过 Protobuf 解码 → 全是乱码
      2. 前端安装解析库，然后用 Protobuf JS 库 解码，才可以
      3. 内部服务grpc自动解析
   2. 什么时候才用 c.ProtoBuf？只有一种场景：
      1. 前后端约定使用 Protobuf 通信！
      2. 必须满足：
         1. 后端：返回 Protobuf 二进制
         2. 前端：用 Protobuf JS 库 解码
         3. 前后端共用同一份 .proto 文件
3. 返回c.PureJSON(200, gin.H{})：通常情况下，JOSN会将特殊的HTML字符替换为对应的unicode字符，比如标签符号< >替换为\u003c,如果想原样输出json，可以使用PureJSON

## 1-7 登录的表单验证

- 见gin_start/ch07


1. 表单的基本验证，go常用使用`go-playground/validator`库验证参数
   1. Gin 框架内部已经自带了 go-playground/validator
   2. 所以gin中可以直接使用验证器，如c.ShouldBind或c.ShouldBindUri()
2. 需要在参数绑定的结构体字段上设置tag.比如，绑定格式为json，需要这样设置json:"fieldname"
   1. Gin 参数绑定来源严格隔离：json ≠ form ≠ query ≠ uri**
      1. 只写 json:"user" → 只能收 JSON
      2. 写 json:"user" form:"user" → JSON + 表单都能收
      3. 只收一种会导致另一种传参方式直接失败
   2. tag内的验证规则：具体看文档全，下面是重点注意的规则
      1. 支持跨字段验证
3. 此外，Gin还提供了两套绑定方法:
   1. Must bind
      1. Methods - Bind, BindJSON, BindXML, BindQuery, BindYAML。
      2. Behavior-这些方法底层使用MustBindwith，
      3. 特点：如果存在绑定错误直接返回响应和状态码，请求将被以下指令中i上 c.AbortWithError(400, err).SetType (ErrorTypeBind), 响应状态f代码会被设置为400, 请求头Content-Type被设置为text/plain;charset=utf-8.注意，如果你试图在此之后设置响应代码，将会发出一个警告[GIN-debug][WARNING]HeadersWere already written. Wanted to override status code 400 with 422, l果你希望更好地控制行为，请使用ShouldBind相关的方法
   2. Should bind (----常用----)
      1. Methods - ShouldBind, ShouldBindJSON, ShouldBindXML, ShouldBindQuery,ShouldBindYAML
      2. Behavior -这些方法底层使用ShouldBindWith，
      3. 特点：如果存在绑定错误，则返回错误不直接返回响应和状态码，由开发人员可以正确处理请求和错误。

当我们使用绑定方法时，Gin会根据Content-Type推断出使用哪种绑定器，如果你确定你绑定的是什么,你可以使用MustBindwith或者Bindingwith.
你还可以给字段指定特定规则的修饰符，如果一个字段用binding:“required"修节资源网绑定时该字段的值为空，那么将返回一个错误，

## 1-8 注册的表单验证

- 继续见gin_start/ch07

## 1-9 表单验证错误翻译成中文

- 继续见gin_start/ch07

gin内置的validator库支持国际化，可以将错误信息翻译成中文

## 1-10 表单中文翻译的json格式化细节

- 继续见gin_start/ch07

## 1-11 自定义gin中间件

- 继续见gin_start/ch08

- gin内置中间件：Logger、Recovery、BasicAuth、Static
- 第三方必备：CORS、JWT、RequestID、限流(防止接口被刷)、Gzip

## 1-12 通过abort终止中间件后续逻辑的执行

- 奇怪的现象
- c.Request.Headers() 的header头的大小写默认规范化问题注意下

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 1. 测试1：只用 return → 拦不住！---- 依然会执行后面的路由链中间件，，，会返回msg：1和msg：hello
func testReturnBlock(c *gin.Context) {
	println("中间件1：我执行了，准备 return")
	// ❌ 错误：return 只能退出当前函数，不能阻止后面执行
    c.JSON(http.StatusOK, gin.H{
				"msg": "1",
			})
	return
}

// 2. 测试2：只用 Abort → 能拦住后续 ---- 不会执行后面的路由链中间件，，，，，只会返回msg：1
func testAbort(c *gin.Context) {
	println("中间件2：我执行了，准备 Abort")
    c.JSON(http.StatusOK, gin.H{
				"msg": "1",
			})
	c.Abort() // ✅ 中断整个执行链
	// ⚠️ 注意：本行依然会执行！Abort不退出当前函数
	println("中间件2：我还在执行！")
}

// 3. 测试3：正确写法 → Abort + return。   -----  拦截了后续路由链中间件，，，，，只返回msg：1
func testAbortReturn(c *gin.Context) {
	println("中间件3：我执行了，准备拦截")
    c.JSON(http.StatusOK, gin.H{
				"msg": "1",
			})
	// ✅ 正确：中断链条 + 退出函数
	c.Abort()
	return

	// 下面永远不会执行
	println("我不会执行")
}

// 真正的接口处理函数
func apiHandler(c *gin.Context) {
	println("✅ 接口逻辑执行了！！！")
	c.JSON(http.StatusOK, gin.H{"msg": "hello"})
}

func main() {
	r := gin.Default()

	// ============ 测试1：return 拦不住 ============
	// 访问：http://localhost:8080/test1
	r.GET("/test1", testReturnBlock, apiHandler)

	// ============ 测试2：Abort 能拦住 ============
	// 访问：http://localhost:8080/test2
	r.GET("/test2", testAbort, apiHandler)

	// ============ 测试3：Abort + return ============
	// 访问：http://localhost:8080/test3
	r.GET("/test3", testAbortReturn, apiHandler)

	_ = r.Run(":8080")
}
```