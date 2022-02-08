package main

import (
	"fmt"
	"go_shop/go_shop_api/user_web/global"
	"go_shop/go_shop_api/user_web/initialize"
	"go_shop/go_shop_api/user_web/utils"
	"go_shop/go_shop_api/user_web/utils/consul"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	myvalidator "go_shop/go_shop_api/user_web/validator"
)

func main() {

	// 1.初始化logger
	initialize.InitLogger()

	// 2. 初始化配置文件
	initialize.InitConfig()

	// 3.初始化 routers
	Router := initialize.Routers()

	// 4.初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		zap.S().Errorf("初始化翻译 err:", err)
		return
	}

	// 5.初始化srv的连接
	initialize.InitSrvConn()

	viper.AutomaticEnv()
	//如果是本地开发环境端口号固定，线上环境启动获取端口号
	debug := viper.GetBool("GO_SHOP_DEBUG")
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	// port := 9023
	/*
		1. S()可以获取一个全局的sugar，可以让我们自己设置一个全局的logger
		2. 日志是分级别的，debug， info， warn， error， fatal
		3. s函数和l函数，已经加了锁
	*/

	// 可以省略代码，同样拿到一个全局的sugger
	zap.S().Debugf("启动服务器， 端口： %d", global.ServerConfig.Port)

	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())

	err := register_client.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败:", err.Error())
	}

	// if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
	// 	zap.S().Panic("启动服务器失败: ", err.Error())
	// }

	go func() {
		if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动失败:", err.Error())
		}
	}()
	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = register_client.DeRegister(serviceId); err != nil {
		zap.S().Info("注销失败:", err.Error())
	} else {
		zap.S().Info("注销成功:")
	}

}
