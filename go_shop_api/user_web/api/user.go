package api

import (
	"context"
	"fmt"
	"go_shop/go_shop_api/user_web/forms"
	"go_shop/go_shop_api/user_web/global"
	"go_shop/go_shop_api/user_web/proto"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-playground/validator/v10"

	"go_shop/go_shop_api/user_web/global/response"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	// 将grpc 的 code 转换成 http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "其它错误, code: " + e.Code().String() + "msg: " + e.Message(),
				})
				return
			}
		}
	}
}

func HandlerValidatorError(ctx *gin.Context, err error) {

	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})

}

func GetUserList(ctx *gin.Context) {

	addr := fmt.Sprintf("%s:%d", global.ServerConfig.UserServerInfo.Host, global.ServerConfig.UserServerInfo.Port)

	// ip := "127.0.0.1"
	// port := 50051
	// fmt.Sprintf("%s:%d", ip, port)
	userConn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())

	if err != nil {
		zap.S().Errorw("[GetUserList] 连接用户服务失败",
			"msg", err.Error)
	}

	userSrvCli := proto.NewUserClient(userConn)

	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "0")
	pSizeInt, _ := strconv.Atoi(pSize)

	rsp, err := userSrvCli.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),    //0,
		PSize: uint32(pSizeInt), // 3,
	})

	if err != nil {
		zap.S().Warnf("[GetUserList] 查询用户列表失败",
			"msg", err.Error)
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, val := range rsp.Data {
		/*
			data := make(map[string]interface{})

			data["id"] = val.Id
			data["name"] = val.NickName
			data["birthday"] = val.BirthDay
			data["gender"] = val.Gender
			data["mobile"] = val.Mobile
		*/

		user := response.UserResponse{
			Id:       val.Id,
			NickName: val.NickName,
			// Birthday: time.Time(time.Unix(int64(val.BirthDay), 0)).Format("2006-01-02"),
			Birthday: response.JsonTime(time.Unix(int64(val.BirthDay), 0)),
			Gender:   val.Gender,
			Mobile:   val.Mobile,
		}

		result = append(result, user)
	}

	ctx.JSON(http.StatusOK, result)

}

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func PassWordLogin(ctx *gin.Context) {

	// 进行表单验证
	passWordLoginForm := forms.PassWordLoginForm{}
	if err := ctx.ShouldBind(&passWordLoginForm); err != nil {
		HandlerValidatorError(ctx, err)
		return
	}

}
