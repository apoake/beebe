package service

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
	"beebe/config"
	"beebe/log"
	"os"
	"time"
)

var db *gorm.DB

func init() {
	var err error
	conf := config.GetConfig().DbConfig
	db, err = gorm.Open(conf.Dialect,
		fmt.Sprintf("%s:%s@%s/%s?%s", conf.UserName, conf.Password, conf.Host, conf.DbName, conf.ConfigStr))
	db.DB().SetMaxOpenConns(conf.MaxOpen)
	db.DB().SetMaxIdleConns(conf.MaxIdle)
	db.DB().SetConnMaxLifetime(time.Minute * conf.MaxLifeTime)
	if err != nil {
		log.Log.Fatal("db config error")
		os.Exit(-1)
	}
}

func DB() *gorm.DB {
	return db;
}
