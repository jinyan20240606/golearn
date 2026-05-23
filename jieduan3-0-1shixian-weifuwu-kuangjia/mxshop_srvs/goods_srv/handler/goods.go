package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic/v7"
	"google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/model"
	"mxshop_srvs/goods_srv/proto"
)

type GoodsServer struct {
	// 只用于快速测试接口的连通性，可以用这个自动生成的UnimplementedGoodsServer结构体
	proto.UnimplementedGoodsServer
}

func ModelToResponse(goods model.Goods) proto.GoodsInfoResponse {
	return proto.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brands.ID,
			Name: goods.Brands.Name,
			Logo: goods.Brands.Logo,
		},
	}
}

func (s *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	// 使用es的目的是搜索出商品的id来，通过id拿到具体的字段信息是通过mysql来完成

	//我们使用es是用来做搜索的， 是否应该将所有的mysql字段全部在es中保存一份
	//es用来做搜索，这个时候我们一般只把搜索和过滤的字段信息保存到es中
	//es可以用来当做mysql使用， 但是实际上mysql和es之间是互补的关系， 一般mysql用来做存储使用，es用来做搜索使用
	//es想要提高性能， 就要将es的内存设置的够大， 1k 2k

	//关键词搜索、查询新品、查询热门商品、通过价格区间筛选， 通过商品分类筛选
	goodsListResponse := &proto.GoodsListResponse{}

	//match bool 复合查询
	// 创建 ES 的布尔查询（bool 查询）
	q := elastic.NewBoolQuery()
	localDB := global.DB.Model(model.Goods{})
	if req.KeyWords != "" {
		// 用 GORM 做普通数据库查询
		// 用 ElasticSearch 做关键词全文搜索（搜商品名称、简介）
		// 不使用where("name like ?") 语法， 不使用 MySQL 的 LIKE 模糊查询，改用 ES 全文检索
		// 下面等价于WHERE name LIKE '%手机%' OR goods_brief LIKE '%手机%'
		// 我们对3个字段进行搜索关键词，使用 MultiMatchQuery，多字段分词查询
		q = q.Must(elastic.NewMultiMatchQuery(req.KeyWords, "name", "goods_brief"))
	}
	if req.IsHot {
		localDB = localDB.Where(model.Goods{IsHot: true})
		// 精确字段查询，使用filter不参与算分，不要使用must，业务层面只筛选不需要算分排名
		q = q.Filter(elastic.NewTermQuery("is_hot", req.IsHot))
	}
	if req.IsNew {
		// 都不参与算分，只筛选，精确查询
		q = q.Filter(elastic.NewTermQuery("is_new", req.IsNew))
	}

	if req.PriceMin > 0 {
		// 范围查询
		q = q.Filter(elastic.NewRangeQuery("shop_price").Gte(req.PriceMin))
	}
	if req.PriceMax > 0 {
		q = q.Filter(elastic.NewRangeQuery("shop_price").Lte(req.PriceMax))
	}

	if req.Brand > 0 {
		// 精确查询
		q = q.Filter(elastic.NewTermQuery("brands_id", req.Brand))
	}

	// 通过category去查询商品：电商商品三级分类 + ES terms 查询
	// 功能：根据用户传的 商品分类ID（1级/2级/3级）
	// 自动查询出 所有下级分类ID
	// 然后用 ES 的 terms 查询 匹配这些分类下的所有商品id
	var subQuery string
	categoryIds := make([]interface{}, 0)
	if req.TopCategory > 0 {
		var category model.Category
		if result := global.DB.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}

		if category.Level == 1 {
			// 如果是 1 级分类，要查 所有 3 级子分类
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category WHERE parent_category_id=%d)", req.TopCategory)
		} else if category.Level == 2 {
			// 如果是 2 级分类，直接查 所有下级 3 级分类 ID
			subQuery = fmt.Sprintf("select id from category WHERE parent_category_id=%d", req.TopCategory)
		} else if category.Level == 3 {
			// 如果是 3 级分类，只查自己
			subQuery = fmt.Sprintf("select id from category WHERE id=%d", req.TopCategory)
		}

		type Result struct {
			ID int32 // 可以只写部分字段
		}
		// 执行原生sql
		var results []Result
		// Scan方法：会把数据库查到的值，赋值到指定的结构体中去
		global.DB.Model(model.Category{}).Raw(subQuery).Scan(&results)
		for _, re := range results {
			categoryIds = append(categoryIds, re.ID)
		}

		// es操作：生成terms查询----实现查到分类id属于属于多个categoryId 的多个结果，类似于mysql的in
		// 注意下categoryIds... 的用法：它内部虽然接收的是...interface{} 任意类型，但categoryIds想解构也必须是interface类型，不能是其他类型的解构如：categoryIds := make([]int, 0)
		q = q.Filter(elastic.NewTermsQuery("category_id", categoryIds...))
	}

	//分页
	if req.Pages == 0 {
		req.Pages = 1
	}

	switch {
	case req.PagePerNums > 100:
		req.PagePerNums = 100
	case req.PagePerNums <= 0:
		req.PagePerNums = 10
	}
	// es的分页查询
	result, err := global.EsClient.Search().Index(model.EsGoods{}.GetIndexName()).Query(q).From(int(req.Pages)).Size(int(req.PagePerNums)).Do(context.Background())
	if err != nil {
		return nil, err
	}

	goodsIds := make([]int32, 0)
	goodsListResponse.Total = int32(result.Hits.TotalHits.Value)
	for _, value := range result.Hits.Hits {
		goods := model.EsGoods{}
		_ = json.Unmarshal(value.Source, &goods)
		goodsIds = append(goodsIds, goods.ID)
	}

	//查询id在某个数组中的值
	var goods []model.Goods
	re := localDB.Preload("Category").Preload("Brands").Find(&goods, goodsIds)
	if re.Error != nil {
		return nil, re.Error
	}

	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}

	return goodsListResponse, nil
}

// 现在用户提交订单有多个商品，你得批量查询商品的信息吧
func (s *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	goodsListResponse := &proto.GoodsListResponse{}
	var goods []model.Goods
	//调用where并不会真正执行sql 只是用来生成sql的 当调用find， first才会去执行sql，

	// req.Id = []int32{1,2,3} 是切片类型，where支持对主键的切片和多参where(1,2,3) 这种方式自动解析，非主键字段切片，需要手动写Where("id IN ?", req.Id)
	result := global.DB.Where(req.Id).Find(&goods)
	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}
	goodsListResponse.Total = int32(result.RowsAffected)
	return goodsListResponse, nil
}
func (s *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var goods model.Goods

	if result := global.DB.Preload("Category").Preload("Brands").First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	goodsInfoResponse := ModelToResponse(goods)
	return &goodsInfoResponse, nil
}

func (s *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}
	//先检查redis中是否有这个token
	//防止同一个token的数据重复插入到数据库中，如果redis中没有这个token则放入redis

	//这里没有看到图片文件是如何上传， 在微服务中 普通的文件上传已经不再使用---需要第三方文件上传服务，专门解决文件上传的问题
	goods := model.Goods{
		Brands:          brand,
		BrandsID:        brand.ID,
		Category:        category,
		CategoryID:      category.ID,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.GoodsBrief,
		ShipFree:        req.ShipFree,
		Images:          req.Images,
		DescImages:      req.DescImages, // 这块拿到的其实上传后的文件url参数
		GoodsFrontImage: req.GoodsFrontImage,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		OnSale:          req.OnSale,
	}

	//srv之间互相调用了
	// 事务：包含mysql操作和es操作
	tx := global.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &proto.GoodsInfoResponse{
		Id: goods.ID,
	}, nil
}

// 删除商品----结合es做了优化
// 优化 1：强制指定了 Model 结构体的参数，否则默认空结构体会导致钩子中取不到参数，只拿了零值
// 优化 2: 删除操作，不能用 RowsAffected == 0 来判断 “数据不存在”！必须用 result.Error != nil 来判断是否执行失败！
// 软删除场景下，已经被删过的数据，RowsAffected = 0
// result.Error != nil：只要 SQL 执行成功，不管有没有真的删除数据，都不会返回 error！
func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Goods{BaseModel: model.BaseModel{ID: req.Id}}, req.Id); result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	var goods model.Goods

	if result := global.DB.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	goods.Brands = brand
	goods.BrandsID = brand.ID
	goods.Category = category
	goods.CategoryID = category.ID
	goods.Name = req.Name
	goods.GoodsSn = req.GoodsSn
	goods.MarketPrice = req.MarketPrice
	goods.ShopPrice = req.ShopPrice
	goods.GoodsBrief = req.GoodsBrief
	goods.ShipFree = req.ShipFree
	goods.Images = req.Images
	goods.DescImages = req.DescImages
	goods.GoodsFrontImage = req.GoodsFrontImage
	goods.IsNew = req.IsNew
	goods.IsHot = req.IsHot
	goods.OnSale = req.OnSale

	// 为mysql和es添加事务
	tx := global.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}
