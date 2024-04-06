package entity

import (
	"searcher/internal/user/model/dto"
)

type UpdateUserPassword struct {
	Id       int
	Password string
}

type UpdateUser struct {
	Id    int
	Email string
	Phone string
}

func NewUpdateUser(user dto.UpdateUserRequest) UpdateUser {
	return UpdateUser{
		Id:    user.Id,
		Email: user.Email,
		Phone: user.Phone,
	}
}

func NewUpdateUserPassword(user dto.UpdateUserPasswordRequest) UpdateUserPassword {
	return UpdateUserPassword{
		Id:       user.Id,
		Password: user.Password,
	}
}
