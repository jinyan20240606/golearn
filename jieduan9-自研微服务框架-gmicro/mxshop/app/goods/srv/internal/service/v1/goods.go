package v1

import (
	"context"
	proto "mxshop/api/goods/v1"
	v1 "mxshop/app/goods/srv/internal/data/v1"
	v12 "mxshop/app/goods/srv/internal/data_search/v1"
	"mxshop/app/goods/srv/internal/domain/do"
	"mxshop/app/goods/srv/internal/domain/dto"
	"sync"

	metav1 "mxshop/pkg/common/meta/v1"
	"mxshop/pkg/log"

	"github.com/zeromicro/go-zero/core/mr"
)

type GoodsSrv interface {
	// 商品列表
	List(ctx context.Context, opts metav1.ListMeta, req *proto.GoodsFilterRequest, orderby []string) (*dto.GoodsDTOList, error)

	// 商品详情
	Get(ctx context.Context, ID uint64) (*dto.GoodsDTO, error)

	// 创建商品
	Create(ctx context.Context, goods *dto.GoodsDTO) error

	// 更新商品
	Update(ctx context.Context, goods *dto.GoodsDTO) error

	// 删除商品
	Delete(ctx context.Context, ID uint64) error

	//批量查询商品
	BatchGet(ctx context.Context, ids []uint64) ([]*dto.GoodsDTO, error)
}

type goodsService struct {
	//工厂
	data v1.DataFactory

	searchData v12.SearchFactory
}

func newGoods(srv *service) *goodsService {
	return &goodsService{
		data:       srv.data,
		searchData: srv.dataSearch,
	}
}

// 遍历树结构
func retrieveIDs(category *do.CategoryDO) []uint64 {
	ids := []uint64{}
	if category == nil || category.ID == 0 {
		return ids
	}
	ids = append(ids, uint64(category.ID))
	for _, child := range category.SubCategory {
		subids := retrieveIDs(child)
		ids = append(ids, subids...)
	}
	return ids
}

func (gs *goodsService) List(ctx context.Context, opts metav1.ListMeta, req *proto.GoodsFilterRequest, orderby []string) (*dto.GoodsDTOList, error) {
	searchReq := v12.GoodsFilterRequest{
		GoodsFilterRequest: req,
	}
	if req.TopCategory > 0 {
		category, err := gs.data.Categorys().Get(ctx, uint64(req.TopCategory))
		if err != nil {
			log.Errorf("categoryData.Get err: %v", err)
			return nil, err
		}

		var ids []interface{}
		for _, value := range retrieveIDs(category) {
			ids = append(ids, value)
		}
		searchReq.CategoryIDs = ids
	}

	goodsList, err := gs.searchData.Goods().Search(ctx, &searchReq)
	if err != nil {
		log.Errorf("serachData.Search err: %v", err)
		return nil, err
	}

	log.Debugf("Search es data: %v", goodsList)

	goodsIDs := []uint64{}
	for _, value := range goodsList.Items {
		goodsIDs = append(goodsIDs, uint64(value.ID))
	}

	//通过id批量查询mysql数据
	goods, err := gs.data.Goods().ListByIDs(ctx, goodsIDs, orderby)
	if err != nil {
		log.Errorf("data.ListByIDs err: %v", err)
		return nil, err
	}
	var ret dto.GoodsDTOList
	ret.TotalCount = int(goodsList.TotalCount)
	for _, value := range goods.Items {
		ret.Items = append(ret.Items, &dto.GoodsDTO{
			GoodsDO: *value,
		})
	}
	return &ret, nil
}

func (gs *goodsService) Get(ctx context.Context, ID uint64) (*dto.GoodsDTO, error) {
	goods, err := gs.data.Goods().Get(ctx, ID)
	if err != nil {
		log.Errorf("data.Get err: %v", err)
		return nil, err
	}
	return &dto.GoodsDTO{
		GoodsDO: *goods,
	}, nil
}

func (gs *goodsService) Create(ctx context.Context, goods *dto.GoodsDTO) error {
	/*
		数据先写mysql，然后写es
	*/
	_, err := gs.data.Brands().Get(ctx, uint64(goods.BrandsID))
	if err != nil {
		return err
	}

	_, err = gs.data.Categorys().Get(ctx, uint64(goods.CategoryID))
	if err != nil {
		return err
	}

	//方案1:之前的入es的方案是给gorm添加aftercreate
	//方案2:原来做的 分布式事务， 异构数据库的事务， 基于可靠消息最终一致性方案 ---- 也是比较重的方案： 每次都要发送一个事务消息
	//方案3:下面这种方案通过获取全局的数据库事务对象来解决数据库和es的同步问题------ 这种方案是不是就没有问题了呢，还是有缺点的，最终一致性比可靠消息稍微弱一点
	//  ---- 因为操作es时，万一es接口1秒超时了，但网速延迟2秒后已经通知到es了，那么此时就不一致了，还是没有RocketMQ的事务消息机制可靠方案强，
	//  ---- 对于一致性不高的场景，可以简单用这个
	txn := gs.data.Begin() // 只要你开启事务，非常小心， 千万不要忘记回滚和commit，
	defer func() {         //防止万一程序此时挂了， 事务没有回滚， 这个地方很重要，一定要紧接者加上这个recover的回滚逻辑
		if err := recover(); err != nil {
			txn.Rollback()
			log.Errorf("goodsService.Create panic: %v", err)
			return
		}
	}()

	err = gs.data.Goods().CreateInTxn(ctx, txn, &goods.GoodsDO)
	if err != nil {
		log.Errorf("data.CreateInTxn err: %v", err)
		txn.Rollback()
		return err
	}
	searchDO := do.GoodsSearchDO{
		ID:          goods.ID,
		CategoryID:  goods.CategoryID,
		BrandsID:    goods.BrandsID,
		OnSale:      goods.OnSale,
		ShipFree:    goods.ShipFree,
		IsNew:       goods.IsNew,
		IsHot:       goods.IsHot,
		Name:        goods.Name,
		ClickNum:    goods.ClickNum,
		SoldNum:     goods.SoldNum,
		FavNum:      goods.FavNum,
		MarketPrice: goods.MarketPrice,
		GoodsBrief:  goods.GoodsBrief,
		ShopPrice:   goods.ShopPrice,
	}

	// 万一es接口1秒超时了，但网速慢2秒后已经通知到es了，那么此时就不一致了，还是没有RocketMQ的事务消息机制可靠方案强，
	err = gs.searchData.Goods().Create(ctx, &searchDO) //这个接口如果超时了
	if err != nil {
		txn.Rollback()
		return err
	}
	txn.Commit()
	return nil
}

func (gs *goodsService) Update(ctx context.Context, goods *dto.GoodsDTO) error {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsService) Delete(ctx context.Context, ID uint64) error {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsService) BatchGet(ctx context.Context, ids []uint64) ([]*dto.GoodsDTO, error) {
	// 如果不想用底层提供的批量接口，可以用这个包进行并发查询，帮我封装了很多细节
	//go-zero提供了一个map-reduce包 非常好用， 但是我们自己去做并发的话 - 一次性启动多个goroutine
	var ret []*dto.GoodsDTO
	var callFuncs []func() error
	var mu sync.Mutex
	for _, value := range ids {
		//tmp := value 必须这样写，// 必须临时变量拷贝
		tmp := value
		callFuncs = append(callFuncs, func() error {
			goodsDTO, err := gs.Get(ctx, tmp)
			mu.Lock() // 保证线程并发安全
			ret = append(ret, goodsDTO)
			mu.Unlock()
			return err
		})
	}
	// 	逐个执行传入的函数切片，每个函数单独起 goroutine
	// 等待所有协程执行完毕
	// 聚合所有错误：只要有一个函数返回 error，mr.Finish 就返回错误
	// 内置协程池、panic 捕获、超时控制（可配置）
	err := mr.Finish(callFuncs...)
	if err != nil {
		return nil, err
	}
	// 使用底层提供的批量接口
	//ds, err := gs.data.ListByIDs(ctx, ids, []string{})
	//if err != nil {
	//	return nil, err
	//}
	//for _, value := range ds.Items {
	//	ret = append(ret, &dto.GoodsDTO{
	//		GoodsDO: *value,
	//	})
	//}
	return ret, nil
}

var _ GoodsSrv = &goodsService{}
