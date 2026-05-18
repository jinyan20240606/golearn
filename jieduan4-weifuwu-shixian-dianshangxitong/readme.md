

# 阶段4-微服务实现电商系统

## 11周 商品微服务的grpc服务

第一章-商品服务-service服务

- 课件代码见商品的grpc服务：`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop_srvs/goods_srv` 目录
  - 目录结构与用户服务保持一致，能复用的就复用
- grpc服务写好接口，需要调试测试，因为web服务还没有开发，只能自己用tests文件运行测试

### 1-1 需求分析-数据库实体分析

需要分析下需求和前端界面，分析下有哪些数据实体需要纳入到商品微服务中来

通过需求分析商品服务这块大概需要有5张表

![alt text](image.png)
1. 轮博图管理实体表
2. 商品分类管理实体表
3. 商品管理信息实体表 -- goods表
4. 品牌表  - log图片和名称
5. 商品分类和品牌关系中间表  -- category_brand_relation表
   1. 品牌也有可能属于多个分类，多对多关系

#### 表结构对应的数据库表数据示例

> 根据jieduan3-0-1shixian-weifuwu-kuangjia/mxshop_srvs/goods_srv/model表文件实际对应的数据库示例
```go
// 1、categories 商品分类表（自关联三级分类）
id	created_at	updated_at	deleted_at	name	parent_category_id	level	is_tab
1	2025-01-01	2025-01-01	NULL	手机	0	1	1
2	2025-01-01	2025-01-01	NULL	安卓	1	2	0
3	2025-01-01	2025-01-01	NULL	小米	2	3	0

// 2、brands 商品品牌表
id	name	logo
1	小米	https://xxx.png
2	华为	https://xxx.png

// 3、goodscategorybrand 分类 <-> 品牌 中间表（多对多）一个分类可以绑定多个品牌，一个品牌可以属于多个分类
id	category_id	brands_id
1	1	1
2	1	2
3	2	1

// 4、goods 商品表
id	name	category_id	brands_id	shop_price	goods_brief	on_sale	ship_free
1	小米 14	3	1	3999	最新小米旗舰	1	1
2	华为 Mate70	2	2	4999	华为旗舰	1	1

// 5、banners 轮博图表
id	created_at	updated_at	deleted_at	image	url	index
1	2026-05-18 10:00:00	2026-05-18 10:00:00	NULL	https://img1.jpg	/goods/1001	1
2	2026-05-18 10:05:00	2026-05-18 10:05:00	NULL	https://img2.jpg	/category/1	2
3	2026-05-18 10:10:00	2026-05-18 10:10:00	NULL	https://img3.jpg	/activity/sale	3
```

### 1-2 需求分析- 商品微服务接口分析

分析下商品相关的需求和前端界面，都有哪些功能，需要提供哪些接口

1. 商品列表页
   1. 通过分类筛选
   2. 通过价格区间筛选
   3. 通过搜索框搜索
   4. 新品筛选
   5. 热卖商品的筛选
   6. 通过品牌筛选
2. 商品详情
3. 批量查询商品的接口
4. 后台管理系统的功能
   1. 添加商品
   2. 修改商品
   3. 删除商品
5. 品牌列表
   1. 后台系统中的增删改查
6. 商品的分类接口
   1. 通过1级分类查询出所有的二级分类
   2. 后台系统的分类增删改查
7. 品牌的分类
   1. 列表 
   2. 后台系统的分类增删改查
   3. 通过分类查找品牌

- 接口注意事项
  - 一旦设计修改数据库都需要必须加上权限校验，查询类的可选鉴权

### 1-3 商品多级分类表结构设计及细节注意
- 设计多级分类表：`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop_srvs/goods_srv/model/goods.go`的Category表结构体定义
  - 重点讲解多级分类表的设计细节和新手避坑的地方
- 重点注意
  - gorm的自关联语法写法
  - GORM 的 AutoMigrate 默认会创建真实数据库外键约束。foreignKey / references 等 tag 是告诉 GORM 关联逻辑的，但并不阻止外键创建。如果不想创建真实外键，需要设置 DisableForeignKeyConstraintWhenMigrating: true。constraint tag 的作用是自定义外键的级联规则，而非控制是否创建外键。
    - foreignKey / references tag 是给 GORM 查询引擎看的；DisableForeignKeyConstraintWhenMigrating 是给 GORM 迁移引擎看的。两者完全独立，互不影响。
  - 但是请注意：**生产项目几乎都禁用 GORM 自动创建外键约束**
    -  性能问题：真实外键约束在每次 INSERT/UPDATE/DELETE 时数据库都要做约束校验，高并发场景下是性能瓶颈。


### 1-4 品牌，轮博图表结构设计

- `jieduan3-0-1shixian-weifuwu-kuangjia/mxshop_srvs/goods_srv/model/goods.go`的brand表和GoodsCategoryBrand表结构体定义
  - GoodsCategoryBrand：多对多关系是另起一张表
  - 注意细节
    - 联合唯一索引的创建：限制数据库中：(CategoryID + BrandsID) 这一组数据，绝对不能重复出现！
      - 防止同一个分类下，重复绑定同一个品牌！比如：手机分类 → 小米品牌，不能再绑定一次手机分类 → 小米品牌，这就是联合唯一索引的作用。
    - 多对多关联的写法
- `jieduan3-0-1shixian-weifuwu-kuangjia/mxshop_srvs/goods_srv/model/goods.go`的Banner表结构体定义

### 1-5 商品表结构的设计


- `jieduan3-0-1shixian-weifuwu-kuangjia/mxshop_srvs/goods_srv/model/goods.go`的Goods表结构体定义
- 细节注意
  - 在gorm表结构体中如何定义数组类型：如用在 Images字段：商品轮播图列表（JSON数组）
    - gorm默认没有数组类型，需要自定义一个gorm的数组类型，见`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop_srvs/goods_srv/model/base.go`的GormList类型自定义
    - 这里用自定义类型GormList---在数据库中存的是字符串类型，或者用官方提供的 datatypes.JSON，支持数组对象，只不过它在数据库中存的是json类型

### 1-6 生成表结构和导入数据

- 见 `jieduan3-0-1shixian-weifuwu-kuangjia/mxshop_srvs/goods_srv/model/main/main.go` 来生成数据库创建新表


### 1-7 定义proto接口

- 建 `jieduan3-0-1shixian-weifuwu-kuangjia/mxshop_srvs/goods_srv/proto/goods.proto`

### 1-8 快速启动grpc服务

快速启动学会用`proto.UnimplementedGoodsServer`方法，进行初期快速连通接口测试

- 先完善`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop_srvs/goods_srv/main.go` 中的启动proto的grpc服务实例
  - `proto.RegisterGoodsServer(server, &handler.GoodsServer{})`
- 接着定义handler层，实现proto接口：`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop_srvs/goods_srv/handler/goods.go`
  - 细节注意：
  - 开发初期：如果你想快速测试接口的连通性，可以用这个自动生成的UnimplementedGoodsServer结构体：使用`proto.UnimplementedGoodsServer` proto自动生成 --- 只是做初期测试用
  - 后面开发：就需要自己实现完整具体的GoodsServer结构体方法
- 添加对应的nacos配置，就可以快速启动了

### 1-9 品牌列表的实现

- 本节实现轮博图和品牌接口的完成
- 先实现`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop_srvs/goods_srv/tests/brands.go` 测试文件进行接口开发测试
  - 由于目前是用UnimplementedGoodsServer快速实现的，接口发现，测试发现会它默认会帮我们自动返回错误的状态码，未实现的错误信息
- 品牌列表细节注意
  - 学习到一个返回品牌分页时，返回分页的数据，总数，总数的获取学习使用gorm的`Count()`方法
    - 文件见`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop_srvs/goods_srv/handler/brands.go`

### 1-10 品牌的其他接口：新建删除更新

- 文件见`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop_srvs/goods_srv/handler/brands.go`

### 1-11 轮播图的增删改查crud

- 文件见`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop_srvs/goods_srv/handler/banner.go`



### 1-12 商品分类的列表接口1

这个需要重点看下关联表 和多级分类的实现

- 见`mxshop_srvs/goods_srv/handler/category.go`
  - 重点实现返回适合前端展示的分类大数组结构，包含子分类（一级分类二级分类）
- 分类列表接口`GetAllCategorysList`：希望返回如下这种嵌套的大json结构给前端，拼装好这种结构我们选择在goods_srv层做，不交给web-服务层做，
  - 因为这种结构用gorm来做是非常简单的，web-服务一般不与数据库交互，自行拼接是很麻烦的
    - goods_srv层做我们返回Data原始结构和JsonData拼好的大json结构
- 先涉及完善分类表结构的自关联语法，子分类切片结构体关联父分类结构体， -- 详见`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop_srvs/goods_srv/model/goods.go`的Category表结构体,自关联外键相关写法
- 完善后，在我们的接口处理函数`GetAllCategorysList`中使用这种表结构
  - 使用Preload预加载语法
- 最后在`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop_srvs/goods_srv/tests/category/category.go`的`TestGetCategoryList`测试文件中测试接口响应

### 1-13 商品分类的列表接口2

- 前面遇到问题：调用分类的接口，接口响应中只能往下加载一级，不能无限递归加载，
  - 解决：正确用Preload预加载语法，使用点号语法`Preload("SubCategory.SubCategory")`
    - 语法：点的数量 = 子分类层数，
    - 如果有 四级分类就写2个点3段：`Preload("SubCategory.SubCategory.SubCategory")`

### 1-14 获取商品分类的子分类

- 实现`mxshop_srvs/goods_srv/handler/category.go`的GetSubCategory接口方法
  - 获取一个分类下的子分类
- 然后在test下测试

### 1-15 商品分类的新建删除和更新

- 实现相关剩余的分类接口

### 1-16 品牌分类相关接口

- 略

### 1-17&18&19 商品列表页接口
- 实现GoodsList接口
  - 涉及多维度筛选，复合查询，分类查询等实现的难点
  - 这个比较重要稍微复杂一些
- 然后tests文件下测试该接口


### 1-20 批量获取商品信息、商品详情接口

- 注意下批量查询时的where语法-见BatchGetGoods方法

### 1-21 新增修改和删除商品接口

- CreateGoods接口方法等
- 删除方法
  - 使用的是逻辑删除

## 12周 商品微服务的gin层和oss图片上传

### 1章 gin完成商品服务的http接口
> 主要完成 `jieduan3-0-1shixian-weifuwu-kuangjia/mxshop-api/goods-web`目录，完成与goods_srv的联调

#### 1-1 快速将用户的web服务复刻出商品的web服务
**重点注意**：
1. 迁移中间件时涉及鉴权功能，这个可能每个服务可能都需要统一的这种代码逻辑 ---- 涉及公共代码的优缺点
##### 公共代码的优缺点
   1. 优点：抽离出来，各个微服务公用，是有好处的，抽离，减少代码量 
   2. 缺点：公共的代码作为多个微服务公用的话，有一个很大的缺点，一旦修改可能会影响其他服务，而且不知道当时有哪些服务在用，不知道影响面，一旦有问题，其他微服务都会受到影响
   3. 解决：公共的代码给大家统一使用，必须要加版本管理，方便风险控制

### 1-2&3 商品列表页接口1&2

- `jieduan3-0-1shixian-weifuwu-kuangjia/mxshop-api/goods-web/api/goods/goods.go` 的List方法

web服务层的接口都是跟业务强相关的，

### 1-4 如何设计一个符合go风格的注册中心接口？

- 我们需要把我们的商品web服务注册到注册中心去，如consul去，但是我们必须统一封装一个适配器，解耦和扩展性，支持扩展后面能随时切换注册中心平台
- 我们建立`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop-api/goods-web/utils/register` 单独的注册中心模块，可以方便扩展
- 这节课重点：不是把商品gin层项目接入注册中心，而是设计一个通用扩展性高的集成注册中心功能
- 然后在`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop-api/goods-web/main.go` 引入（集成注册中心模块）使用它注册到consule中

### 1-5 gin的退出后的服务注销

- 见`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop-api/goods-web/main.go`
- 1. 优雅的退出功能，一定要先把这个Router服务启动放到协程里
- 2. 在系统退出信号中--然后调用前面封装的统一注册中心接口暴露的注销方法register_client.DeRegister
### 1-6 用户的web服务注册和优雅退出

- 给user-web服务也封装统一注册集成接口，添加注册功能和退出注销功能

- 见`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop-api/user-web/main.go`


### 1-7 新建商品

- 完成`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop-api/goods-web/router/goods.go` 的 `goods.New` 相关路由的实现
  - 以及`api/goods/goods.go 的 func New(ctx *gin.Context) 方法`
  - 重点关注
    - 但凡这种post，put改数据库类的请求，在微服务中都是比较重要关注的，稍复杂，因为设计跨微服务的数据库交互，必须要考虑的如最重要的分布式事务问题。。。
    - 如何设置库存 --- 这块是重点后面单独讲，TODO 商品的库存 - 分布式事务

### 1-8 获取商品详情

- 完成`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop-api/goods-web/router/goods.go`的`GoodsRouter.GET("/:id", goods.Detail) //获取商品的详情`handler逻辑
  - handelr中的Detail方法
- 重点细节
  - 只有需要在 gRPC 服务端里拿到 Gin 上下文信息时，才用 WithValue
    - 如 grpc服务端中想获取请求头 token，用户id什么的
      - ginCtx := ctx.Value("ginContext").(*gin.Context)
      - userId := ginCtx.GetInt("userId")
    - 那么webgin层，调用时一定要
      - rsp, err := goodsClient.GetGoodsDetail(
        - context.WithValue(context.Background(), "ginContext", ctx), // 把 gin ctx 塞进去
        - &proto.GoodInfoRequest{Id: int32(i)},
      - )
  - 普通接口、不需要传递信息时，直接用grpc的上下文即可： context.Background ()*

### 1-9 商品删除更新

- 完成路由处理方法：`GoodsRouter.DELETE("/:id",middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Delete) //删除商品`
- 再增加一个获取商品库存的接口`GoodsRouter.GET("/:id/stocks", goods.Stocks)   `
  - 目的：此时主要是在gin层留出一个口子，后续在goods_srv中实现苦寻接口
- 完善一个更新路由接口`GoodsRouter.PATCH("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.UpdateStatus)`
  - 更新部分状态，接口-- 部分更新用PATCH方法
  - 使用抽离form表单类型：`forms/goods.go`
- 完善`GoodsRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)         // 全量更新用PUT方法`
  - 使用抽离的form表单类型：`forms/goods.go`

### 1-10 商品的分类接口

- 分类接口单独拆分文件模块，与商品本身的增删改查不能混在一起

- 新建分类路由接口`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop-api/goods-web/router/category.go`
- **重点注意的**
  - 新建分类时需要新建个form表单类型`forms/category.go`，from里的bool类型必须用指针类型，原因如下：
    - bool 零值 = false，即使前端根本没传这个字段，Gin 也会认为它的值是 false，用户不传 is_tab → 系统以为是 false → 验证通过，验证器没有拦住
    - 用 *bool 指针就完美解决 IsTab *bool `binding:"required"`
      - 因为指针的零值 = nil（空）
  - api下goods和category文件存在公共帮助代码重复使用，可以抽离到api/base.go文件自定义api包名下，直接api包名引用

### 1-11 轮博图接口和yapi的快速测试

1. `jieduan3-0-1shixian-weifuwu-kuangjia/mxshop-api/goods-web/router/banner.go`
2. 建个banner的form类型：`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop-api/goods-web/forms/banner.go`
3. 新建`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop-api/goods-web/api/banners/banner.go`

### 1-12 品牌列表页接口

- 完成`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop-api/goods-web/router/brand.go`
  - 品牌列表的CRUD接口
### 1-13 品牌分类CRUD接口

- 接着完成`goods-web/router/brand.go`品牌分类接口

**例如**
分类1：手机 → 绑定品牌：华为、小米
分类2：电脑 → 绑定品牌：联想、苹果
1 GetCategoryBrandList(ctx *gin.Context)。---- /rpc/category/1/brands接收分类id
[
  { "id": 1, "name": "华为", "logo": "xxx" },
  { "id": 2, "name": "小米", "logo": "xxx" }
]

2 CategoryBrandList(ctx *gin.Context) --- /rpc/category/brands// 调用时不用传参
{
  "total": 4,
  "data": [
    { "category": {id:1,name:"手机"}, "brand": {id:1,name:"华为"} },
    { "category": {id:1,name:"手机"}, "brand": {id:2,name:"小米"} },
    { "category": {id:2,name:"电脑"}, "brand": {id:3,name:"联想"} },
    { "category": {id:2,name:"电脑"}, "brand": {id:4,name:"苹果"} }
  ]
}

### 2章 阿里云的oss服务集成


## 13周 库存服务和分布式锁


## 14周 购物车微服务


## 15周 支付宝支付，用户操作微服务


## 16周 用ElasticSearch实现搜索微服务