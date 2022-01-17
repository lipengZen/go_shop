package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"go_shop/go_shop_api/user_web/models"

)

func IsAdminAuth() gin.HandlerFunc{
	return func(ctx *gin.Context){
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*models.CustomClaims)

		if currentUser.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg":"无权限",
			})
			ctx.Abort()   // 要记得abort
			return
		}
		ctx.Next()
	}

}