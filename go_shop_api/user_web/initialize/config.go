package initialize

import (
	"encoding/json"
	"fmt"
	"go_shop/go_shop_api/user_web/global"

	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
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
		zap.S().Fatal("读取 viper 配置信息: ", err)
		return
	}

	// 这个对象如何在其他文件中使用 - 全局变量,放 global里面
	// serverConfig := global.ServerConfig{}

	/*
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
	*/

	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("初始化 Nacos 配置信息: ", global.NacosConfig)

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置信息: %s %v", e.Name, global.NacosConfig)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(&global.NacosConfig)
		// fmt.Println("changed: ", serverConfig)
	})

	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace, //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})

	if err != nil {
		zap.S().Fatal("初始化 Nacos 配置信息错误: ", err)
		return
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group})

	if err != nil {
		zap.S().Fatal("转化成结构体 error : ", err)
	}

	zap.S().Info("nacos 配置字符串: ", content)

	// 这里 global.Config
	// config := &config.ServerConfig{}
	json.Unmarshal([]byte(content), global.ServerConfig)
	zap.S().Info("转化成结构体: ", global.ServerConfig)

}
