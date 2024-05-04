package entity

import (
	"fmt"
	"math/rand"
	"searcher/internal/user/model/dto"
	"strconv"
	"time"
)

type CreateUserInfo struct {
	IdUser            int `db:"id_user"`
	Name              string
	MiddleName        string `db:"middle_name"`
	LastName          string `db:"last_name"`
	CourseNumber      int    `db:"course_number"`
	GroupName         string `db:"group_name"`
	Proffession       string
	ProffessionNumber string `db:"proffession_number"`
}

func NewCreateUserInfo(user dto.CreatUserInfoRequest) CreateUserInfo {
	userId, err := strconv.Atoi(user.IdUser)
	if err != nil {
		return CreateUserInfo{}
	}

	courseNumber, err := strconv.Atoi(user.CourseNumber)
	if err != nil {
		return CreateUserInfo{}
	}
	return CreateUserInfo{
		IdUser:            userId,
		Name:              user.Name,
		MiddleName:        user.MiddleName,
		LastName:          user.LastName,
		CourseNumber:      courseNumber,
		GroupName:         user.GroupName,
		Proffession:       user.Proffession,
		ProffessionNumber: user.ProffessionNumber,
	}
}

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
		Avatar:     fmt.Sprintf("template/user/profile/avatars/default_%d.jpg", rand.Intn(66)+1),
	}
}
