package course

import (
	"database/sql"
	"fmt"
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

func (r *psqlCourseRepo) CreateCourse(course Course) error {
	str := `select * from courses where link = $1`

	_, err := r.psql.Exec(str, course.Link)
	if err == nil {
		if err != sql.ErrNoRows{
			return newm_helper.Trace(err, str)
		}
	}

	str = `insert into courses(name, description, language, author, duration, rating, platform, money, link) values($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	result, err := r.psql.Exec(str, course.Name, course.Description, course.Language, course.Author,
		course.Duration, course.Rating, course.Platform, course.Money, course.Link)
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	row, err := result.RowsAffected()
	if err != nil {
		return newm_helper.Trace(err)
	}

	if row == 0 {
		return newm_helper.Trace(fmt.Errorf("bad insert course"))
	}

	return nil
}

func (r *psqlCourseRepo) GetCourse(searchValue string) ([]Course, error) {
	var courses []Course

	str := `select * from courses where name ilike $1`

	if err := r.psql.Select(&courses, str, "%"+strings.Replace(searchValue, " ", "%", -1)+"%"); err != nil {
		return nil, newm_helper.Trace(err, str)
	}

	return courses, nil
}
