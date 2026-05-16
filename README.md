# golearn

笔记参考百度网盘的6-go语言工程师体系课



## 易混积累

1. go中双引号跟单引号区别
   1. 单引号：表示字符类型，定义的类型是字节类型
   2. 双引号：表示字符串类型，定义的类型是string
2. utf-8语法规则：Go 字符串的底层编码是 UTF-8，不同字符的 UTF-8 字节数固定：
   1. 英文字母、数字、符号（ASCII 字符）：占 1 个字节；
   2. 中文、日文等东亚字符：占 3 个字节；
   3. 极少数特殊字符（如 emoji）：占 4 个字节；
3. Go 里直接用 string(数字) 不是把数字转成字符串，而是把数字当成 Unicode 字符编码转成字符
   1. 比如：string(65) → 不是 "65"，而是字符 A，
   2. 所以绝对不能用 string(i) 来把数字转成数字字符串！
   3. 正确写法（2 种，任选）
      1. strconv.Itoa(i)
      2. fmt.Sprintf（一行搞定）：result := fmt.Sprintf("%s %d %d", req.Data, i, time.Now().Unix())
4. - **重点记录**：一般不建议不用外键，只留逻辑外键，数据的一致性由应用层（代码）保证，不由数据库保证
   1.  不是不使用外键，数据库中只要定义了外键，就默认有约束规则
   2.  gorm中只有在标签里 显式写出 constraint: 相关内容，GORM 才会生成 物理外键约束；否则，哪怕你写了 foreignKey:、references:，都只是逻辑关联，不会生成数据库物理外键！，如见`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop_srvs/goods_srv/model/goods.go`
5. 数据库中的字符集：utf8mb4 = 完整版 UTF-8（能存 emoji 😄），默认实现的utf8mb4是不支持emoji的
   1. 全世界所有标准 UTF-8 都能存 Emoji，只有 MySQL 的 utf8 是假的、阉割版，存不了

### 常用遗忘方法
1. strings.SplitN(字符串, 分隔符, 切成几段)
2. json序列化和反序列化方法
   1. `// 序列化：必须1个参数（只给数据）`
   2. `b, _ := json.Marshal(&categorys)`
   3. `// 反序列化：必须2个参数（给JSON + 给空切片）`
   4. `var data []model.Category`
   5. `json.Unmarshal(b, &data)`
3. time方法
   1. `time.Unix(int64(value.BirthDay), 0)`
      1. func Unix(sec int64, nsec int64) Time
      2. 作用： 把秒级时间戳 + 纳秒，转成 Go 的 time.Time 对象
      3. 10位秒级时间戳如： timestamp := int64(1735689600)
   2. time.Now()：打印结果展示字符串为2025-04-01 15:30:45 +0800 CST
      1. 返回值：time.Time 结构体，含义：系统当前完整时间对象（年月日时分秒时区都有）2025-04-01 15:30:45 +0800 CST
   3. time.Now().Unix() ：获取当前秒级时间戳，返回值：int64
   4. time.Now().UnixMilli() ：获取当前毫秒级时间戳，返回值：int64
   5. time.Now().UnixNano() ：获取当前纳秒级时间戳，返回值：int64
   6. time.Time：它是一个对象 / 结构体，里面存了：年、月、日、时、分、秒、纳秒、时区……，所有时间操作都靠它
4. Go JSON库的标签常见语法
   1. 
   ```go
   type LoginForm struct {
	User     string `json:"user,omitempty" binding:"required,min=3,max=10"`
	Password string `json:"password,-" binding:"required"`
   }
   // 最常用的 2 个：
   // omitempty：字段为空时，不返回给前端
   // -：忽略这个字段，不序列化、不返回
   ```
### 在课程中提到的docker技巧

1. 服务器重启，自动启动容器`docker run -d --name nginx -p 80:80 --restart=always nginx`
## go语言中的常见规定

### Go 语言中命名和结构体初始化规定

#### 核心规则：大小写决定可见性

Go 语言没有 public、private、`protected关键字，而是通过标识符的首字母大小写来控制访问权限：

1. 大写开头（Exported/导出）：

含义：公共的，可以被其他包访问。
适用对象：变量、常量、函数、方法、结构体字段、接口等。
示例：Name, GetName, Course。
你的代码问题：course.Course 中的字段 name 是小写开头，因此它是私有的。main 包无法直接访问或初始化它。

2. 小写开头（Unexported/未导出）：

含义：私有的，只能在当前包内部访问。
示例：name, getName, course (如果它是局部变量则无所谓，如果是包级变量则其他包不可见)。

#### 其他重要的 Go 语言规定

除了可见性规则，Go 还有以下几类关键规定，特别是在使用结构体和包时：

A. 结构体初始化规定
键值对初始化（Key-Value）：

必须使用导出字段（大写开头）。
字段名必须与结构体定义完全一致。
示例：User{Name: "Alice", Age: 30}。
顺序初始化（Positional）：

可以忽略字段名，按定义顺序赋值。
限制：只能用于所有字段都导出或者在同一包内的情况。如果结构体中有未导出字段，跨包时不能使用这种简写方式，必须使用键值对且只填导出字段。
示例：User{"Alice", 30} （不推荐，易出错）。
部分初始化：

使用键值对时，可以只初始化部分字段，未初始化的字段为零值（如 "", 0, nil）。
示例：User{Name: "Alice"} （Age 默认为 0）。


B. 包导入与命名规定
1. 包名唯一性：

在同一文件中，导入的包别名不能重复。
如果两个包名相同（如都叫 user），必须给其中一个起别名：
```go
import (
    user1 "example.com/project/v1/user"
    user2 "example.com/project/v2/user"
)
```

2. 导入路径：

必须使用 go.mod 中定义的模块路径。
不支持相对路径（如 ../user）当项目初始化了 go.mod 后。
3. 空白导入（Blank Import）：

使用 _ 导入包，仅为了执行该包的 init() 函数，而不使用其任何导出标识符。
示例：import _ "github.com/lib/pq" （常用于数据库驱动注册）。


C. 方法与接口规定
1. 接收者（Receiver）：

方法必须定义在同一包内的类型上。
不能为其他包的基本类型（如 int, string）添加方法，只能为自定义类型或指针添加。
2. 接口实现：

Go 是隐式实现接口。只要一个类型实现了接口中的所有方法，它就自动实现了该接口，无需 implements 关键字。
通常建议接口名以 er 结尾（如 Reader, Writer, Stringer）。

D. 错误处理规定

1. error 类型：

Go 没有异常机制（try-catch），而是通过返回 error 类型来处理错误。
惯例：最后一个返回值通常是 error。`func OpenFile(name string) (*File, error)`

2. panic 与 recover：

panic 用于严重错误（如数组越界、空指针解引用），会导致程序崩溃。
recover 只能在 defer 函数中捕获 panic。
规定：不要在普通业务逻辑中使用 panic 代替错误返回。