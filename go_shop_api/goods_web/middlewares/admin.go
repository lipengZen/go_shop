package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go_shop/go_shop_api/goods_web/models"
)

func IsAdminAuth() gin.HandlerFunc {
	// 将一些共用的代码抽出来然后共用 -> 版本管理
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*models.CustomClaims)

		if currentUser.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			ctx.Abort() // 要记得abort
			return
		}
		ctx.Next()
	}

}
