package initialize

import (
	"github.com/gin-gonic/gin"

	"go_shop/go_shop_api/user_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	ApiGroup := Router.Group("/u/v1")
	router.InitUserRouter(ApiGroup)

	return Router
}
