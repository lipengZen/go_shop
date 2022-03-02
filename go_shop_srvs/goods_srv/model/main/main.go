package main

import (
	"go_shop/go_shop_srvs/goods_srv/model"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {

	dsn := "root:root@tcp(127.0.0.1:3306)/shop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	// 全局模式
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 这个是将表生成为单数，否则生成 users
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	//设置全局的logger，这个logger在我们执行每个sql语句的时候会打印每一行sql
	//sql才是最重要的，本着这个原则我尽量的给大家看到每个api背后的sql语句是什么

	//定义一个表结构， 将表结构直接生成对应的表 - migrations
	// 迁移 schema
	_ = db.AutoMigrate(&model.Category{}, &model.Brands{}, &model.GoodsCategoryBrand{}, &model.Banner{}, &model.Goods{}) //此处应该有sql语句

	// fmt.Println(genMd5("123457"))

	// salt, encodedPwd := password.Encode("generic password", nil)
	// check := password.Verify("generic password", salt, encodedPwd, nil)
	// fmt.Println(check) // true

	// fmt.Println(salt, encodedPwd)

	// Using custom options
	// options := &password.Options{16, 100, 32, sha512.New}
	// salt, encodedPwd := password.Encode("generic password", options)

	// newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	// fmt.Println("newPassword", newPassword)

	// passwordInfo := strings.Split(newPassword, "$")
	// fmt.Println(password Info)
	// check := password.Verify("generic password", passwordInfo[2], passwordInfo[3], options)
	// fmt.Println(check) // true

	// check = password.Verify("generic password", salt, encodedPwd, options)

	// // fmt.Println(salt)
	// // fmt.Println(encodedPwd)

	// fmt.Println(check) // true

}
