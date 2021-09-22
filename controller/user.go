package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/yrvsyh/gin_demo/service"
	"github.com/yrvsyh/gin_demo/utils"
)

func GetUser(c *gin.Context) {
	id := c.Param("id")
	user := service.User.GetUserById(id).ToDTO()
	log.Debug().Interface("user", user).Msg("USER")
	utils.Success(c, "", user)
}
