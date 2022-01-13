package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"strings"
	"time"

	"go_shop/go_shop_srvs/shop_srvs/user_srv/global"
	"go_shop/go_shop_srvs/shop_srvs/user_srv/model"
	"go_shop/go_shop_srvs/shop_srvs/user_srv/proto"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserServer struct {
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func ModelToResponse(user model.User) proto.UserInfoResponse {

	res := proto.UserInfoResponse{
		Id:       user.ID,
		PassWord: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
		Mobile:   user.Mobile,
	}

	// 给 grpc 的message 传 nil 容易异常
	if user.Birthday != nil {
		res.BirthDay = uint64(user.Birthday.Unix())
	}

	return res

}

func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	rsp := &proto.UserListResponse{}
	rsp.Total = uint32(result.RowsAffected)

	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)

	for _, user := range users {
		userInfoRsp := ModelToResponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}

	return rsp, nil

}

func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User

	// result := global.DB.Where("mobile=", req.Mobile).First(&user)
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil

}

func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User

	// result := global.DB.Where("mobile=", req.Mobile).First(&user)

	result := global.DB.First(&user, req.Id)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil

}

func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	//新建用户
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}

	user.Mobile = req.Mobile
	user.NickName = req.NickName

	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(req.PassWord, options)

	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil

}

func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*empty.Empty, error) {
	//个人中心更新用户
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	birthDay := time.Unix(int64(req.BirthDay), 0)
	user.NickName = req.NickName
	user.Birthday = &birthDay
	user.Gender = req.Gender

	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &empty.Empty{}, nil
}

func (s *UserServer) CheckPassword(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.PasswordCheckResponse, error) {

	// 校验用户输入的密码和 数据库中加盐的密码

	passwordInfo := strings.Split(req.EncryptedPassword, "$")

	options := &password.Options{16, 100, 32, sha512.New}
	checkFlag := password.Verify(req.Password, passwordInfo[2], passwordInfo[3], options)

	return &proto.PasswordCheckResponse{Success: checkFlag}, nil

}

/*

	UpdateUser(context.Context, *UpdateUserInfo) (*emptypb.Empty, error)
	CheckPassword(context.Context, *PasswordCheckInfo) (*PasswordCheckResponse, error)
*/
