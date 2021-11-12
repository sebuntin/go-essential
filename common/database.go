package common

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"sebuntin/ginessential/model"
)

var DB *gorm.DB

func InitDb() *gorm.DB {
	driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	dbName := viper.GetString("datasource.database")
	userName := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true", userName, password, host, port, dbName, charset)
	var err error
	DB, err = gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}
	// 自动创建表
	DB.AutoMigrate(&model.User{})
	return DB
}

func GetDb() *gorm.DB {
	return DB
}
