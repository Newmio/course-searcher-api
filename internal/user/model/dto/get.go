package dto

type GetUserProfileResponse struct {
	Login string `json:"login" xml:"login"`
	Email string `json:"email" xml:"email"`
	Phone string `json:"phone" xml:"phone"`
	Avatar string `json:"avatar" xml:"avatar"`
}