package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go_shop/go_shop_api/goods_web/api/goods"
	"go_shop/go_shop_api/goods_web/middlewares"
)

func InitGoodsRouter(Router *gin.RouterGroup) { //*gin.RouterGroup {

	GoodsRouter := Router.Group("goods") //.Use(middlewares.JWTAuth()) 这样的话,这一组都需要登录

	zap.S().Info("配置用户相关的url")
	{
		GoodsRouter.GET("list", goods.List)

		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New) //改接口需要管理员权限

		GoodsRouter.GET("/:id", goods.Detail) //获取商品的详情

		GoodsRouter.DELETE("/:id",middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Delete) //删除商品
		GoodsRouter.PUT("/:id",middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)


	}

	//return UserRouter
}
