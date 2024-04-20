package entity

import (
	"fmt"
	"math/rand"
	"searcher/internal/user/model/dto"
	"time"
)

type CreateUser struct {
	Login      string
	Password   string
	Email      string
	Role       string
	DateCreate string
	Avatar     string
}

func NewCreateUser(user dto.RegisterUserRequest) CreateUser {
	return CreateUser{
		Login:      user.Login,
		Password:   user.Password,
		Email:      user.Email,
		Role:       "user",
		DateCreate: time.Now().Format("2006-01-02 15:04:05"),
		Avatar: fmt.Sprintf("template/user/profile/avatars/default_%d.jpg", rand.Intn(66)+1),
	}
}
