package api

import (
	"context"
	"fmt"
	"go_shop/go_shop_api/user_web/proto"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

func GetUserList(ctx *gin.Context) {

	ip := "127.0.0.1"
	port := 50051
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())

	if err != nil {
		zap.S().Errorw("[GetUserList] 连接用户服务失败",
			"msg", err.Error)
	}

	userSrvCli := proto.NewUserClient(userConn)

	rsp, err := userSrvCli.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    0,
		PSize: 3,
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
			Gender: val.Gender,
			Mobile: val.Mobile,
		}

		result = append(result, user)
	}

	ctx.JSON(http.StatusOK, result)

}
