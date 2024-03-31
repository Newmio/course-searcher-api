package entity

type GetUser struct {
	Id int
	Login string
	Password string
	Email string
	Phone string
	Role string
	Active bool
	DateCreate string `db:"date_create"`
	DateUpdate string `db:"date_update"`
}