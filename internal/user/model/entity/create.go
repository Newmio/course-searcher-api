package entity

import (
	"searcher/internal/user/model/dto"
	"time"
)

type CreateUser struct {
	Login      string
	Password   string
	Email      string
	Role       string
	DateCreate string
}

func NewCreateUser(user dto.RegisterUserRequest) CreateUser {
	return CreateUser{
		Login:      user.Login,
		Password:   user.Password,
		Email:      user.Email,
		Role:       "user",
		DateCreate: time.Now().Format("2006-01-02 15:04:05"),
	}
}
