package user

import "github.com/jmoiron/sqlx"

type IPsqlUserRepo interface {
	
}

type psqlRepo struct {
	db *sqlx.DB
}

func NewPsqlUserRepo(db *sqlx.DB) *psqlRepo {
	r := &psqlRepo{db: db}
	r.initTables()
	return r
}
