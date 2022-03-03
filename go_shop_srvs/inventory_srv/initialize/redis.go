package initialize

import (
	"fmt"
	"go_shop/go_shop_srvs/inventory_srv/global"

	goredislib "github.com/go-redis/redis/v8"
)

func InitRedis() {

	c := global.ServerConfig.RedisInfo

	global.RedisClient = goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d", c.Host, c.Port),
	})
}
