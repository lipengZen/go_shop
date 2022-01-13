package main

import (
	"flag"
	"fmt"
	"go_shop/go_shop_srvs/shop_srvs/user_srv/handler"
	"go_shop/go_shop_srvs/shop_srvs/user_srv/proto"
	"net"

	"google.golang.org/grpc"
)

func main() {

	IP := flag.String("ip", "0.0.0.0", "")
	Port := flag.Int("port", 50051, "")

	// flag.Parse()

	fmt.Println("ip: ", *IP)
	fmt.Println("port: ", *Port)

	server := grpc.NewServer()
	// 写错了
	// proto.RegisterUserServer(server, &proto.User{})
	proto.RegisterUserServer(server, &handler.UserServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	err = server.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}

}
