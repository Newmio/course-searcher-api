package user

import "github.com/jmoiron/sqlx"

type psqlRepo struct {
	psql *sqlx.DB
}

func NewPsqlUserRepo(psql *sqlx.DB) IUserRepo {
	r := &psqlRepo{psql: psql}
	r.initTables()
	return r
}
