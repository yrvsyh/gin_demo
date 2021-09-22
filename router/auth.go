package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yrvsyh/gin_demo/controller"
	"github.com/yrvsyh/gin_demo/middleware"
)

func authRouter(group *gin.RouterGroup) {
	g := group.Group("/auth")
	g.POST("/register", controller.Auth.Register)
	g.POST("/login", controller.Auth.Login)
	g.GET("/refresh", controller.Auth.Refresh)
	g.GET("/logout", middleware.JWTAuthMiddleware(), controller.Auth.Logout)
}
