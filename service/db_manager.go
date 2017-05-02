package service

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
	"beebe/config"
	"beebe/log"
	"os"
)

var db *gorm.DB

func init() {
	var err error
	config := config.GetConfig().DbConfig
	db, err = gorm.Open(config.Dialect,
		fmt.Sprintf("%s:%s@%s/%s?%s", config.UserName, config.Password, config.Host, config.DbName, config.ConfigStr))
	if err != nil {
		log.Logger().Fatal("db config error")
		os.Exit(-1)
	}
}

func DB() *gorm.DB {
	return db;
}
