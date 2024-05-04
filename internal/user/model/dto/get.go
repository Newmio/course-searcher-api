package dto

type GetUserInfoResponse struct {
	Id                int    `json:"id" xml:"id"`
	IdUser            int    `json:"id_user" xml:"id_user"`
	Name              string `json:"name" xml:"name"`
	MiddleName        string `json:"middle_name" xml:"middle_name"`
	LastName          string `json:"last_name" xml:"last_name"`
	CourseNumber      int    `json:"course_number" xml:"course_number"`
	GroupName         string `json:"group_name" xml:"group_name"`
	Proffession       string `json:"proffession" xml:"proffession"`
	ProffessionNumber string `json:"proffession_number" xml:"proffession_number"`
}

type GetUserProfileResponse struct {
	Login  string `json:"login" xml:"login"`
	Email  string `json:"email" xml:"email"`
	Phone  string `json:"phone" xml:"phone"`
	Avatar string `json:"avatar" xml:"avatar"`
}
