package course

import "github.com/jmoiron/sqlx"

type ICourseRepo interface {
}

type courseRepo struct {
	db *sqlx.DB
}

func NewPsqlCourseRepo(db *sqlx.DB) *courseRepo {
	return &courseRepo{db: db}
}
