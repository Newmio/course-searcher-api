package dto

import (
	"fmt"
	"strings"
)

type CreatUserInfoRequest struct {
	IdUser            string    `json:"id_user"`
	Name              string `json:"name"`
	MiddleName        string `json:"middle_name"`
	LastName          string `json:"last_name"`
	CourseNumber      string    `json:"course_number"`
	GroupName         string `json:"group_name"`
	Proffession       string `json:"proffession"`
	ProffessionNumber string `json:"proffession_number"`
}

type RegisterUserRequest struct {
	Login    string `json:"login" xml:"login"`
	Email    string `json:"email" xml:"email"`
	Password string `json:"password" xml:"password"`
}

func (u RegisterUserRequest) Validate() error {
	if u.Login == "" {
		return fmt.Errorf("login is empty")
	}

	if u.Email == "" || !strings.Contains(u.Email, "@") || !strings.Contains(u.Email, ".") {
		return fmt.Errorf("email is empty")
	}

	if u.Password == "" {
		return fmt.Errorf("password is empty")
	}

	return nil
}

type RegisterUserResponse struct {
	Id int `json:"id" xml:"id"`
}

type LoginUserRequest struct {
	Login    string `json:"login" xml:"login"`
	Password string `json:"password" xml:"password"`
}

func (u LoginUserRequest) Validate() error {
	if u.Login == "" {
		return fmt.Errorf("login is empty")
	}

	if u.Password == "" {
		return fmt.Errorf("password is empty")
	}

	return nil
}

type LoginUserResponse struct {
	Access     string `json:"access" xml:"access"`
	Refresh    string `json:"refresh" xml:"refresh"`
	Exp        int    `json:"exp" xml:"exp"`
	ExpRefresh int    `json:"exp_refresh" xml:"exp_refresh"`
}

func NewLoginUserResponse(accessToken, refreshToken string, exp, expRefresh int) LoginUserResponse {
	return LoginUserResponse{
		Access:     accessToken,
		Refresh:    refreshToken,
		Exp:        exp,
		ExpRefresh: expRefresh,
	}
}
