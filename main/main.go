package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(20);not null"`
	Phone    string `gorm:"varchar(11);not null;unique"`
	Password string `gorm:"size:255;not null"`
}

func main() {
	db := InitDb()
	fmt.Println("hello world")
	r := gin.Default()
	r.POST("/api/auth/register", func(ctx *gin.Context) {
		// 获取参数
		userName := ctx.PostForm("name")
		userPhone := ctx.PostForm("phone")
		userPassword := ctx.PostForm("password")

		// 数据认证
		if len(userPhone) != 11 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "手机号必须为11位",
			})
		}
		if len(userPassword) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "密码不能少于6位",
			})
		}

		// 如果名称没传,则生成10位随机字符串
		if len(userName) == 0 {
			userName = RandomString(10)
		}

		// 判断手机号是否存在
		if isPhoneNumberExist(db, userPhone) {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "用户已存在",
			})
		}
		log.Println(userName, userPhone, userPassword)
		// 创建用户
		newUser := User{
			Name:     userName,
			Phone:    userPhone,
			Password: userPassword,
		}
		createUser(db, &newUser)
		// 返回结果
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "注册成功",
		})
	})
	panic(r.Run(":8008"))
}

func createUser(db *gorm.DB, user *User) {
	db.Create(user)
}

func isPhoneNumberExist(db *gorm.DB, phone string) bool {
	var user User
	db.Where("phone = ?", phone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

func RandomString(num int) string {
	rand.Seed(time.Now().Unix())
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]byte, num)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func InitDb() *gorm.DB {
	driveName := "mysql"
	host := "localhost"
	userName := "root"
	port := "3306"
	DBname := "gin-essential"
	password := "654232"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true", userName, password, host, port, DBname, charset)
	db, err := gorm.Open(driveName, args)
	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}
	// 自动创建表
	db.AutoMigrate(&User{})
	return db
}
