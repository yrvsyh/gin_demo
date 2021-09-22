package utils

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type (
	Result struct {
		Code int         `json:"code,omitempty"`
		Msg  string      `json:"msg,omitempty"`
		Data interface{} `json:"data,omitempty"`
	}

	e struct {
		code int
		msg  string
	}
)

var (
	// Common
	ERR_COMMON_ERR = e{10001, "内部错误"}
	// Auth
	ERR_USER_EXIST          = e{10001, "用户已存在"}
	ERR_USER_NOT_EXIST      = e{10002, "用户不存在"}
	ERR_PASSWORD_ERROR      = e{10003, "密码错误"}
	ERR_TOKEN_GEN_FAILD     = e{10004, "token生成失败"}
	ERR_TOKEN_DEL_FAILD     = e{10005, "token删除失败"}
	ERR_TOKEN_REFRESH_FAILD = e{10006, "token刷新失败"}
	ERR_TOKEN_PARSE_FAILD   = e{10007, "token解析失败"}
	ERR_TOKEN_INVALID       = e{10008, "token已失效"}
	ERR_REGISTER_ERROR      = e{10009, "注册失败"}
	ERR_ALREADY_LOGIN       = e{10010, "用户已登陆"}
)

var RetCode = map[int]string{
	10000: "refresh token",
}

func Success(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Result{Code: 0, Msg: msg, Data: data})
}

func Error(c *gin.Context, errCode e) {
	c.JSON(http.StatusOK, Result{Code: errCode.code, Msg: errCode.msg})
}

func ErrorWithParams(c *gin.Context, errCode e, params interface{}) {
	if params == nil {
		c.JSON(http.StatusOK, Result{Code: errCode.code, Msg: errCode.msg})
		return
	}
	var err error
	t := template.New("msg")
	t, err = t.Parse(errCode.msg)
	if err != nil {
		log.Error().Err(err).Msg("解析错误信息失败")
		ErrorWithMsg(c, ERR_COMMON_ERR.code, ERR_COMMON_ERR.msg)
		return
	}
	buf := &bytes.Buffer{}
	err = t.Execute(buf, params)
	if err != nil {
		log.Error().Err(err).Msg("渲染错误信息失败")
		ErrorWithMsg(c, ERR_COMMON_ERR.code, ERR_COMMON_ERR.msg)
		return
	}
	msg := buf.String()
	c.JSON(http.StatusOK, Result{Code: errCode.code, Msg: msg})
}

func ErrorWithMsg(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Result{Code: code, Msg: msg, Data: nil})
}

func ErrorWithData(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, Result{Code: code, Msg: msg, Data: data})
}
