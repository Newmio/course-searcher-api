package repository

import "github.com/jmoiron/sqlx"

type psqlFileRepo struct {
	db *sqlx.DB
}

func NewPsqlFileRepo(db *sqlx.DB) IPsqlFileRepo {
	r := &psqlFileRepo{db: db}
	r.initTables()
	return r
}
