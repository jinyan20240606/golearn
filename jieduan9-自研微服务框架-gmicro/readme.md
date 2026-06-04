# 阶段9-自研微服务框架-gmicro

- 代码目录见`jieduan9-自研微服务框架-gmicro/mxshop`

## 26周 三层代码结构

- 示例代码目录见
1章 3层代码结构规范

### 1-1 导入common和app包
见`jieduan9-自研微服务框架-gmicro/mxshop/pkg`下的common和app公共包，可以在mxshop/app应用目录各个微服务应用中快速引入初始化项目
### 1-2 通过app启动配置文件映射和flag映射
### 1-3 重构app启动项目
### 1-4 app启动的原理
### 1-5 已有代码存在哪些耦合

我们前面开发的mxshop_srvs项目和mxshop-api项目，每个handler中逻辑内部有很多耦合地方，handler直接依赖的grpc的代码，耦合了很多模块硬编码直接用：下面想换哪一个都得全量改代码

- RPC 想换（zRPC → 其他）
- ORM 想换（GORM → 原生 SQL / 其他）
- Web 框架想换（Gin → go-zero/kratos）
- 注册中心想换（Consul → Nacos/K8s）
- 缓存想换（Redis → 内存 /memcache）
- 一句话终极解决方案：你缺的是：接口抽象 + 依赖倒置（DI）+ 分层架构
- 只要把底层全部抽成接口，上层只依赖接口，不依赖具体实现 → 想换啥换啥，一行业务代码不改！
### 1-6 三层代码结构降低代码耦合
```js
controller层（参数校验，调用servicec接口）
    接收 HTTP / RPC 请求
    参数校验
    组装请求 → 传给 service
    返回统一响应
    只和前端交互
    暴露结构：VO (View Object) 视图对象 (Request/Response)，专门给页面接口展示用，美化、聚合、格式化数据
service层（具体的业务逻辑）
    真正的业务逻辑
    事务
    调用 data 层（DB/RPC/Cache）
    组合多个 data 数据源
    完全不关心底层用什么 GORM/RPC/Redis
    暴露结构：DTO / BO（Business Object）跨层传输数据，前端请求、接口返回都用它。
data层（数据库的接口）
    DB 操作（GORM / 原生 SQL）
    RPC 调用（zRPC/grpc）
    缓存（Redis）
    只做数据读写，无业务逻辑
    暴露结构：DO (Data Object)，和数据库表一一对应，纯粹存数据。
// 每层的对外暴漏的数据接口类型是不一样的
```

- 我们以原来的user服务和对应的user-web服务，为例，重构到当前的规范微服务目录架构`app/user`下
- 如在我们实际的微服务项目下`jieduan9-自研微服务框架-gmicro/mxshop/app/user`微服务下，首先分成srv和client两个大目录,每个目录就按照3层结构进行目录划分
- api目录下只放接口文档，如proto文件，swagger文档等，还可以加个版本号，作为二级目录（因为有可能不同的服务依赖我们不同的接口版本）

### 1-7 service层和data层的解耦

- 继续完善`jieduan9-自研微服务框架-gmicro/mxshop/app/user`服务的重构迁移，本节重点：将service层和data层的代码分离

- service层最终是要调data层，data下也是建个v1

### 1-8 DO、DTO、VO这些概念是什么

![alt text](image.png)
![alt text](image-1.png)

- 这些概念基于上面的3层结构进行划分的。
  - controller层主要负责接收请求和返回结果：前端请求进入系统时通常先绑定成 Request DTO，service处理完成后既可以返回 Response DTO 给 controller，再由 controller 转成 VO 响应给前端；也可以在简单项目里由 service 直接返回 VO。
  - service层内部主要使用 BO 承载业务语义和聚合数据，然后与 data 层通过 BO 转成 DO 或 PO 进行交互
  - data 层内部主要基于 DO 或 PO 最后再与外部依赖数据库服务进行 DAO 对接

1. DTO = 数据传输对象（跨层）
   1. Controller ↔ Service 跨层之间传输
   2. 请求 DTO：前端传过来的参数，controller 接收后传给 service
   3. 响应 DTO：service 返回给 controller 的数据，controller 再决定是否转成 VO 返回前端
   4. 更广义上，服务与服务之间的 RPC / HTTP 请求和响应对象，也都可以算 DTO
2. VO = 调用方展示专用
   1. 给前端页面看、给调用方展示”。它一般用于响应结果
3. BO（Business Object）业务对象
   1. 给谁用：Service 层（核心业务）
   2. 作用：承载业务逻辑、事务、聚合数据
   3. 特点：真正的业务核心，与前端 / 数据库无关
   4. 例子：商品 + 库存 + 价格 = 订单 BO
4. DAO（Data Access Object）数据访问接口
   1. 给谁用：DAO 层（数据访问层）
   2. 作用：定义数据库操作接口
   3. 特点：接口，不关心实现（可换 MyBatis/JPA/GORM）
5. DO（Data Object）数据对象
   1. 给谁用：DAO 实现层
   2. 作用：与数据库表一一对应
   3. 特点：纯数据库映射
6. PO 一般指 Persistent Object，DO 一般指 Data Object：在很多项目里两者都表示持久化层对象，通常和数据库表结构对应，实际常常不严格区分。区别更多取决于团队命名规范，而不是统一标准。为了避免混乱，项目里最好只保留一种命名
   1. 如果你们团队已经有 OrderDO，就别再来一个 OrderPO；如果已经叫 UserPO，那就统一都叫 PO

#### 一个经典业务场景：创建订单时整条链怎么流转

下面用“创建订单”这个场景，把 DTO、BO、DO、VO 一次串起来。重点不是只看字段，而是看每一段到底在干什么业务。

##### 1. 前端传参：用户发起“下单请求”

用户在页面上：
- 选了收货地址
- 选了商品和数量
- 选了优惠券
- 点击“提交订单”

这时候前端发来的只是“用户输入”，还不是完整订单：

```json
{
  "user_id": 1001,
  "address_id": 88,
  "coupon_id": 10,
  "items": [
    {"goods_id": 20001, "num": 2}
  ]
}
```

- 这一段的业务逻辑：
  - 前端只负责收集用户输入
  - 调接口把数据传给后端
  - 此时还不知道商品价格、库存、优惠券是否可用、最终实付金额是多少

##### 2. 进入 Request DTO：后端把前端参数接住

后端 Controller 收到请求后，先绑定成请求 DTO：

```go
CreateOrderDTO {
    UserID: 1001,
    AddressID: 88,
    CouponID: 10,
    Items: [{GoodsID: 20001, Num: 2}],
}
```

- 这一段的业务逻辑：
  - 做参数绑定
  - 做基础参数校验，比如 user_id、address_id、items 是否为空
  - 做最基础的合法性校验，比如数量必须大于 0
- DTO 这一层解决的是：你传得对不对

##### 3. Service 补全成 BO：进入真正的业务处理阶段

到了 Service，系统开始围绕“下单”做真正的业务计算。这里很关键：一个 DTO 在 service 层，完全可能被拆成多个 BO，不一定只对应一个 BO。
- DTO 是前端传过来的 “大包裹”，里面可能包含多个独立业务域的数据。
- 例子：用户提交 “创建订单 + 上传商品”
  - 前端只发 1 个请求 → 1 个 CreateOrderAndProductDTO
  - 但后端必须拆成 2 个独立 BO：
    - OrderBO（订单业务对象）
    - ProductBO（商品业务对象）
  - 因为店铺和商品是两个独立业务域，不能混在一个 BO 里

```go
CreateOrderDTO {
    UserID: 1001,
    AddressID: 88,
    CouponID: 10,
    Items: [{GoodsID: 20001, Num: 2}],
}
```

在 service 内部，可能会先拆成多个 BO：

```go
// 订单主业务 BO
OrderBaseBO {
    UserID: 1001,
    AddressID: 88,
    CouponID: 10,
    ReceiverName: "张三",
    ReceiverMobile: "13800000000",
    AddressDetail: "上海市浦东新区xxx路",
}
// 订单明细商品 BO
OrderGoodsBO {
   //  OrderID: 占位，当插入订单成功后，会追加
    GoodsID: 20001,
    GoodsName: "机械键盘",
    Price: 150,
    Num: 2,
    Stock: 99,
    Subtotal: 300,
}
```

如果项目不大，也可以聚合成一个总 BO：

```go
CreateOrderBO {
    UserID: 1001,
    AddressID: 88,
    CouponID: 10,
    ReceiverName: "张三",
    ReceiverMobile: "13800000000",
    AddressDetail: "上海市浦东新区xxx路",
    Items: [{
        GoodsID: 20001,
        GoodsName: "机械键盘",
        Price: 150,
        Num: 2,
        Stock: 99,
        Subtotal: 300,
    }],
    GoodsAmount: 300,
    CouponAmount: 20,
    FreightAmount: 8,
    PayAmount: 288,
}
```

- 这一段负责每个业务领域的业务逻辑：
  - 根据 address_id 查询地址，补全收件人姓名、手机号、详细地址
  - 根据 goods_id 查询商品，补全商品名称、当前单价、库存、是否上架
  - 校验库存是否充足
  - 计算商品总额、优惠金额、运费、实付金额
  - 把“用户输入”变成“系统已确认、可执行”的完整业务上下文
- BO 这一层解决的是：这笔业务到底该怎么做
- 这里的关键补充：
  - 一个 DTO 可以拆成多个 BO，因为 service 处理的往往不是“传参结构”，而是多个业务子问题
  - 比如下单业务里，地址信息、商品明细、金额计算，本身就可以拆成不同 BO

##### 4. 转成 DO 落库：把业务结果存进数据库

当 Service 完成校验和计算后，就会把 BO 转成 DO / PO，准备入库。这里也很关键：一个 BO 往往也不只对应一个 DO，因为数据库通常本来就是多张表。

例如下单时，至少可能拆成下面几个 DO：

```go
// 订单主表
OrderDO {
    UserID: 1001,
    ReceiverName: "张三",
    ReceiverMobile: "13800000000",
    AddressDetail: "上海市浦东新区xxx路",
    GoodsAmount: 300,
    CouponAmount: 20,
    FreightAmount: 8,
    PayAmount: 288,
    CouponID: 10,
    Status: 1,
}
// 订单明细表
OrderItemDO {
    OrderID: 90001,
    GoodsID: 20001,
    GoodsName: "机械键盘",
    Price: 150,
    Num: 2,
    Subtotal: 300,
}
```

- 这一段的业务逻辑：
  - 写入订单主表
  - 写入订单明细表
  - 扣减库存
  - 更新优惠券使用状态
  - 这些操作通常放在同一个事务里，任一步失败都要回滚
- DO 这一层解决的是：这些业务结果怎么安全、准确地保存下来
- 这里的关键补充：
  - 一个 DTO 进入 service 后，可能拆成多个 BO
  - 一个 BO 在持久化时，也完全可能拆成多个 DO
  - 在订单场景里非常常见，因为订单主表、订单明细表、优惠券记录表、库存流水表，本来就是不同存储对象

##### 5. Service 响应 DTO 与 VO：先返回 service 响应 DTO，再按需转成 VO

很多人容易漏掉一件事：DTO 不只有请求 DTO，也可以有 service 的响应 DTO。

比如 service 处理完下单后，先返回给 controller 一个响应 DTO：

```go
CreateOrderRespDTO {
    OrderID: 90001,
    GoodsAmount: 300,
    CouponAmount: 20,
    FreightAmount: 8,
    PayAmount: 288,
    Status: 1,
}
```

然后 controller 再把它转成给前端看的 VO：

```go
OrderVO {
    OrderID: 90001,
    GoodsAmount: 300,
    CouponAmount: 20,
    FreightAmount: 8,
    PayAmount: 288,
    StatusText: "待支付",
}
```

- 这一段的业务逻辑：
  - service 响应 DTO 更偏“跨层传输”，方便 controller 理解业务处理结果
  - controller / presenter 再根据前端展示需要，把 `Status: 1` 转成 `StatusText: "待支付"`
  - 隐藏内部字段，不直接暴露数据库结构
  - 按前端页面需要输出字段，可能还会做时间、金额格式化
- VO 这一层解决的是：前端最适合看到什么样的数据

#### 总结：你最该记住的不是定义，而是“职责变化”

这几个对象不是为了显得高级，而是为了让每一层只管自己的事：
- DTO 只关心“传进来什么”
  - 前端 → 后端（大而全）
- BO 只关心“业务怎么算、怎么执行”
  - 业务领域对象（按业务拆分）
  - 1 个 BO = 处理一块完整业务
  - 1 个 BO = 对应多张表
- DO 只关心“数据库怎么存”
  - DO：数据库表映射（一张表一个 DO）
- VO 只关心“前端怎么展示”
- DAO 操作库的统一抽象接口
```js
// 举例：
DTO（前端大对象）
↓ 拆分
ShopBO（店铺 BO）
ProductBO（商品 BO）
↓ 再拆分入库
ShopDO（店铺表）
ShopExtDO（店铺扩展表）
ProductDO（商品表）
ProductSkuDO（商品 SKU 表）
ProductDetailDO（商品详情表）
```

- 为什么要分这么多层？（你最关心）
  - 解耦！解耦！解耦！
  - 换数据库 → 主要改 DO
  - 换 ORM → 主要改 DAO
  - 换前端展示 → 主要改 VO
  - 改业务规则 → 主要改 Service / BO
  - 接口参数变 → 主要改 DTO
  - 分层的价值不是“业务逻辑永远不动”，而是“改动范围尽量可控”
## 27周 grpc服务封装更方便的rpc服务

## 28周 深入grpc的服务注册与负载均衡原理
## 29周 基于gin封装api服务

## 30周 可观测的终极解决方案

## 31周 系统监控核心

