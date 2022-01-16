package global

import (
	ut "github.com/go-playground/universal-translator"

	"go_shop/go_shop_api/user_web/config"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	Trans ut.Translator
)
