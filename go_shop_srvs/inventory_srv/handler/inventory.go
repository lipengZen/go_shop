package handler

import (
	"context"
	"go_shop/go_shop_srvs/inventory_srv/global"
	"go_shop/go_shop_srvs/inventory_srv/model"
	"go_shop/go_shop_srvs/inventory_srv/proto"

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

func (*InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {

	// 扣减库存, 本地事务 [1:10, 2:5, 3:20] 批量扣减需要满足事务的特性,必须都扣减成功/失败
	// 并发场景下,可能会出现超卖
	tx := global.DB.Begin() // 手动事务
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := tx.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
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
