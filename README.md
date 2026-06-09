# golearn

笔记参考百度网盘的6-go语言工程师体系课

## go语言原理相关

1. Node.js 是单线程异步非阻塞，io时必须用 await，否则会卡死。
2. Go 是多线程同步阻塞，天生支持高并发，io时可以直接同步写法，性能无敌！
3. 并发安全问题
   1. 在 Node.js 中，普通的单例写法不用担心多线程并发问题（因为它是单线程的）；但如果你在单例的创建过程中有异步操作，就必须加上 Promise 锁来防止重复初始化
      1. 绝大部分没有异步代码情况下你不需要考虑像 Go 那样的并发（多线程）安全问题，同一时刻，Node.js 只会执行一段 JavaScript 代码，
   2. Go不一样，go必须考虑并发安全问题， go的多 Goroutine 会真正并行执行。当多个 Goroutine 同时读写同一块共享资源（如全局变量、Map、切片）且没有加锁保护时，就会发生数据竞争（Data Race）
      1. 如单例模式的创建也得考虑并发安全问题，如`jieduan6-开发规范设计模式单元测试/pattern/ch02/main.go`

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
4. 省略号语法2种用法
   1. 参数用：表示可变参数类型，作为类型
      1. `func callback(msgs ...*MessageExt) ` msgs就是可变参数类型的变量
   3. 传值用：表示代表把切片打散，一个个传进去，作为值用
      1. `callback(msgList...) // ✅ 正确！打散切片传入`
5. go中的存在多个defer时，defer的执行顺序是后进先出，最后一个定义的defer先执行
6. 创建文件时，一定要记得设置文件权限，避免文件权限类错误
7. select case语法： 是 Go 里专门配合 channel 使用的并发语法，作用是：同时等待多个 channel 操作，谁先就绪就执行谁
   1. 普通的 switch 不一样：switch 判断值，select 判断 channel 通信事件
8. 什么是优雅退出：即最关键：server.Shutdown (ctx) 到底做什么？
  - 温柔版关闭流程：
    - 不再接受新连接
    - 已经进来的请求继续跑完
    - 跑完后再退出
    - 如果 ctx是一个超时上下文， 超时（比如 5 秒），就强制关闭
      - 5 秒内优雅关闭完 → 正常退出
      - 5 秒还没关闭完 → 强制退出
  - 对比暴力关闭：
    - server.Close() → 直接断开所有请求，用户报错
9. wait group用法：主要用于goroutine的执行等待，Add方法必须与Done方法配套成对使用
  - 你执行一次 wg.Add(1)，计数器 +1
  - 你执行一次 wg.Done()，计数器 -1
  - 当计数器 = 0 时，wg.Wait() 就会停止阻塞，继续往下走

### 类型的表达式运算
1. `int / int` → 结果是 **整数**（小数直接截断，不会四舍五入）
   1. 只能`float64(totalErr) / float64(total) `
2. Go 语言有一个规则：
   1. 如果两个数类型不一样，会自动把「范围小 / 表示能力弱」的类型，转换成「强」的那个
   2. float64 比 int 强
   3. 所以计算`float64与int相互乘除时` → 结果是 float64（带小数，正确）时，int 会 自动变成 float64

### 数据库相关

1. - **重点记录**：一般不建议不用外键，只留逻辑外键，数据的一致性由应用层（代码）保证，不由数据库保证
   1.  不是不使用外键，数据库中只要定义了外键，就默认有约束规则
   2.  gorm中只有在标签里 显式写出 constraint: 相关内容，GORM 才会生成 物理外键约束；否则，哪怕你写了 foreignKey:、references:，都只是逻辑关联，不会生成数据库物理外键！，如见`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop_srvs/goods_srv/model/goods.go`
2. 数据库中的字符集：utf8mb4 = 完整版 UTF-8（能存 emoji 😄），默认实现的utf8mb4是不支持emoji的
   1. 全世界所有标准 UTF-8 都能存 Emoji，只有 MySQL 的 utf8 是假的、阉割版，存不了
3. 加索引的原则：我们需要根据这个字段查询时候才会加， 1. 会影响插入性能 2. 会占用磁盘
   1. 一般不要随意加，会严重拖慢性能
### 常用遗忘方法
2. json序列化和反序列化方法
   1. `序列化：只接收1个参数（任何类型）转成json字节切片类型`
      1. `b, _ := json.Marshal(&categorys)`
      2. json.Marshal(任何类型) 返回值永远是 []byte
         1. json.Marshal 可以接收 Go 里几乎所有常用类型转成json字节切片
         2. Marshal = 把 Go 数据 → 变成 JSON 字节（序列化）
   2. `反序列化：只接收2个参数（[]byte (字节切片), 指针变量）`
      1. `var data []model.Category`-`json.Unmarshal(b, &data)`
      2. Unmarshal = 把 JSON 字节 → 变回 Go 数据（反序列化）
         1. json.Unmarshal 第一个参数 只接受一种类型：[]byte (字节切片)，2参是指针变量
         2. 转换成 Go 里几乎所有常用类型：bool，int / uint / float，string，slice / 数组，map，结构体，interface{}
            ```go
               // 转成map类型
               bs := []byte(`{"name":"张三","age":20}`)

               var m map[string]interface{}
               json.Unmarshal(bs, &m)

               fmt.Println(m["name"]) // 张三
               // 转成 结构体类型
               type User struct {
                  Name string `json:"name"`
                  Age  int    `json:"age"`
               }

               bs := []byte(`{"name":"张三","age":20}`)

               var u User
               json.Unmarshal(bs, &u)

               fmt.Println(u.Name) // 张三
            ```
         3. 将原始json字符串解析：把原始json字符串 转 [] byte，用 json.Unmarshal 解析
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
   6. time.Now().Sub(time.Now()) ：获取当前时间与之前时间的差值时间，返回值：time.Duration
   7. time.Now().Add(time.Duration(1)) ：获取当前时间加上指定时间间隔的时间，返回值：time.Time
   8. time.Time：它是一个对象 / 结构体，里面存了：年、月、日、时、分、秒、纳秒、时区……，所有时间操作都靠它
   9. time.After (时间) = 定时触发的通道，它返回一个 只读 channel，
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
5. 随机数方法
   1. rand.Uint64 ()：生成一个 0 ~ 极大的随机无符号整数（随机 64 位数字）
      1. 计算机里：
         1. 8 位 = 0~255
         2. 16 位 = 0~65535
         3. 32 位 = 0~42 亿
         4. 64 位 = 0 ~ 18446744073709551615
      2. rand.Uint64 () % 20
         1. 0 ~ 19 之间的随机整数
         2. 不会小于 0
         3. 不会大于等于 20


#### 字符串相关方法

1. strconv方法
   1. strconv.ParseInt(id, 10, 64) // 把字符串 → 转成 int64 类型，2参是进制，3参是目标位数
   2. strconv.Itoa(123) // 返回 "123"
   3. strconv.Atoi("123") // 返回 int, error
   4. strconv.ParseInt("123", 10, 32)
   5. strconv.ParseFloat("3.14", 64)
   6. strconv.FormatFloat(float64(rsp.Total), 'f', 2, 64) // 将float64 → 转成字符串
      1. 要转换的浮点数（订单金额）
      2. 'f' 字符格式只占1个字节：普通小数（不是科学计数法），用字符传，节省内存占用提升性能
      3. 保留 2 位小数时，默认执行的就是 四舍五入
      4. // 源数据是 float64
2. strings.TrimPrefix(字符串, 要去掉的前缀)
   1. strings.SplitN(字符串, 分隔符, 切成几段)

#### 事务和锁

1. 只要涉及批量修改或组合处理数据，一定要有事务
   1. 单机本地事务
      1. 控制本地机器里读写事务
   2. 分布式事务
      1. 涉及一组事务含跨服务读写资源，如创建订单里的各个模块处理逻辑
2. 只要涉及并发读写同一份数据，一定要有锁的介入
   1. 单机本地锁
   2. 分布式锁



### shell语法积累

1. 参考`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop_srvs/goods_srv/start.sh`

### 在课程中提到的docker技巧

1. 服务器重启，自动启动容器`docker run -d --name nginx -p 80:80 --restart=always nginx`
## go语言中的常见规定
### 代码规范 
uber开源的代码规范:https://github.com/xxjwxc/uber_go_guide_cn
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