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