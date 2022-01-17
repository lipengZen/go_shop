package initialize

import (
	"github.com/gin-gonic/gin"

	"go_shop/go_shop_api/user_web/middlewares"
	"go_shop/go_shop_api/user_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	// 处理跨域问题
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("/u/v1")
	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)

	return Router
}
