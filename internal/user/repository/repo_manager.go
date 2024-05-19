package repository

import (
	"searcher/internal/user/model/entity"

	"github.com/Newmio/newm_helper"
	"github.com/jmoiron/sqlx"
)

type IUserRepo interface {
	CreateUser(user entity.CreateUser) error
	GetUser(login, password string) (entity.GetUser, error)
	UpdateUser(user entity.UpdateUser) error
	UpdatePassword(user entity.UpdateUserPassword) error
	GetUserById(id int) (entity.GetUser, error)
	UpdateUserInfo(info entity.CreateUserInfo) error
	GetUserInfo(userId int) (entity.GetUserInfo, error)
	UpdateUserAvatar(userId int, avatar string) error
	GetAllAdmins()([]entity.GetUser, error)
	GetUsersByGroupName(group string)([]entity.GetUser, error)
}

type IPsqlUserRepo interface {
	CreateUser(user entity.CreateUser) error
	GetUser(login, password string) (entity.GetUser, error)
	UpdateUser(user entity.UpdateUser) error
	UpdatePassword(user entity.UpdateUserPassword) error
	GetUserById(id int) (entity.GetUser, error)
	CreateUserInfo(info entity.CreateUserInfo) error
	UpdateUserInfo(info entity.CreateUserInfo) error
	GetUserInfo(userId int) (entity.GetUserInfo, error)
	UpdateUserAvatar(userId int, avatar string) error
	GetAllAdmins()([]entity.GetUser, error)
	GetUsersByGroupName(group string)([]entity.GetUser, error)
}

type managerUserRepo struct {
	psql IPsqlUserRepo
}

func NewManagerUserRepo(psql *sqlx.DB) IUserRepo {
	return &managerUserRepo{psql: NewPsqlUserRepo(psql)}
}

func (r *managerUserRepo) GetUsersByGroupName(group string)([]entity.GetUser, error){
	return r.psql.GetUsersByGroupName(group)
}

func (r *managerUserRepo) GetAllAdmins()([]entity.GetUser, error){
	return r.psql.GetAllAdmins()
}

func (r *managerUserRepo) UpdateUserInfo(info entity.CreateUserInfo) error {
	if err := r.psql.CreateUserInfo(info); err != nil {
		if err.Error() != "info exists" {
			return newm_helper.Trace(err)
		}
	}
	return r.psql.UpdateUserInfo(info)
}

func (r *managerUserRepo) UpdateUserAvatar(userId int, avatar string) error{
	return r.psql.UpdateUserAvatar(userId, avatar)
}

func (r *managerUserRepo) GetUserInfo(userId int) (entity.GetUserInfo, error){
	return r.psql.GetUserInfo(userId)
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
