package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go_shop/go_shop_api/user_web/api"
	"go_shop/go_shop_api/user_web/middlewares"
)

func InitUserRouter(Router *gin.RouterGroup) { //*gin.RouterGroup {

	UserRouter := Router.Group("user") //.Use(middlewares.JWTAuth()) 这样的话,这一组都需要登录

	zap.S().Info("配置用户相关的url")
	{
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("pwd_login", api.PassWordLogin)

		UserRouter.POST("register", api.Register)

		UserRouter.GET("detail", middlewares.JWTAuth(), api.GetUserDetail)
		UserRouter.PATCH("update", middlewares.JWTAuth(), api.UpdateUser)

	}

	//return UserRouter
}
