package initialize

import (
	"fmt"
	"go_shop/go_shop_api/user_web/global"
	"go_shop/go_shop_api/user_web/proto"

	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitSrvConn() {

	consulInfo := global.ServerConfig.ConsulInfo

	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name), // tag不能写错
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("用户服务不可达", err)
	}

	global.UserSrvClient = proto.NewUserClient(userConn)

}

func InitSrvConn2() {

	// 从注册中心获取用户服务的信息
	// 服务注册
	cfg := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)

	userSrvHost := ""
	userSrvPort := 0
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service=="%s"`, global.ServerConfig.UserSrvInfo.Name))

	if err != nil {
		panic(err)
	}

	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}
	if userSrvHost == "" {
		// ctx.JSON(http.StatusBadRequest, gin.H{
		// 	"msg": "用户服务不可达",
		// })

		zap.S().Fatal("[InitSrvConn] 连接 [用户服务失败]")
		return
	}

	// addr := fmt.Sprintf("%s:%d", global.ServerConfig.UserServerInfo.Host, global.ServerConfig.UserServerInfo.Port)
	addr := fmt.Sprintf("%s:%d", userSrvHost, userSrvPort)

	// ip := "127.0.0.1"
	// port := 50051
	// fmt.Sprintf("%s:%d", ip, port)
	userConn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())

	if err != nil {
		zap.S().Errorw("[GetUserList] 连接用户服务失败",
			"msg", err.Error)
	}

	// 1. 后续的服务下线了  2.改端口了 3.改ip了  这个在后面负载均衡来做
	// 提前创立好连接,后续不用进行tcp的三次握手,性能更高
	// 问题: 1.一个连接多个goroutine共用,有性能问题 --->需要连接池, 有开源库: grpc-connection-pool, grpc-go-pool； 后面通过负载均衡能解决这个问题
	userSrvCli := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvCli
}
