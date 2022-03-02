package main

import (
	"context"
	"fmt"
	"go_shop/go_shop_srvs/goods_srv/proto"

	"google.golang.org/grpc"
)

var brandClient proto.GoodsClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	brandClient = proto.NewGoodsClient(conn)
}

func TestGetBrandList() {
	rsp, err := brandClient.BrandList(context.Background(), &proto.BrandFilterRequest{
		// Pages:       0,
		// PagePerNums: 3,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, brand := range rsp.Data {
		fmt.Println(brand.Name)
	}
}

/*
  init方法中的 50051 已经是 获取可用端口了
*/

func main() {

	Init()
	TestGetBrandList()

	conn.Close()

}
