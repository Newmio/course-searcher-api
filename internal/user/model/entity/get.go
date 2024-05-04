package entity

import (
	"searcher/internal/user/model/dto"
)

type GetUserInfo struct{
	Id int
	IdUser            int `db:"id_user"`
	Name              string
	MiddleName        string `db:"middle_name"`
	LastName          string `db:"last_name"`
	CourseNumber      int    `db:"course_number"`
	GroupName         string `db:"group_name"`
	Proffession       string
	ProffessionNumber string `db:"proffession_number"`
}

func NewGetUserInfoResponse(user GetUserInfo) dto.GetUserInfoResponse {
	return dto.GetUserInfoResponse{
		Id:                user.Id,
		IdUser:            user.IdUser,
		Name:              user.Name,
		MiddleName:        user.MiddleName,
		LastName:          user.LastName,
		CourseNumber:      user.CourseNumber,
		GroupName:         user.GroupName,
		Proffession:       user.Proffession,
		ProffessionNumber: user.ProffessionNumber,
	}
}

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
