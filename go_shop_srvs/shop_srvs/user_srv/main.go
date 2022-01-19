package main

import (
	"flag"
	"fmt"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"go_shop/go_shop_srvs/shop_srvs/user_srv/global"
	"go_shop/go_shop_srvs/shop_srvs/user_srv/handler"
	"go_shop/go_shop_srvs/shop_srvs/user_srv/initialize"
	"go_shop/go_shop_srvs/shop_srvs/user_srv/proto"
	"go_shop/go_shop_srvs/shop_srvs/user_srv/utils"

	"net"
)

func main() {

	IP := flag.String("ip", "0.0.0.0", "")
	Port := flag.Int("port", 0, "") //50051

	flag.Parse()

	// fmt.Println("ip: ", *IP)
	// fmt.Println("port: ", *Port)
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	zap.S().Info("ip: ", *IP)

	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	zap.S().Info("port: ", *Port)

	server := grpc.NewServer()
	// 写错了
	// proto.RegisterUserServer(server, &proto.User{})
	proto.RegisterUserServer(server, &handler.UserServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	// 注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("172.18.0.1:%d", *Port), // 50051", // 之后改到配置中心拿
		Timeout:                        "5s",
		Interval:                       "1s",
		DeregisterCriticalServiceAfter: "100s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	registration.ID = global.ServerConfig.Name
	registration.Port = *Port
	registration.Tags = []string{"user_srv", "lee", ""}
	registration.Address = "172.18.0.1"
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	// client.Agent().ServiceDeregister(id)
	if err != nil {
		panic(err)
	}

	err = server.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}

}
