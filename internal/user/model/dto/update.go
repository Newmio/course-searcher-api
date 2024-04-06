package dto

import (
	"fmt"
	"strings"
)

type UpdateUserRequest struct {
	Id    int    `json:"id" xml:"id"`
	Email string `json:"email" xml:"email"`
	Phone string `json:"phone" xml:"phone"`
}

type UpdateUserPasswordRequest struct {
	Id       int    `json:"id" xml:"id"`
	Password string `json:"password" xml:"password"`
}

func (u UpdateUserRequest) Validate() error {
	if u.Email == "" ||
		!strings.Contains(u.Email, "@") ||
		!strings.Contains(u.Email, ".") {
		return fmt.Errorf("email is empty")
	}

	if u.Phone == "" {
		return fmt.Errorf("phone is empty")
	}

	return nil
}

func (u UpdateUserPasswordRequest) Validate() error {
	if len(u.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	return nil
}
