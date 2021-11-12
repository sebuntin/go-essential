package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"net/http"
	"sebuntin/ginessential/common"
	"sebuntin/ginessential/model"
	"time"
)

func Register(ctx *gin.Context) {
	db := common.GetDb()
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
		return
	}
	if len(userPassword) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "密码不能少于6位",
		})
		return
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
		return
	}
	log.Println(userName, userPhone, userPassword)
	// 创建用户
	newUser := model.User{
		Name:     userName,
		Phone:    userPhone,
		Password: userPassword,
	}
	createUser(db, &newUser)

	// 发放token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		log.Printf("token generate error: %v\n", err)
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		return
	}
	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功",
	})
	return
}

// Info get user info
func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	response.Success(ctx, gin.H{"user": user}, "success")
}

func randomString(num int) string {
	rand.Seed(time.Now().Unix())
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]byte, num)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func createUser(db *gorm.DB, user *model.User) {
	db.Create(user)
}

func isPhoneNumberExist(db *gorm.DB, phone string) bool {
	var user model.User
	db.Where("phone = ?", phone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
