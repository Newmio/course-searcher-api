package entity

type GetUser struct {
	Id int
	Login string
	Password string
	Email string
	Phone string
	Role string
	Active bool
	DateCreate string
	DateUpdate string
}