package service

import (
	"context"
	proto2 "mxshop/api/goods/v1"
	proto "mxshop/api/inventory/v1"
	proto3 "mxshop/api/order/v1"
	v12 "mxshop/app/order/srv/internal/data/v1"
	"mxshop/app/order/srv/internal/domain/do"
	"mxshop/app/order/srv/internal/domain/dto"
	"mxshop/app/pkg/code"
	"mxshop/app/pkg/options"
	v1 "mxshop/pkg/common/meta/v1"
	"mxshop/pkg/errors"
	"mxshop/pkg/log"

	"github.com/dtm-labs/client/dtmgrpc"
)

type OrderSrv interface {
	Get(ctx context.Context, orderSn string) (*dto.OrderDTO, error)
	List(ctx context.Context, userID uint64, meta v1.ListMeta, orderby []string) (*dto.OrderDTOList, error)
	Submit(ctx context.Context, order *dto.OrderDTO) error
	Create(ctx context.Context, order *dto.OrderDTO) error
	CreateCom(ctx context.Context, order *dto.OrderDTO) error //这是create的补偿
	Update(ctx context.Context, order *dto.OrderDTO) error
}

type orderService struct {
	data    v12.DataFactory
	dtmOpts *options.DtmOptions
}

func (os *orderService) CreateCom(ctx context.Context, order *dto.OrderDTO) error {
	/*
		1. 删除orderinfo表
		2. 删除ordergoods表
		3. 删除order找到对应的购物车条目，删除购物车条目
	*/
	//其实不用回滚
	//你应该先查询订单是否已经存在，如果已经存在删除相关记录即可， 同时删除购物车记录
	return nil
}

// 真正执行订单创建的
func (os *orderService) Create(ctx context.Context, order *dto.OrderDTO) error {
	/*
		1. 生成orderinfo表
		2. 生成ordergoods表
		3. 根据order找到对应的购物车条目，删除购物车条目
	*/

	var goodsids []int32
	for _, value := range order.OrderGoods {
		goodsids = append(goodsids, value.Goods)
	}
	// return status.Error(codes.Aborted, "订单创建失败") //测试下是否会触发补偿
	// 远程调用商品服务，批量查询商品信息
	goods, err := os.data.Goods().BatchGetGoods(context.Background(), &proto2.BatchGoodsIdInfo{Id: goodsids})
	if err != nil {
		// 日志记录异常：查询失败+商品ID+错误信息
		log.Errorf("批量获取商品信息失败，goodids: %v, err:%v", goodsids, err)
		// 会触发dtm的重试
		return err
	}
	// 校验：传入的商品ID数量 和 查询返回的商品数量不一致 → 部分/全部商品不存在
	if len(goods.Data) != len(goodsids) {
		log.Errorf("批量获取商品信息失败，goodids: %v, 返回值：%v, err:%v", goodsids, goods.Data, err)
		// 抛出自定义业务错误：商品不存在
		return errors.WithCode(code.ErrGoodsNotFound, "商品不存在或者部分不存在")
	}
	// 构建商品ID → 商品信息的Map，方便后续快速取值
	var goodsMap = make(map[int32]*proto2.GoodsInfoResponse)
	for _, value := range goods.Data {
		goodsMap[value.Id] = value
	}
	// 初始化订单总金额
	var orderAmount float32
	// 遍历订单商品，计算总金额、补全商品名称/价格/图片等展示字段
	for _, value := range order.OrderGoods {
		orderAmount += goodsMap[value.Goods].ShopPrice * float32(value.Nums)
		value.GoodsName = goodsMap[value.Goods].Name
		value.GoodsPrice = goodsMap[value.Goods].ShopPrice
		value.GoodsImage = goodsMap[value.Goods].GoodsFrontImage
	}
	// 开启本地数据库事务，保证订单创建、购物车删除原子性
	txn := os.data.Begin()
	defer func() {
		if err := recover(); err != nil {
			_ = txn.Rollback()
			log.Error("新建订单事务进行中出现异常，回滚")
			return
		}
	}()

	err = os.data.Orders().Create(ctx, txn, &order.OrderInfoDO)
	if err != nil {
		// 写入失败，手动回滚事务
		txn.Rollback()
		log.Errorf("创建订单失败，err:%v", err)
		// 重点：直接返回原始error，未返回 gRPC codes.Aborted
		// DTM 判定为临时异常，会**不断重试**当前分支
		return err //这个不是abort 也就是说会不停的重试
	}
	// 事务内：根据用户ID+商品ID，删除对应购物车记录
	err = os.data.ShopCarts().DeleteByGoodsIDs(ctx, txn, uint64(order.User), goodsids)
	if err != nil {
		txn.Rollback()
		log.Errorf("删除购物车失败，goodids:%v, err:%v", goodsids, err)
		// 同样返回原始error，DTM 会重试
		return err
	}
	// 所有操作正常，提交本地事务，数据正式落地
	txn.Commit()
	//这里没有逻辑
	return nil
	// 这块正常业务逻辑来看，其实不需要补偿，不用返回abort错误的
}

func (os *orderService) Get(ctx context.Context, orderSn string) (*dto.OrderDTO, error) {
	order, err := os.data.Orders().Get(ctx, orderSn)
	if err != nil {
		return nil, err
	}
	return &dto.OrderDTO{*order}, nil
}

func (os *orderService) List(ctx context.Context, userID uint64, meta v1.ListMeta, orderby []string) (*dto.OrderDTOList, error) {
	orders, err := os.data.Orders().List(ctx, userID, meta, orderby)
	if err != nil {
		return nil, err
	}
	var ret dto.OrderDTOList
	ret.TotalCount = orders.TotalCount
	for _, value := range orders.Items {
		ret.Items = append(ret.Items, &dto.OrderDTO{
			*value,
		})
	}
	return &ret, nil
}

// 主要用来提交saga事务的
func (os *orderService) Submit(ctx context.Context, order *dto.OrderDTO) error {
	//先从购物车中获取商品信息
	list, err := os.data.ShopCarts().List(ctx, uint64(order.User), true, v1.ListMeta{}, []string{})
	if err != nil {
		log.Errorf("获取购物车信息失败，err:%v", err)
		return err
	}

	if len(list.Items) == 0 {
		log.Errorf("购物车中没有商品，无法下单")
		return errors.WithCode(code.ErrNoGoodsSelect, "没有选择商品")
	}

	var orderGoods []*do.OrderGoods
	var orderItems []*proto3.OrderItemResponse
	for _, value := range list.Items {
		orderGoods = append(orderGoods, &do.OrderGoods{
			Goods: value.Goods,
			Nums:  value.Nums,
		})

		orderItems = append(orderItems, &proto3.OrderItemResponse{
			GoodsId: value.Goods,
			Nums:    value.Nums,
		})
	}
	order.OrderGoods = orderGoods

	//基于可靠消息最终一致性的思想， saga事务来解决订单生成的问题
	var goodsInfo []*proto.GoodsInvInfo
	for _, value := range order.OrderGoods {
		goodsInfo = append(goodsInfo, &proto.GoodsInvInfo{
			GoodsId: value.Goods,
			Num:     value.Nums,
		})
	}
	req := &proto.SellInfo{
		GoodsInfo: goodsInfo,
		OrderSn:   order.OrderSn,
	}
	oReq := &proto3.OrderRequest{
		OrderSn:    order.OrderSn,
		UserId:     order.User,
		Address:    order.Address,
		Name:       order.SignerName,
		Mobile:     order.SingerMobile,
		Post:       order.Post,
		OrderItems: orderItems,
	}

	qsBusi := "discovery:///mxshop-inventory-srv"
	gBusi := "discovery:///mxshop-order-srv"
	saga := dtmgrpc.NewSagaGrpc(os.dtmOpts.GrpcServer, order.OrderSn).
		Add(qsBusi+"/Inventory/Sell", qsBusi+"/Inventory/Reback", req).
		Add(gBusi+"/Order/CreateOrder", gBusi+"/Order/CreateOrderCom", oReq)
	saga.WaitResult = true
	err = saga.Submit()
	//查询最终状态的话：通过OrderSn查询一下， 当前的状态如何状态一直值Submitted那么就你一直不要给前端返回， 如果是failed那么你提示给前端说下单失败，重新下单
	return err
}

func (os *orderService) Update(ctx context.Context, order *dto.OrderDTO) error {
	//TODO implement me
	panic("implement me")
}

func newOrderService(sv *service) *orderService {
	return &orderService{
		data:    sv.data,
		dtmOpts: sv.dtmopts,
	}
}

var _ OrderSrv = &orderService{}
