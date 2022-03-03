package global

import (
	"go_shop/go_shop_srvs/inventory_srv/config"

	"gorm.io/gorm"
	goredislib "github.com/go-redis/redis/v8"

)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	// NacosConfig config.NacosConfig
	NacosConfig *config.NacosConfig = &config.NacosConfig{}
	RedisClient *goredislib.Client
)

func init() {
	/*
		dsn := "root:root@tcp(127.0.0.1:3306)/shop_inventory_srv?charset=utf8mb4&parseTime=True&loc=Local"

		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // 慢 SQL 阈值
				LogLevel:      logger.Info, // Log level
				Colorful:      true,        // 禁用彩色打印
			},
		)

		// 全局模式
		var err error
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			// 这个是将表生成为单数，否则生成 users
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger: newLogger,
		})
		if err != nil {
			fmt.Println("get DB Error:", err)
		}
	*/
}
