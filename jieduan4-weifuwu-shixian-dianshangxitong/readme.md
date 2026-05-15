

# 阶段4-微服务实现电商系统

## 11周 商品微服务的grpc服务

第一章-商品服务-service服务

- 课件代码见商品的grpc服务：`jieduan3-0-1shixian-weifuwu-kuangjia/mxshop_srvs/goods_srv` 目录
  - 目录结构与用户服务保持一致，能复用的就复用

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


### 1-7 

## 12周 商品微服务的gin层和oss图片上传



## 13周 库存服务和分布式锁


## 14周 购物车微服务


## 15周 支付宝支付，用户操作微服务


## 16周 用ElasticSearch实现搜索微服务