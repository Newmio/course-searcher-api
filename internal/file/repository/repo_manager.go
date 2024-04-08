package repository

import "github.com/jmoiron/sqlx"

type IFileRepo interface {
}

type IPsqlFileRepo interface {
}

type managerFileRepo struct {
	psql IPsqlFileRepo
}

func NewManagerFileRepo(psql *sqlx.DB) IFileRepo {
	psqlRepo := NewPsqlFileRepo(psql)
	return &managerFileRepo{psql: psqlRepo}
}
