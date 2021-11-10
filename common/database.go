package common

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"sebuntin/ginessential/model"
)

var DB *gorm.DB

func InitDb() *gorm.DB {
	driveName := "mysql"
	host := "localhost"
	userName := "root"
	port := "3306"
	DBname := "gin-essential"
	password := "654232"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true", userName, password, host, port, DBname, charset)
	var err error
	DB, err = gorm.Open(driveName, args)
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
