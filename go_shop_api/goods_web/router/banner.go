package router

import (
	"go_shop/go_shop_api/goods_web/api/banners"
	"go_shop/go_shop_api/goods_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitBannerRouter(Router *gin.RouterGroup) {
	BannerRouter := Router.Group("banners") //.Use(middlewares.Trace())
	{
		BannerRouter.GET("", banners.List)                                                            // 轮播图列表页
		BannerRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), banners.Delete) // 删除轮播图
		BannerRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), banners.New)          //新建轮播图
		BannerRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), banners.Update)    //修改轮播图信息
	}
}
