package initialize

import (
	"fmt"
	"go_shop/go_shop_api/user_web/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GetEnvInfo(env string) bool {

	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {

	// fmt.Println("debug_env:", GetEnvInfo("GO_SHOP_DEBUG"))

	debug := GetEnvInfo("GO_SHOP_DEBUG")

	configFilePrefix := "config"
	configFileName := fmt.Sprintf("%s-pro.yaml", configFilePrefix)

	if debug {
		configFileName = fmt.Sprintf("%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()

	// 如何设置路径
	v.SetConfigFile(configFileName)

	if err := v.ReadInConfig(); err != nil {
		fmt.Println(err)
		return
	}

	// 这个对象如何在其他文件中使用 - 全局变量,放 global里面
	// serverConfig := global.ServerConfig{}

	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息: ", global.ServerConfig)

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置信息: %s %v", e.Name, global.ServerConfig)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(&global.ServerConfig)
		// fmt.Println("changed: ", serverConfig)
	})
}
