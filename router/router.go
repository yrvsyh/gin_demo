package router

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yrvsyh/gin_demo/middleware"
)

func InitRouter() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	r.Use(middleware.LoggerMiddleware(&log.Logger))
	r.Use(gin.Recovery())

	rootGroup := r.Group("/")
	authRouter(rootGroup)

	apiGroup := r.Group("/api/v1")
	apiGroup.Use(middleware.JWTAuthMiddleware())
	userRouter(apiGroup)

	return r
}
