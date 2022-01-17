package api

import (
	"context"
	"fmt"

	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/go-redis/redis/v8"


	"go_shop/go_shop_api/user_web/forms"
	"go_shop/go_shop_api/user_web/global"
	"go_shop/go_shop_api/user_web/global/response"
	"go_shop/go_shop_api/user_web/middlewares"
	"go_shop/go_shop_api/user_web/models"
	"go_shop/go_shop_api/user_web/proto"
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

func HandleValidatorError(ctx *gin.Context, err error) {

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

	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户: %d", currentUser.ID)

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
		HandleValidatorError(ctx, err)
		return
	}

	// 校验验证码
	// 这个逻辑登录不了if store.Verify(passWordLoginForm.CaptchaId, passWordLoginForm.Captcha, false) {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"captcha": "验证码错误",
	// 	})
	// 	return
	// }
	if !store.Verify(passWordLoginForm.CaptchaId, passWordLoginForm.Captcha, true) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	addr := fmt.Sprintf("%s:%d", global.ServerConfig.UserServerInfo.Host, global.ServerConfig.UserServerInfo.Port)
	userConn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())

	if err != nil {
		zap.S().Errorw("[GetUserList] 连接用户服务失败",
			"msg", err.Error)
	}

	userSrvCli := proto.NewUserClient(userConn)

	// 登录逻辑
	if rsp, err := userSrvCli.GetUserByMobile(ctx, &proto.MobileRequest{
		Mobile: passWordLoginForm.Mobile,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登录失败",
				})
			}
			return
		}
	} else {
		if passRsp, pasErr := userSrvCli.CheckPassword(ctx, &proto.PasswordCheckInfo{
			Password:          passWordLoginForm.PassWord,
			EncryptedPassword: rsp.PassWord,
		}); pasErr != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]string{
				"password": "登录失败",
			})
		} else {
			if passRsp.Success {

				// 生成token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),               //签名的生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
						Issuer:    "shop_project",                  // 哪个机构进行的签名
					},
				}
				token, err := j.CreateToken(claims)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
				}

				ctx.JSON(http.StatusOK, gin.H{
					"id":         rsp.Id,
					"nick_name":  rsp.NickName,
					"token":      token,
					"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000, // 毫秒级别
				})
				// 	map[string]string{
				// 	"msg": "登录成功",
				// })
			} else {
				ctx.JSON(http.StatusBadRequest, map[string]string{
					"msg": "登录失败,密码错误",
				})
			}

		}
	}

}


func Register(c *gin.Context){
	//用户注册
	registerForm := forms.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	//验证码
	rdb := redis.NewClient(&redis.Options{
		Addr:fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	value, err := rdb.Get(context.Background(), registerForm.Mobile).Result()
	if err == redis.Nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"code":"验证码错误",
		})
		return
	}else{
		if value != registerForm.Code {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":"验证码错误",
			})
			return
		}
	}

	user, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: registerForm.Mobile,
		PassWord: registerForm.PassWord,
		Mobile:   registerForm.Mobile,
	})

	if err != nil {
		zap.S().Errorf("[Register] 查询 【新建用户失败】失败: %s", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}

	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:             uint(user.Id),
		NickName:       user.NickName,
		AuthorityId:    uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(), //签名的生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
			Issuer: "imooc",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":"生成token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": user.Id,
		"nick_name": user.NickName,
		"token": token,
		"expired_at": (time.Now().Unix() + 60*60*24*30)*1000,
	})
}

func GetUserDetail(ctx *gin.Context){
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户: %d", currentUser.ID)

	rsp, err := global.UserSrvClient.GetUserById(context.Background(), &proto.IdRequest{
		Id: int32(currentUser.ID),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"name":rsp.NickName,
		"birthday": time.Unix(int64(rsp.BirthDay), 0).Format("2006-01-02"),
		"gender":rsp.Gender,
		"mobile":rsp.Mobile,
	})
}


func UpdateUser(ctx *gin.Context){
	updateUserForm := forms.UpdateUserForm{}
	if err := ctx.ShouldBind(&updateUserForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户: %d", currentUser.ID)

	//将前端传递过来的日期格式转换成int
	loc, _ := time.LoadLocation("Local") //local的L必须大写
	birthDay, _ := time.ParseInLocation("2006-01-02", updateUserForm.Birthday, loc)
	_, err := global.UserSrvClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Id:       int32(currentUser.ID),
		NickName: updateUserForm.Name,
		Gender:   updateUserForm.Gender,
		BirthDay: uint64(birthDay.Unix()),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}