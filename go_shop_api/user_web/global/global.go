package global

import (
	ut "github.com/go-playground/universal-translator"

	"go_shop/go_shop_api/user_web/config"
	"go_shop/go_shop_api/user_web/proto"

)



var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	Trans ut.Translator

	UserSrvClient proto.UserClient
)


