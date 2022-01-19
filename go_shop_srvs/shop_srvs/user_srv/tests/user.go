package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"go_shop/go_shop_srvs/shop_srvs/user_srv/proto"
)

var conn *grpc.ClientConn
var userClient proto.UserClient

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())
	//grpc.WithInsecure)

	if err != nil {
		fmt.Println("conn failed: ", err)
		return
	}

	userClient = proto.NewUserClient(conn)

}

func TestGetUserList() { //t *testing.T

	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    2,
		PSize: 4,
	})
	if err != nil {
		// t.Fatal(err)
		fmt.Println(err)
	}

	fmt.Println(rsp.Total)
	for _, user := range rsp.Data {

		fmt.Println(user.Mobile, user.NickName, user.PassWord)

		checkRsp, err := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:          "123456",
			EncryptedPassword: user.PassWord,
		})
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(checkRsp.Success)

	}

}

/* 
  init方法中的 50051 已经是 获取可用端口了
*/

// func main() {

// 	Init()
// 	defer conn.Close()

// 	TestGetUserList() //&testing.T{})

// }
