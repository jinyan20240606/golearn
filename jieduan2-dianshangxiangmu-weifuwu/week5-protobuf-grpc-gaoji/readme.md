# protobuf和grpc高级进阶

grpc底层就是基于protobuf的

## 1-1 protobuf的基本类型和默认值

> protobuf的官方文档可参考

1. 定义一个消息类型
2. 标量数值类型
   1. 一个标量消息字段可以含有一个如下的类型protobuf中的每个数值类型都对应不同后端语言中的具体数值类型
3. 默认值：当一个消息被解析的时候，如果你不传值的话，那么这个字段的值将被设置为默认值，程序不会缺值报错，对于不同类型指定如下:
   1. 对于strings，默认是一个空string
   2. 对于bytes，默认是一个空的bytes
   3. 对于bools，默认是false
   4. 对于数值类型，默认是0
   5. 对于枚举，默认是第一个定义的枚举值，必须为0;
   6. 消息类型(message)，域没有被设置，确切的消息是根据语言确定的，详见generatedcode guide
      1. 对于可重复域的默认值是空(通常情况下是对应语言中空列表)。注:对于标量消息域，一旦消息被解析，就无法判断域释放被设置为默认值(例如，例如boolean值是否被设置为false)还是根本没有被设置。你应该在定义你的消息类型时非常注意。例如，比如你不应该定义boolean的默认值false作为任何行为的触发方式。也应该注意如果一个标量消息域被设置为标志位，这个值不应该被序列化传输。
      2. 查看generated codeguide选择你的语言的默认值的工作细节。
4. 枚举
5. 使用其他消息类型
6. ...等等，详细查看搜索文档


## 1-2 option go_package的作用

- option go_package = "./;helloworld"; // 👈 就加这一行！
- // ./ 表示生成到当前目录 和 ;helloworld 表示生成的 Go package 叫 helloworld
- 如：option go_package = "common/stream;v1";
- 不加斜线，是相对目录，生成到当前的common/stream目录下，v1子目录作为包目录，包名叫v1
- 一般微服务中，都统一生成到公共目录下，多个微服务之间，可以公共引用

## 1-3 proto文件同步时的坑

建grpc_proto_test目录


一般实际开发中，只在这个目录下写proto和client代码不写server，去调用其他项目如python写的微服务中的server代码


1. 两个项目必须要同步proto文件，才能代码一致
```go
message HelloRequest {
  string name = 1; // 1是编号 ，不是值,复制proto时两边同步proto时，这个编号一定不能搞错，错了消息就顺序就会反
  int32 age = 2;
}
```


## 1-4 proto文件中import另一个proto文件

方便代码拆分复用

如part1-protobuf/grpc_proto_test/proto/base.proto，写了公共的Empty消息体，在helloworld.proto中引用

- 这时候，编译命令需要编译多个文件：`protoc --go_out=. --go-grpc_out=. helloworld.proto base.proto`


- protobuf中 还提供了许多内置消息类型，如：`google.protobuf.Empty`，使用见文件
- 内置消息类型见包源码目录里：github.com/golang/protobuf/ptypes下面是所有支持的内置类型
- 内置类型然后在client中使用时，引入路径也是github.com下的源码目录，见 `grpc_proto_test/client/client.go`


## 1-5 嵌套的message对象


见helloworld中的嵌套的Result类型相关使用 和client中的引用方式

## 1-6 protobuf中的enum枚举类型

## 1-7 map类型

## 1-8 使用protobuf内置的timestamp类型

时间戳类型

## 1-9 grpc的metadata机制

gRPC让我们可以像本地调用一样实现远程调用，对于每一次的RPC调用中，都可能会有一些有用的数据，而这些数据就可以通过metadata来传递，metadata是以key-value的形式存储数据的，其中key是string类型，而value是[]string，即一个字符串数组类型。metadata使得client和server能够为对方提供关于本次调用的一些信息，就像一次http请求的RequestHeader和ResponseHeader一样。http中header的生命周周期是一次http请求，那么metadata的生命周期就是一次RPC调用。

- 类型`type MD map[string][]string` // value为string切片类型
- **metadata的核心规则**
  - Metadata = gRPC 的请求头
  - 必须放在 context 里传递

### 新建metadata

- 新建的时候可以像创建普通的map类型一样，使用new关键字进行创建，也可以使用Pairs方法
  - metadata.New()
  - metadata.Pairs()
```go
// 使用New的方式：创建空 metadata
md := metadata.New(nil)

// 往里面放值:用 Set 会覆盖，不能追加
md.Set("token", "123456")
md.Set("uid", "10086")
// 往同一个 key 追加多个值：-----> "hobby": ["eat", "sleep"]
md.Append("hobby", "eat")
md.Append("hobby", "sleep")

// 或者直接赋值
md := metadata.MD{
    "token": []string{"123456", "654321"},
    "uid":   []string{"10086"},
}

// 或者使用Paris的方式
md := metadata.Pairs(
    "token", "123456", // 给同一个 key 放多个值的方式实现数组写法
    "token", "789", // token: [tk1, tk2, tk3]
    "uid", "10086",
    "version", "1.0.0",
)

// 获取全部值 → 返回 []string
hobbies := md["hobby"]

// 或者用 Get → 返回第一个值
first := md.Get("hobby") // "eat"
```


### 发送metadata

```go
// 1. 构造 metadata（单值 + 多值都可以）
md := metadata.Pairs(
    "token", "123456",
    "uid", "10086",
    "hobby", "eat",    // 多值
    "hobby", "sleep",  // 多值
)

// 2. 把 md 放入 context（关键方法）
ctx := metadata.NewOutgoingContext(context.Background(), md)

// 3. 调用 gRPC 方法时，把 ctx 传进去
resp, err := client.SayHello(
    ctx, // 👈 必须传这个 ctx
    &proto.HelloRequest{Name: "test"},
)
```

### 接收metadata
```go

func (s *server) SayHello(
    ctx context.Context,
    req *proto.HelloRequest,
) (*proto.HelloReply, error) {

    // 👇 从 ctx 里取出 metadata
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, errors.New("获取元数据失败")
    }

    // =====================================
    // 取值方式1：Get → 拿第一个值（常用）
    // =====================================
    token := md.Get("token") // string
    uid := md.Get("uid")     // string

    // =====================================
    // 取值方式2：直接取数组 → 拿全部值（多值）
    // =====================================
    hobbies := md["hobby"] // []string{"eat", "sleep"}

    fmt.Println("token:", token)
    fmt.Println("uid:", uid)
    fmt.Println("hobbies:", hobbies)

    return &proto.HelloReply{Message: "ok"}, nil
}
```

## 1-10 grpc拦截器-go版

- 拦截器，可以理解为grpc中间件，在grpc调用之前和之后执行。
  - 客户端可以加拦截器
  - 服务端也加个拦截器
- 简单模式和流模式的拦截器写法也不一样


见part1-protobuf/grpc_interpreter文件夹



- 有一个现成的开源中间件，可以参考：https://github.com/grpc-ecosystem/go-grpc-middleware
- 也可以自己写一个中间件库

## 1-11 通过拦截器和metadata实现grpc的auth认证

- 实现方式：直接在客户端的拦截器中添加认证信息的metadata元数据，在服务端的拦截器中验证元数据
- 另一种实现方式：go语言中可以不用在client中使用原生拦截器那样做，可以换一种方式做，go语言中针对认证场景有一个专门的intercepter
  - 可以把代码写的更加简单呢