package user

import "github.com/jmoiron/sqlx"

type IPsqlUserRepo interface {
	CreateUser(user User) error
	GetUser(login, password string) (User, error)
}

type managerUserRepo struct {
	psql IPsqlUserRepo
}

func NewManagerUserRepo(psql *sqlx.DB) IUserRepo {
	psqlRepo := NewPsqlUserRepo(psql)
	return &managerUserRepo{psql: psqlRepo}
}

func (r *managerUserRepo) GetUser(login, password string) (User, error) {
	return r.psql.GetUser(login, password)
}

func (r *managerUserRepo) CreateUser(user User) error {
	return r.psql.CreateUser(user)
}