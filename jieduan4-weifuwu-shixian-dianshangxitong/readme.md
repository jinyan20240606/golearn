

# 阶段4-微服务实现电商系统

## 11周 商品微服务的grpc服务

第一章-商品服务-service服务

- 课件代码见商品的grpc服务：`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop_srvs/goods_srv` 目录
  - 目录结构与用户服务保持一致，能复用的就复用
- grpc服务写好接口，需要调试测试，因为web服务还没有开发，只能自己用tests文件运行测试

### 1-1 需求分析-数据库实体分析

需要分析下需求和前端界面，分析下有哪些数据实体需要纳入到商品微服务中来

商品服务这块大概需要有5张表

![alt text](image.png)
1. 轮博图管理实体表
2. 商品分类管理实体表
3. 商品管理信息实体表 -- goods表
4. 品牌表  - log图片和名称
5. 商品分类和品牌关系中间表  -- category_brand_relation表
   1. 品牌也有可能属于多个分类，多对多关系

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



## 13周 库存服务和分布式锁


## 14周 购物车微服务


## 15周 支付宝支付，用户操作微服务


## 16周 用ElasticSearch实现搜索微服务