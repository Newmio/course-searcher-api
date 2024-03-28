package repository

import (
	"searcher/internal/user/model/entity"

	"github.com/jmoiron/sqlx"
)

type IUserRepo interface {
	CreateUser(user entity.CreateUser) error
	GetUser(login, password string) (entity.GetUser, error)
}

type IPsqlUserRepo interface {
	CreateUser(user entity.CreateUser) error
	GetUser(login, password string) (entity.GetUser, error)
}

type managerUserRepo struct {
	psql IPsqlUserRepo
}

func NewManagerUserRepo(psql *sqlx.DB) IUserRepo {
	psqlRepo := NewPsqlUserRepo(psql)
	return &managerUserRepo{psql: psqlRepo}
}

func (r *managerUserRepo) GetUser(login, password string) (entity.GetUser, error) {
	return r.psql.GetUser(login, password)
}

func (r *managerUserRepo) CreateUser(user entity.CreateUser) error {
	return r.psql.CreateUser(user)
}
