package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"net/http"
	"sebuntin/ginessential/common"
	"sebuntin/ginessential/model"
	"sebuntin/ginessential/response"
	"time"
)

func Register(ctx *gin.Context) {
	db := common.GetDb()
	// 获取参数
	var req = model.User{}
	invalidate := userRequestValidate(ctx, &req)
	if invalidate {
		return
	}

	// 如果名称没传,则生成10位随机字符串
	if len(req.Name) == 0 {
		req.Name = randomString(10)
	}

	// 判断手机号是否存在
	if isPhoneNumberExist(db, req.Telephone) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已存在")
		return
	}
	// 创建用户
	newUser := req
	createUser(db, &newUser)

	// 发放token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		log.Printf("token generate error: %v\n", err)
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		return
	}
	// 返回结果
	response.Success(ctx, gin.H{"token": token}, "注册成功")
	return
}

// Login user login
func Login(ctx *gin.Context) {
	db := common.GetDb()
	// 获取参数
	var req = model.User{}
	invalidate := userRequestValidate(ctx, &req)
	if invalidate {
		return
	}

	user := getUserByTelephone(db, req.Telephone)
	if user.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}

	// validate login info
	if user.Password != req.Password {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户名或密码错误")
		return
	}
	// 发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		log.Printf("token generate error: %v\n", err)
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		return
	}
	// 返回结果
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

func userRequestValidate(ctx *gin.Context, req *model.User) bool {
	err := ctx.Bind(req)
	if err != nil {
		log.Println("request error, json parse failed!")
		response.Response(ctx, http.StatusBadRequest, 400, nil, "参数有误")
		return true
	}

	// 数据验证
	if len(req.Telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return true
	}

	if len(req.Password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码长度不得小于6位数")
		return true
	}
	return false
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
	db.Where("telephone = ?", phone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

func getUserByTelephone(db *gorm.DB, phone string) model.User {
	var user model.User
	db.Where("telephone = ?", phone).First(&user)
	return user
}
