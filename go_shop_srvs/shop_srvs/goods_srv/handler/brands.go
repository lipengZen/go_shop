package handler

import (
	"context"

	"go_shop/go_shop_srvs/shop_srvs/goods_srv/global"
	"go_shop/go_shop_srvs/shop_srvs/goods_srv/model"
	"go_shop/go_shop_srvs/shop_srvs/goods_srv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (*GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {

	brandListResponse := &proto.BrandListResponse{}

	// brands := make([]*proto.BrandInfoResponse, 0)  注:要用model去查
	var brands []*model.Brands
	result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}
	// zap.S().Info(result)
	// brandListResponse.Data = brands

	var total int64
	global.DB.Model(&model.Brands{}).Count(&total)

	brandListResponse.Total = int32(total)

	for _, brand := range brands {
		brandListResponse.Data = append(brandListResponse.Data, &proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}

	return brandListResponse, nil
}

func (s *GoodsServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {

	if result := global.DB.Where("name=?", req.Name).First(&model.Brands{}); result.RowsAffected != 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌已存在")
	}

	brand := model.Brands{
		Name: req.Name,
		Logo: req.Logo,
	}
	global.DB.Save(brand) // create一样, save既能创建,又能更新

	return &proto.BrandInfoResponse{Id: brand.ID}, nil

}

func (s *GoodsServer) DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {

	if result := global.DB.Delete(&model.Brands{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {

	if result := global.DB.Where("id=?", req.Id).First(&model.Brands{}); result.RowsAffected == 0 { //.Where("name=?", req.Name)
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	brand := model.Brands{}
	if req.Name != "" {
		brand.Name = req.Name
	}
	if req.Logo != "" {
		brand.Logo = req.Logo
	}
	global.DB.Save(brand)

	return &emptypb.Empty{}, nil
}
