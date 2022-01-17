package main

import (
	"fmt"
	"go_shop/go_shop_api/user_web/global"
	"go_shop/go_shop_api/user_web/initialize"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
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

	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动服务器失败: ", err.Error())
	}

}
