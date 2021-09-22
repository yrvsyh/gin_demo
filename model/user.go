package model

import (
	"github.com/yrvsyh/gin_demo/database"
	"github.com/yrvsyh/gin_demo/dto"
)

type (
	t struct{}

	UserModel struct {
		Name     string `gorm:"primarykey;notNull;default:''" json:"name,omitempty" form:"name"`
		Password string `gorm:"notNull;default:''" json:"password,omitempty" form:"password"`
		Email    string `gorm:"notNull;default:''" json:"email,omitempty" form:"email"`
		Role     int    `gorm:"notNull;default:1" json:"role,omitempty" form:"role"`
	}
)

const (
	USER_ROLE_ADMIN  = 0
	USER_ROLE_NORMAL = 1
)

var (
	User = t{}
)

func init() {
	database.DB.AutoMigrate(&UserModel{})
}

func (*UserModel) TableName() string { return "user" }

func (user *UserModel) ToDTO() *dto.User {
	return &dto.User{
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}

func (t) GetUserById(name string) *UserModel {
	user := &UserModel{Name: name}
	if database.DB.First(user).Error != nil {
		return nil
	}
	return user
}

func (t) InsertUser(user *UserModel) error {
	return database.DB.Create(user).Error
}
