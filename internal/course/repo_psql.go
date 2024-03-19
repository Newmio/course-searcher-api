package course

import (
	"strings"

	"github.com/Newmio/newm_helper"
	"github.com/jmoiron/sqlx"
)

type psqlCourseRepo struct {
	psql *sqlx.DB
}

func NewPsqlCourseRepo(psql *sqlx.DB) IPsqlCourseRepo {
	r := &psqlCourseRepo{psql: psql}
	r.initTables()
	return r
}

func (r *psqlCourseRepo) GetCourseInPostgres(searchValue string) ([]Course, error) {
	var courses []Course

	str := `select * from courses where name ilike $1`

	if err := r.psql.Select(&courses, str, "%"+strings.Replace(searchValue, " ", "%", -1)+"%"); err != nil {
		return nil, newm_helper.Trace(err)
	}

	return courses, nil
}
