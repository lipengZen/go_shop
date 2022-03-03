package handler

import (
	"context"
	"fmt"
	"go_shop/go_shop_srvs/inventory_srv/global"
	"go_shop/go_shop_srvs/inventory_srv/model"
	"go_shop/go_shop_srvs/inventory_srv/proto"
	"sync"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

func (*InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {

	//设置库存， 如果我要更新库存
	var inv model.Inventory
	global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num

	global.DB.Save(&inv)
	return &emptypb.Empty{}, nil
}

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

var m sync.Mutex

func (*InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {

	// 扣减库存, 本地事务 [1:10, 2:5, 3:20] 批量扣减需要满足事务的特性,必须都扣减成功/失败
	// 并发场景下,可能会出现超卖
	tx := global.DB.Begin() // 手动事务
	// 单实例 锁
	// m.Lock()
	// defer m.Unlock()

	
	pool := goredis.NewPool(global.RedisClient) // or, pool := redigo.NewPool(...)
	rs := redsync.New(pool)

	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		// 悲观锁
		// if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
		// 	tx.Rollback() // 回滚之前的操作
		// 	return nil, status.Errorf(codes.NotFound, "没有库存信息")
		// }

		// 乐观锁
		// for {
		// 乐观锁

		// redis 分布式锁
		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))

		if err := mutex.Lock(); err != nil {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}
		// 在commit前, 或者返回之前 , 释放锁
		defer func() {
			if ok, err := mutex.Unlock(); !ok || err != nil {
				tx.Rollback() // 回滚之前的操作
				zap.S().Warn(status.Errorf(codes.Internal, "释放redis分布式锁异常"))
			}
		}()

		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.NotFound, "没有库存信息")
		}
		// 判断库存是否充足
		if inv.Stocks < goodInfo.Num {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		// 扣减, 出现数据不一致的问题  - 锁,分布式锁
		inv.Stocks -= goodInfo.Num

		// 乐观锁
		//update inventory set stocks = stocks-1, version=version+1
		// where goods=goods and version=version
		// model里面不能使用 inv, 否则inv里的值会成为 where后面的查询语句
		// 零值会被 gorm 给忽略掉, 导致更新为0被忽略 -> 更新选定字段
		// if result := tx.Model(&model.Inventory{}).Select("Stocks", "Version").
		// 						Where("goods = ? and version= ?", goodInfo.GoodsId, inv.Version).
		// 						Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version + 1}); result.RowsAffected == 0 {
		// 	zap.S().Info("库存扣减失败")
		// } else {
		// 	break
		// }
		// }

		tx.Save(&inv)

	}
	tx.Commit() // 需要手动提交

	return &emptypb.Empty{}, nil
}

func (*InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//库存归还： 1：订单超时归还 2. 订单创建失败，归还之前扣减的库存(分布式事务) 3. 手动归还

	tx := global.DB.Begin() // 手动事务
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := tx.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.NotFound, "没有库存信息")
		}

		// 扣减, 出现数据不一致的问题  - 锁,分布式锁
		inv.Stocks += goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit() // 需要手动提交

	return &emptypb.Empty{}, nil

}
