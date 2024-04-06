package entity

type GetUser struct {
	Id int
	Login string
	Password string
	Email string
	Phone string
	Role string
	DateCreate string `db:"date_create"`
	DateUpdate string `db:"date_update"`
}