package user

import "github.com/golang-jwt/jwt"

const (
	SALT = "fwfjsndfwdqwdqwuriiotncna23219nsncjancasncuenfen834832u0423u094239jdjsanjsiqepee33425e1rqwftdyvghsuqw78e6trgdhbsuw3e7ref"
	SIGNKEY = "ncaeuwbcewr43943qfb8340hdq4t93q48ugmx9bgbfbydsbufxy6g37b2qg6fxbg67b4gfbxq6xf7x349q6gf76gew7gqf67xg4qf76g437fggf6gwefg"
)

type User struct {
	Id       int    `db:"id" json:"id" xml:"id"`
	Login    string `db:"login" json:"login" xml:"login"`
	Password string `db:"password" json:"password" xml:"password"`
	Email    string `db:"email" json:"email" xml:"email"`
	Phone    string `db:"phone" json:"phone" xml:"phone"`
	Role     string `db:"role" json:"role" xml:"role"`
}

type Person struct {
	Id             int    `db:"id" json:"id" xml:"id"`
	FirstName      string `db:"first_name" json:"first_name" xml:"first_name"`
	MiddleName     string `db:"middle_name" json:"middle_name" xml:"middle_name"`
	LastName       string `db:"last_name" json:"last_name" xml:"last_name"`
	City           string `db:"city" json:"city" xml:"city"`
	Street         string `db:"street" json:"street" xml:"street"`
	UniversityRole string `db:"university_role" json:"university_role" xml:"university_role"`
	Id_user        int    `db:"id_user" json:"id_user" xml:"id_user"`
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int    `json:"id_user"`
	Role   string `json:"role"`
}
