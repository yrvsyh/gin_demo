package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yrvsyh/gin_demo/controller"
)

func userRouter(group *gin.RouterGroup) {
	g := group.Group("/user")
	g.GET("/:id", controller.GetUser)
}
