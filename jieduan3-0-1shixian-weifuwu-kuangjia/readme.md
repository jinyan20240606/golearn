# 阶段3 从0到1实现完整的微服务框架

> 整体课件项目代码就见`mxshop_srvs`

## 8周 用户服务的grpc服务

- 该周只有1章-用户服务的service开发
- 代码见`mxshop_srvs/user_srv`

### 1-1 定义用户表结构

> 创建`user_src/model`专门存放用户的表结构字段

### 1-2 同步数据库的表结构

> 创建`user_src/model/main`文件夹， 专门用来同步数据库的表结构

下面语句会创建一个user用户信息表，注释不用这个的话，直接上面db.save也会默认创建表
`_ = db.AutoMigrate(&model.User{}) //此处应该有sql语句`

### 1-3 密码字段的md5加密

1. 一般数据库中密码字段，不能存明文密码，一旦数据库丢失，密码就丢了，一般都是密文保存，而且需要密文不可反解
2. 加密算法一般分为以下几种方式：
   1. 对称加密：加密和解密用的是同一把钥匙，这种一把钥匙泄漏风险也很大，也不能满足不可反解的要求
   2. 非对称加密：一般采用非对称加密，加密和解密用不同的钥匙，但是他不能满足密码不可反解的要求
   3. md5 信息摘要算法：最常用的是这个，它不能反解，它严格上说不是加密算法，而是信息摘要算法，但是一般用做密码加密
3. 密码如果不可以反解，用户忘记了密码找回密码怎么办？？
   1. 首先不能反解的话，我们拿到数据库存的密文，也是无法反解的，就算反解了发给用户邮件万一邮件被拦截，万一丢了泄漏了也不安全
   2. 所以一般是给用户一个链接，让用户去重置一个新密码

#### md5信息摘要算法加密

首先哈希算法就是摘要算法，它是大类，md5只是Hash的一种。下面md5的特性就是哈希算法的特性，不可反解。

1. Hash 家族常见成员，全都属于 Hash 摘要算法：
   1. MD5（最老，不安全），值固定32个字符
   2. SHA-1：值固定40个字符
   3. SHA-256：值固定64个字符
   4. SHA-512：值固定128个字符
   5. bcrypt、PBKDF2、Argon2（密码专用哈希）
   6. Hash算法的基本特点：
      1. 任意长度输入 → 固定长度输出
      2. 单向不可逆，不能解密还原原文
      3. 用来做：校验完整性、密码存储、签名
2. md5:摘要算法可以将任意长度的字符串转换成固定长度的16进制字符串
   1. 压缩性：任意长度的数据，算出md5值的长度都是固定的，永远是32个16进制字符
      1. md5的底层就是输出128位二进制串，32位16进制字符串
   2. 容易计算：从原数据计算出MD5值很容易
   3. 抗修改性：对原数据进行任何修改，哪怕1个字节，md5值差异也很大
   4. 强碰撞：想找到2个完全不同的数据，使得它们的MD5值相同，这是不可能的
   5. 相同内容 → 永远相同 MD5，不同内容 → 几乎不可能相同
   6. 不可逆性：不能反解，单向不可逆（不能从密文还原原文）
3. md5盐值加密
   1. 加盐
      1. 通过生成随机数和md5生成字符串进行组合
      2. 数据库同时存储md5值和salt盐值，验证正确性使用salt进行md5即可


- go中md5加密代码示例见：user_srv/model/main/main.go 的 genMd5 函数
- 默认的md5加密，得到的密文是非常不安全的因为是不可反解的，会提前将任意常见密码用md5加密下存下来，相当于md5值和你的密码是一一对应的，因此可以被暴力破解-彩虹表直接反向映射到，所以一般都会加盐加密

### 1-4 md5盐值加密解决用户密码安全问题


- 使用开源库已经封装了密码加密加盐的方法："github.com/anaskhan96/go-password-encoder"
  - 库内部源码，就是这么写的：PBKDF2 本身不是哈希算法，它是一个 “加密框架 / 流程”！它自己不会算哈希，必须靠 HashFunc（SHA512/SHA1）来干活！PBKDF2 只是一个 “重复加密的流程框架”，它需要一个真正的哈希算法来执行每一次加密
    ```go
    derivedKey := pbkdf2.Key([]byte(password), 
        salt, 
        opts.Iterations,  // 迭代次数
        opts.KeyLen,      // 密钥长度
        opts.HashFunc,    // SHA512
    )
    ```
  - 该函数内部自动计算生成一个加密后 的随机盐值和 加密后的密文
    - `salt, encodedPwd := password.Encode("generic password", options)`
  - 验证用法：验证用户的密码对不对  
    - `password.Verify("原始密码", passwordInfo[2]/*盐值*/, passwordInfo[3]/*密文密码*/, options)`
- 问题：有盐值，那这个盐值存在哪里呢，用户登录后用户名密码得到后，咋取到盐值进行校验呢
  - 一般不建议salt存在用户表中，一般是直接存在密文密码中,存储到数据库的密码字符串，格式是：`$pbkdf2-sha512$随机盐值$加密后的密文`
  - 当用户名登录后，用它的原始密码和数据库中的密文密码包含的盐值和密文提取出来进行方法验证对比，如果相同，则验证成功

### 1-5 定义proto接口

> 定义user_srv/proto/user.proto文件

### 1-6 用户列表接口

- 见 user_srv/handler/user.go 文件，来实现proto中的用户列表接口的具体定义
  - 需要引入gorm数据库实例去查询数据
- 建立 user_srv/global/global.go 全局变量，里面定义了数据库连接等公共引用方法使用

### 1-7 通过id和mobile查询用户

- 见 user_srv/handler/user.go 文件，来实现proto中的GetUserByMobile和GetUserByMobile方法的具体定义
	
    // Go 里面几乎所有：
	// gRPC 响应
	// 业务返回值
	// 较大的结构体
	// 全部统一返回指针，不返回值

### 1-8 新建用户接口
- 见 user_srv/handler/user.go 文件的CreateUser方法

### 1-9 修改用户和校验密码接口

- 见 user_srv/handler/user.go 文件的UpdateUser方法

### 1-10 通过flag启动grpc服务

启动这个grpc服务，测试下前面写的grpc的接口

> 见mxshop_srvs/user_srv/main.go文件

1. 使用go语言内置的flag包，来解析命令行参数，解析用户传入的ip和端口号动态启动grpc-server
   1. `ip := flag.String("ip", "0.0.0.0", "ip地址")：`
      1. 第一个参数："ip"，命令行参数的名字 `./program -ip=192.168.1.100`
      2. 第二个参数：默认值
      3. 第三个参数：参数的描述 用 `./program -help` 命令查看参数的描述
   2. `flag.Parse()  // 必须触发执行解析`
   3. `fmt.Println("IP：", *ip)`
      1. 小细节: flag.方法 返回的是 *string 指针,所以使用时要 加 *
   4. 使用时用`main.exe  -ip=192.168.1.100`


### 1-12 测试用户微服务接口

- 见user_srv/tests 文件夹


## 9周 用户服务的web服务

> 整体见 mxshop-api目录

### 1章 web层开发-基础项目架构

#### 1-1 新建项目和目录结构

- 新建mxshop-api/user-web项目 和对应的目录结构合理划分

#### 1-2 go高性能日志库
- 特点
  - 性能极高：比 gin 默认 log 快 10~100 倍
  - 结构化日志：JSON 格式，方便排查问题
  - 可输出到文件：自动写日志文件
  - 可分级：Debug / Info / Warn / Error / DPanic / Panic / Fatal
  - 生产环境标准库：公司 Go 项目几乎都用 zap
```go
// 安装
go get go.uber.org/zap
// 文件使用
package main

import "go.uber.org/zap"

func main() {
	// 生产环境配置
	logger, _ := zap.NewProduction()
	defer logger.Sync() // 刷新缓冲区

	// 打日志
	logger.Info("服务启动成功",
		zap.String("ip", "0.0.0.0"),
		zap.Int("port", 8080),
	)

	logger.Error("数据库连接失败",
		zap.Error(fmt.Errorf("connection timeout")),
	)
}



// ========= 替换 Gin 框架默认日志 ========
router := gin.Default()

// 把 gin 的日志替换成 zap
router.Use(ginzap.Ginzap(log, time.RFC3339, true))
```

- 使用zap库，来替换gin框架自带的日志中间件库，来实现日志的输出
- Gin 默认 Logger 有什么问题？（重点）
    1. 不能输出到文件，只能打印控制台，不能写文件，生产环境没法用。
    2. 格式固定，不能自定义，只能是它那一种格式，不能改成 JSON。
    3. 没有日志级别，没有 Debug / Info / Error 区分，所有日志混在一起。
    4. 性能一般，小项目没问题，高并发下性能不如 zap。
    5. 不能按天切割、不能自动清理，日志越来越大，撑爆磁盘。
    6. 不方便日志收集，不是结构化 JSON，ELK 之类的工具不好解析
- Zap 有两种 logger：
  - Logger（原味，高性能）
    - 调用：logger.Info("msg", zap.String("k", "v"))
    - 特点：最快、无反射、类型安全；但写起来啰嗦。
  - SugaredLogger（加糖，易用,性能略低一点点）
    - 获取：sugar := logger.Sugar() 或 sugar := zap.NewExample().Sugar()
    - sugar.Info("msg", "k", "v")
    - sugar.Infof("name=%s age=%d", "tom", 18)
- 获取全局loagger的简便写法，初始化一次不用子孙传递获取，直接用简写方法就能获取到全局实例
  - zap.S() 相当于logger.Sugar() .简写，直接获取全局的
  - zap.L() 相当于logger.Logger()，简写，直接获取全局的
```go
package main

import "go.uber.org/zap"

func main() {
	// 1. 初始化一次
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger) // 变成全局
	defer logger.Sync()

	// 2. 直接用！！！
	zap.L().Info("服务启动成功") // L() = 全局 logger
}
```

#### 1-3 zap的文件输出

```go
package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// 1. 创建日志文件（没有会自动创建，有就覆盖）
	logFile, err := os.OpenFile("run.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("日志文件创建失败")
	}

	// 2. 配置日志格式
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 时间格式：2025-01-01 10:00:00
	encoderConfig.TimeKey = "time"                        // 时间字段名

	// 3. 核心：日志输出到【文件】
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 日志格式：JSON
		logFile,                               // 输出目标：文件
		zap.InfoLevel,                         // 日志级别：Info及以上
	)

	// 4. 创建 logger（带文件名+行号）
	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync() // 程序退出前把日志刷入文件

	// ============== 使用 ==============
	logger.Info("服务启动成功", zap.String("ip", "0.0.0.0"))
	logger.Error("数据库连接失败", zap.Int("code", 500))
}
```

#### 1-4&5 集成zap和路由初始到gin的启动过程

- web服务我们使用gin框架，将gin框架安装到我们项目中
- gin-web服务框架的默认启动用法，就不能用简单的demo方式启动了，需要适配现在的工程化的方式启动，我们需要按照我们的目录结构来
  - gin的初始化单独封装在`initialize/router.go`模块中，暴露Router方法，在main文件中初始时调用`initialize.Routers()`，且自顶向下传递得到的路由实例
  - api的统一放到api目录下分模块见子文件
    - api接口的代码专门存放web-server的接口api逻辑，处理用户的请求，调用grpc的服务端接口，返回结果给前端
    - 先建立`api/user.go文件`
  - router路由相关的专门放到router目录下维护，其他目录下都是服务
    - router是路由入口，负责调用api接口
    - 先建立`router/user.go文件`
- 全局logger初始化也封装在`initialize/logger.go`中，暴露Logger方法，在main文件中初始时调用`initialize.Logger()`

#### 1-6&7 gin调用grpc服务

- 见`mxshop-api/user-web/router/user.go` 和 `mxshop-api/user-web/api/user.go`

- 先实现的api/user.go的GetUserList接口方法

#### 1-8 go的配置文件管理库：viper库

Viper = Go 一站式配置管理工具，遵循 12-Factor，统一管理所有配置源，无需硬编码、无需关心格式。

核心功能（全覆盖）

✅ 多格式支持：YAML、JSON、TOML、HCL、INI、env、properties

✅ 多配置源：配置文件、环境变量、命令行参数、远程配置（etcd/Consul）、默认值

✅ 热加载：监听文件变化，自动重新读取（热更新）

✅ 结构体绑定：Unmarshal 到 Go 结构体，类型安全

✅ 优先级合并：自动按优先级覆盖，无需手动处理

✅ 大小写不敏感：Key 不区分大小写

```js
一、核心特点

多格式兼容支持 YAML、JSON、TOML、INI、HCL、Properties 等主流配置格式，不用改代码随意切换格式。
多配置来源统一管理配置来源全覆盖：本地配置文件、环境变量、命令行参数、内存默认值、远程配置中心（etcd/Consul）。
自动配置优先级内置固定优先级：代码 Set > 命令行 > 环境变量 > 配置文件 > 远程配置 > 默认值高层自动覆盖低层，不用自己手写覆盖逻辑。
支持配置热加载可监听配置文件变化，自动重新加载，不用重启服务就能更新配置。
结构体绑定支持直接把配置 Unmarshal 绑定到结构体，类型安全，告别零散 GetString/GetInt 硬编码。
键名大小写不敏感配置 key 大小写不区分，书写更随意，减少大小写报错。
层级配置支持完美支持嵌套层级配置（如 app.port、db.host），适合复杂项目结构。
零侵入、易集成无复杂依赖，接入简单，所有 Go 项目（单体、微服务、CLI 工具）都能直接用。
```

- Viper 优势：多格式、多来源、自动优先级、热加载、结构体绑定，把 Go 项目配置从「手写硬编码」变成「标准化、优雅、可维护」的统一方案
- viper练习目录见：`viper_test/ch01目录`

#### 1-9 viper的配置环境开发环境和生产环境

- 目录见：`viper_test/ch02目录`

#### 1-10 viper集成到gin的web服务中

> 见mxshop-api/user-web

- 创建2个配置文件config-debug和config-pro.yaml，集成到该web项目中
- 支持配置后，`api/user.go`的grpc客户端初始化时下就不用硬编码host和端口号了
- 创建全局配置文件：`mxshop-api/user-web/config/config.go`
- 接着在哪里读取初始化全局配置文件，见单独的模块目录中：`mxshop-api/user-web/initialize/config.go`
  - 在这里封装全局的初始化方法读取全局配置文件，如初始化grpc客户端


### 2章 web层开发-用户接口开发



## 10周 服务注册发现，配置中心，负载均衡