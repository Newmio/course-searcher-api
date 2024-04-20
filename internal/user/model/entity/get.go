package entity

import (
	"searcher/internal/user/model/dto"
)

type GetUser struct {
	Id         int
	Login      string
	Password   string
	Email      string
	Phone      string
	Role       string
	Avatar     string
	DateCreate string `db:"date_create"`
}

func NewUserProfileResponse(user GetUser) dto.GetUserProfileResponse {
	return dto.GetUserProfileResponse{
		Login:  user.Login,
		Email:  user.Email,
		Phone:  user.Phone,
		Avatar: user.Avatar,
	}
}
