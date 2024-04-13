package repository

import (
	"database/sql"
	"searcher/internal/user/model/entity"

	"github.com/Newmio/newm_helper"
	"github.com/jmoiron/sqlx"
)

type psqlUserRepo struct {
	db *sqlx.DB
}

func NewPsqlUserRepo(db *sqlx.DB) IPsqlUserRepo {
	r := &psqlUserRepo{db: db}
	if err := r.initTables(); err != nil{
		panic(err)
	}
	return r
}

func (r *psqlUserRepo) UpdateUserRole(role string, userId int)error{
	str := `update users set role = $1 where id = $2`

	_, err := r.db.Exec(str, role, userId)
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	return nil
}

func (r *psqlUserRepo) UpdatePassword(user entity.UpdateUserPassword) error {
	str := `update users set password = $1 where id = $2`

	_, err := r.db.Exec(str, user.Password, user.Id)
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	return nil
}

func (r *psqlUserRepo) UpdateUser(user entity.UpdateUser) error {
	str := `update users set email = $1, phone = $2 where id = $3`

	_, err := r.db.Exec(str, user.Email, user.Phone, user.Id)
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	return nil
}

func (r *psqlUserRepo) GetUser(login, password string) (entity.GetUser, error) {
	var user entity.GetUser

	str := `select * from users where login = $1 and password = $2`

	if err := r.db.Get(&user, str, login, password); err != nil {
		if err == sql.ErrNoRows {
			return entity.GetUser{}, nil
		}
		return entity.GetUser{}, newm_helper.Trace(err, str)
	}

	return user, nil
}

func (db *psqlUserRepo) CreateUser(user entity.CreateUser) error {
	str := `insert into users(login, password, email, role, date_create) values($1, $2, $3, $4, $5)`

	_, err := db.db.Exec(str, user.Login, user.Password, user.Email, user.Role, user.DateCreate)
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	return nil
}
