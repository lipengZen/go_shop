package config

type MysqlConfig struct{
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"db" json:"db"`
	User string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type ConsulConfig struct{
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}


type ServerConfig struct{
	Name string `mapstructure:"name" json:"name"`  // 服务的名字,用于consul服务注册
	MysqlInfo MysqlConfig `mapstructure:"mysql" json:"mysql"`
	ConsulInfo ConsulConfig `mapstructure:"consul" json:"consul"`
}