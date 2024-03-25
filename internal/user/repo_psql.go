package user

import (
	"database/sql"

	"github.com/Newmio/newm_helper"
	"github.com/jmoiron/sqlx"
)

type psqlUserRepo struct {
	db *sqlx.DB
}

func NewPsqlUserRepo(db *sqlx.DB) IPsqlUserRepo {
	r := &psqlUserRepo{db: db}
	r.initTables()
	return r
}

func (r *psqlUserRepo) GetUser(login, password string) (User, error) {
	var user User

	str := `select * from users where login = $1 and password = $2`

	if err := r.db.Get(&user, str, login, password); err != nil {
		if err == sql.ErrNoRows {
			return User{}, nil
		}
		return User{}, newm_helper.Trace(err, str)
	}

	return user, nil
}

func (db *psqlUserRepo) CreateUser(user User) error {
	str := `insert into users(login, password, email, phone, role) values($1, $2, $3, $4, $5)`

	_, err := db.db.Exec(str, user.Login, user.Password, user.Email, user.Phone, user.Role)
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	return nil
}
