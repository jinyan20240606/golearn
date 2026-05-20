package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"mxshop_srvs/inventory_srv/model"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"

	// redis客户端工具
	goredislib "github.com/go-redis/redis/v8"
	// 分布式锁库（Redis 官方推荐的分布式锁）
	"github.com/go-redsync/redsync/v4"

	// redis连接池
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

	// gorm拼接sql子句的工具包

	"mxshop_srvs/inventory_srv/global"
	"mxshop_srvs/inventory_srv/proto"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

// 设置库存
func (*InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	//设置库存， 如果我要更新库存
	var inv model.Inventory
	// 指定goods这个库存记录更新设置库存
	global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num

	global.DB.Save(&inv)
	return &emptypb.Empty{}, nil
}

// 获取详情
func (*InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有库存信息")
	}
	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

// var m sync.Mutex // 互斥锁

// 扣减库存
func (*InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//没有事务时的漏洞举例子：扣减库存， 本地事务 [1:10,  2:5, 3: 20]，如第一个扣了，第二个库存不足扣减失败，导致异常扣减，数据不一致 --- 只能全部成功或全部失败
	//数据库基本的一个应用场景：就是数据一致性--- 数据库事务，gorm本身是支持事务的
	//并发情况之下 可能会出现超卖 1

	// 使用redis全局锁，解决 原生进程内sync.Mutex锁在分布式下缺点
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "192.168.0.104:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	rs := redsync.New(pool)
	// 开启gorm的手动事务
	tx := global.DB.Begin()
	//m.Lock() //获取锁 这把锁有问题吗？  假设有10w的并发， 这里并不是请求的同一件商品  这个锁就没有问题了吗？

	//这个时候应该先查询表，然后确定这个订单是否已经扣减过库存了，已经扣减过了就别扣减了
	//并发时候会有漏洞， 同一个时刻发送了重复了多次， 使用锁，分布式锁
	sellDetail := model.StockSellDetail{
		OrderSn: req.OrderSn,
		Status:  1,
	}
	var details []model.GoodsDetail
	// for循环拿到每个商品的库存信息
	for _, goodInfo := range req.GoodsInfo {
		details = append(details, model.GoodsDetail{
			Goods: goodInfo.GoodsId,
			Num:   goodInfo.Num,
		})

		var inv model.Inventory
		// 使用 gorm 的 for update 写法实现 MySQL 悲观锁
		// 在单库单表库存扣减场景下，通常可以直接依赖“事务 + 行级悲观锁”解决并发超卖问题，
		// 一般不需要再额外加全局互斥锁或 Redis 分布式锁。
		// 但前提是：查询、判断、更新、提交必须全部放在同一个事务里完成。
		// if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
		// 	tx.Rollback() //回滚之前的操作
		// 	return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		// }

		// for {  // 无限循环：专门用于乐观锁失败重试的
		// 创建商品粒度的锁
		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

		// 1、先查有没有商品对应的库存信息
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		// 2、判断库存是否充足
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		// 3、开始扣减，减个数量，， 会出现数据不一致的问题 - 锁，分布式锁
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)

		// 释放redis分布式锁
		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}
		//原生sql语句：update inventory set stocks = stocks-1, version=version+1 where goods=goods and version=version
		//这种写法有瑕疵，为什么？ --- 零值问题 Stocks: inv.Stocks这个字段对于int类型来说 默认值是0 这种会被gorm给忽略掉，强制更新零值，必须使用Select语法
		// version=？这里面的问号也是占位符，后面由变量填充
		//if result := tx.Model(&model.Inventory{}).Select("Stocks", "Version").Where("goods = ? and version= ?", goodInfo.GoodsId, inv.Version).Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version+1}); result.RowsAffected == 0 {
		//	zap.S().Info("库存扣减失败")
		//  扣减失败时，退出这个无限循环
		//}else{
		//	break
		//}
		// }
		//tx.Save(&inv)
	}
	sellDetail.Detail = details
	//写selldetail表
	if result := tx.Create(&sellDetail); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "保存库存扣减历史失败")
	}
	tx.Commit() // 需要自己手动提交操作
	//m.Unlock() //释放锁
	return &emptypb.Empty{}, nil
}

// 订单归还
func (*InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//库存归还： 1：订单超时归还 2. 订单创建失败，归还之前扣减的库存 3. 手动归还
	// 面临的潜在问题：本地事务 和分布式事务  和并发时的 分布式锁的问题
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}

		//扣减， 会出现数据不一致的问题 - 锁，分布式锁
		// 只要并发场景下面临着同时改一个值，就意味有改错的可能，
		inv.Stocks += goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit() // 需要自己手动提交操作
	return &emptypb.Empty{}, nil
}

func (*InventoryServer) TrySell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//扣减库存， 本地事务 [1:10,  2:5, 3: 20]
	//数据库基本的一个应用场景：数据库事务
	//并发情况之下 可能会出现超卖 1
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "192.168.0.104:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	rs := redsync.New(pool)

	tx := global.DB.Begin()
	//m.Lock() //获取锁 这把锁有问题吗？  假设有10w的并发， 这里并不是请求的同一件商品  这个锁就没有问题了吗？
	for _, goodInfo := range req.GoodsInfo {
		var inv model.InventoryNew
		//if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods:goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
		//	tx.Rollback() //回滚之前的操作
		//	return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		//}

		//for {
		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		//判断库存是否充足
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		//扣减， 会出现数据不一致的问题 - 锁，分布式锁
		//inv.Stocks -= goodInfo.Num
		inv.Freeze += goodInfo.Num
		tx.Save(&inv)

		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}
		//update inventory set stocks = stocks-1, version=version+1 where goods=goods and version=version
		//这种写法有瑕疵，为什么？
		//零值 对于int类型来说 默认值是0 这种会被gorm给忽略掉
		//if result := tx.Model(&model.Inventory{}).Select("Stocks", "Version").Where("goods = ? and version= ?", goodInfo.GoodsId, inv.Version).Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version+1}); result.RowsAffected == 0 {
		//	zap.S().Info("库存扣减失败")
		//}else{
		//	break
		//}
		//}
		//tx.Save(&inv)
	}
	tx.Commit() // 需要自己手动提交操作
	//m.Unlock() //释放锁
	return &emptypb.Empty{}, nil
}

func (*InventoryServer) ConfirmSell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//扣减库存， 本地事务 [1:10,  2:5, 3: 20]
	//数据库基本的一个应用场景：数据库事务
	//并发情况之下 可能会出现超卖 1
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "192.168.0.104:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	rs := redsync.New(pool)

	tx := global.DB.Begin()
	//m.Lock() //获取锁 这把锁有问题吗？  假设有10w的并发， 这里并不是请求的同一件商品  这个锁就没有问题了吗？
	for _, goodInfo := range req.GoodsInfo {
		var inv model.InventoryNew
		//if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods:goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
		//	tx.Rollback() //回滚之前的操作
		//	return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		//}

		//for {
		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		//判断库存是否充足
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		//扣减， 会出现数据不一致的问题 - 锁，分布式锁
		inv.Stocks -= goodInfo.Num
		inv.Freeze -= goodInfo.Num
		tx.Save(&inv)

		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}
		//update inventory set stocks = stocks-1, version=version+1 where goods=goods and version=version
		//这种写法有瑕疵，为什么？
		//零值 对于int类型来说 默认值是0 这种会被gorm给忽略掉
		//if result := tx.Model(&model.Inventory{}).Select("Stocks", "Version").Where("goods = ? and version= ?", goodInfo.GoodsId, inv.Version).Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version+1}); result.RowsAffected == 0 {
		//	zap.S().Info("库存扣减失败")
		//}else{
		//	break
		//}
		//}
		//tx.Save(&inv)
	}
	tx.Commit() // 需要自己手动提交操作
	//m.Unlock() //释放锁
	return &emptypb.Empty{}, nil
}

func (*InventoryServer) CancelSell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//扣减库存， 本地事务 [1:10,  2:5, 3: 20]
	//数据库基本的一个应用场景：数据库事务
	//并发情况之下 可能会出现超卖 1
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "192.168.0.104:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	rs := redsync.New(pool)

	tx := global.DB.Begin()
	//m.Lock() //获取锁 这把锁有问题吗？  假设有10w的并发， 这里并不是请求的同一件商品  这个锁就没有问题了吗？
	for _, goodInfo := range req.GoodsInfo {
		var inv model.InventoryNew
		//if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods:goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
		//	tx.Rollback() //回滚之前的操作
		//	return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		//}

		//for {
		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		//判断库存是否充足
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		//扣减， 会出现数据不一致的问题 - 锁，分布式锁
		inv.Freeze -= goodInfo.Num
		tx.Save(&inv)

		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}
		//update inventory set stocks = stocks-1, version=version+1 where goods=goods and version=version
		//这种写法有瑕疵，为什么？
		//零值 对于int类型来说 默认值是0 这种会被gorm给忽略掉
		//if result := tx.Model(&model.Inventory{}).Select("Stocks", "Version").Where("goods = ? and version= ?", goodInfo.GoodsId, inv.Version).Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version+1}); result.RowsAffected == 0 {
		//	zap.S().Info("库存扣减失败")
		//}else{
		//	break
		//}
		//}
		//tx.Save(&inv)
	}
	tx.Commit() // 需要自己手动提交操作
	//m.Unlock() //释放锁
	return &emptypb.Empty{}, nil
}

func AutoReback(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	type OrderInfo struct {
		OrderSn string
	}
	for i := range msgs {
		//既然是归还库存，那么我应该具体的知道每件商品应该归还多少， 但是有一个问题是什么？重复归还的问题
		//所以说这个接口应该确保幂等性， 你不能因为消息的重复发送导致一个订单的库存归还多次， 没有扣减的库存你别归还
		//如果确保这些都没有问题， 新建一张表， 这张表记录了详细的订单扣减细节，以及归还细节
		var orderInfo OrderInfo
		err := json.Unmarshal(msgs[i].Body, &orderInfo)
		if err != nil {
			zap.S().Errorf("解析json失败： %v\n", msgs[i].Body)
			return consumer.ConsumeSuccess, nil
		}

		//去将inv的库存加回去 将selldetail的status设置为2， 要在事务中进行
		tx := global.DB.Begin()
		var sellDetail model.StockSellDetail
		if result := tx.Model(&model.StockSellDetail{}).Where(&model.StockSellDetail{OrderSn: orderInfo.OrderSn, Status: 1}).First(&sellDetail); result.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil
		}
		//如果查询到那么逐个归还库存
		for _, orderGood := range sellDetail.Detail {
			//update怎么用
			//先查询一下inventory表在， update语句的 update xx set stocks=stocks+2
			if result := tx.Model(&model.Inventory{}).Where(&model.Inventory{Goods: orderGood.Goods}).Update("stocks", gorm.Expr("stocks+?", orderGood.Num)); result.RowsAffected == 0 {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
		}

		if result := tx.Model(&model.StockSellDetail{}).Where(&model.StockSellDetail{OrderSn: orderInfo.OrderSn}).Update("status", 2); result.RowsAffected == 0 {
			tx.Rollback()
			return consumer.ConsumeRetryLater, nil
		}
		tx.Commit()
		return consumer.ConsumeSuccess, nil
	}
	return consumer.ConsumeSuccess, nil
}
