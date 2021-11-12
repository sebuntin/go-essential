package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sebuntin/ginessential/common"
	"sebuntin/ginessential/model"
	"sebuntin/ginessential/response"
	"strings"
)

// AuthMiddleware 鉴权中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取authorization header
		tokenString := ctx.GetHeader("Authorization")
		// validate token format
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			response.Response(ctx, http.StatusUnauthorized, 401, nil, "无权限操作")
			ctx.Abort()
			return
		}
		token, claims, err := common.ParseToken(tokenString[7:])
		if err != nil {
			log.Printf("parse token failed, error: %v", err)
			response.Response(ctx, http.StatusUnauthorized, 401, nil, "无权限操作")
			ctx.Abort()
			return
		}

		// 验证通过后获取claim中的userid
		DB := common.GetDb()
		var user model.User
		DB.First(&user, claims.UserId)

		// 判断用户是否存在
		if user.ID == 0 || !token.Valid {
			response.Response(ctx, http.StatusUnauthorized, 401, nil, "无权限操作")
			ctx.Abort()
			return
		}

		// 用户存在,将user信息写入上下文
		ctx.Set("user", user)
		ctx.Next()
	}
}
