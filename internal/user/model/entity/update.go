package entity

import (
	"searcher/internal/user/model/dto"
	"strconv"
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
	id, err := strconv.Atoi(user.Id)
	if err != nil {
		return UpdateUser{}
	}

	return UpdateUser{
		Id:    id,
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
