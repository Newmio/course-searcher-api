package dto

import "fmt"

type RegisterUserRequest struct {
	Login    string `json:"login" xml:"login"`
	Password string `json:"password" xml:"password"`
	Email    string `json:"email" xml:"email"`
}

func (u RegisterUserRequest) Validate() error {
	if u.Login == "" {
		return fmt.Errorf("login is empty")
	}

	if u.Password == "" {
		return fmt.Errorf("password is empty")
	}

	if u.Email == "" {
		return fmt.Errorf("email is empty")
	}
	
	return nil
}

type RegisterUserResponse struct {
	Id int `json:"id" xml:"id"`
}
