package repository

import (
	"searcher/internal/user/model/entity"

	"github.com/jmoiron/sqlx"
)

type IUserRepo interface {
	CreateUser(user entity.CreateUser) error
	GetUser(login, password string) (entity.GetUser, error)
	UpdateUser(user entity.UpdateUser) error
	UpdatePassword(user entity.UpdateUserPassword) error
	GetUserById(id int) (entity.GetUser, error)
}

type IPsqlUserRepo interface {
	CreateUser(user entity.CreateUser) error
	GetUser(login, password string) (entity.GetUser, error)
	UpdateUser(user entity.UpdateUser) error
	UpdatePassword(user entity.UpdateUserPassword) error
	GetUserById(id int) (entity.GetUser, error)
}

type managerUserRepo struct {
	psql IPsqlUserRepo
}

func NewManagerUserRepo(psql *sqlx.DB) IUserRepo {
	psqlRepo := NewPsqlUserRepo(psql)
	return &managerUserRepo{psql: psqlRepo}
}

func (r *managerUserRepo) GetUserById(id int) (entity.GetUser, error) {
	return r.psql.GetUserById(id)
}

func (r *managerUserRepo) UpdatePassword(user entity.UpdateUserPassword) error {
	return r.psql.UpdatePassword(user)
}

func (r *managerUserRepo) UpdateUser(user entity.UpdateUser) error {
	return r.psql.UpdateUser(user)
}

func (r *managerUserRepo) GetUser(login, password string) (entity.GetUser, error) {
	return r.psql.GetUser(login, password)
}

func (r *managerUserRepo) CreateUser(user entity.CreateUser) error {
	return r.psql.CreateUser(user)
}
