package service

import (
	"github.com/yrvsyh/gin_demo/model"
)

type (
	t struct{}
)

var (
	User = t{}
)

func (t) GetUserById(id string) *model.UserModel {
	return model.User.GetUserById(id)
}

func (t) InsertUser(user *model.UserModel) error {
	return model.User.InsertUser(user)
}
