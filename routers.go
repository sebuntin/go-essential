package main

import (
	"github.com/gin-gonic/gin"
	"sebuntin/ginessential/controller"
	"sebuntin/ginessential/middleware"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("/api/auth/register", controller.Register)
	r.GET("api/auth/info", middleware.AuthMiddleware(), controller.Info)
	return r
}
