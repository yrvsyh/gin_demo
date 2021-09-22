package controller

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/yrvsyh/gin_demo/middleware"
	"github.com/yrvsyh/gin_demo/model"
	"github.com/yrvsyh/gin_demo/service"
	"github.com/yrvsyh/gin_demo/utils"
)

type (
	t struct{}
)

var (
	Auth = t{}
)

func (t) Register(c *gin.Context) {
	user := &model.UserModel{}
	c.Bind(user)

	if dbUser := service.User.GetUserById(user.Name); dbUser != nil {
		utils.Error(c, utils.ERR_USER_EXIST)
		c.Abort()
		return
	}

	user.Password = hashPasswd(user.Password)

	if err := service.User.InsertUser(user); err != nil {
		utils.Error(c, utils.ERR_REGISTER_ERROR)
		c.Abort()
		return
	}

	utils.Success(c, "注册成功", user.ToDTO())
}

func (t) Login(c *gin.Context) {
	tokenString := middleware.JWTAuth.GetToken(c)
	if middleware.JWTAuth.VerifyToken(tokenString) {
		utils.Error(c, utils.ERR_ALREADY_LOGIN)
		return
	}

	user := &model.UserModel{}
	c.Bind(user)

	dbUser := service.User.GetUserById(user.Name)
	if dbUser == nil {
		utils.Error(c, utils.ERR_USER_NOT_EXIST)
		c.Abort()
		return
	}

	if !verifyPassword(user.Password, dbUser.Password) {
		utils.Error(c, utils.ERR_PASSWORD_ERROR)
		c.Abort()
		return
	}

	tokenString, err := middleware.JWTAuth.GenToken(user.Name, user.Role)
	if err != nil {
		utils.Error(c, utils.ERR_TOKEN_GEN_FAILD)
		c.Abort()
		return
	}

	utils.Success(c, "登陆成功", gin.H{"token": tokenString})
}

func (t) Refresh(c *gin.Context) {
	tokenString := middleware.JWTAuth.GetToken(c)
	newTokenString, err := middleware.JWTAuth.RefreshToken(tokenString)
	if err != nil {
		utils.Error(c, utils.ERR_TOKEN_REFRESH_FAILD)
		return
	}
	utils.Success(c, "token刷新成功", gin.H{"token": newTokenString})
}

func (t) Logout(c *gin.Context) {
	tokenString := middleware.JWTAuth.GetToken(c)
	if err := middleware.JWTAuth.DelToken(tokenString); err != nil {
		utils.Error(c, utils.ERR_TOKEN_DEL_FAILD)
		return
	}
	utils.Success(c, "注销成功", nil)
}

func verifyPassword(formPasswd string, dbPasswd string) bool {
	hashPasswd := hashPasswd(formPasswd)
	log.Debug().Str("password", formPasswd).Str("dbPasswd", dbPasswd).Str("hashPasswd", hashPasswd).Msg("PASS")
	return dbPasswd == hashPasswd
}

func hashPasswd(passwd string) string {
	sha1sum := sha1.Sum([]byte(passwd))
	return hex.EncodeToString(sha1sum[:])
}
