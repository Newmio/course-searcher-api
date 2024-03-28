package repository

import (
	"database/sql"
	"fmt"
	"searcher/internal/course/model/entity"
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

func (r *psqlCourseRepo) UpdateCourseByParam(course entity.UpdateCourseByParam) error {
	var values []string
	
	str := "update courses set"

	i := 1
	for k, v := range course.Params {
		str += fmt.Sprintf(" %s = $%d,", k, i)
		values = append(values, v)
		i++
	}

	str = strings.TrimRight(str, ",") + fmt.Sprintf("where id = $%d", i+1)

	_, err := r.psql.DB.Exec(str, values, course.Id)
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	return nil
}

func (r *psqlCourseRepo) UpdateCourse(course entity.UpdateCourse) error {
	str := `update courses set name = $1, description = $2, language = $3, author = $4, 
	duration = $5, rating = $6, platform = $7, money = $8, link = $9, active = $10, date_update = $11 where id = $12`

	_, err := r.psql.Exec(str, course.Name, course.Description, course.Language, course.Author,
		course.Duration, course.Rating, course.Platform, course.Money, course.Link, course.Active, course.DateUpdate, course.Id)
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	return nil
}

func (r *psqlCourseRepo) CreateCourse(course entity.CreateCourse) error {
	str := `select * from courses where link = $1`

	_, err := r.psql.Exec(str, course.Link)
	if err == nil {
		if err != sql.ErrNoRows {
			return newm_helper.Trace(err, str)
		}

		return nil
	}

	str = `insert into courses(name, description, language, author, duration, rating, platform, money, link, active, date_create) 
	values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	result, err := r.psql.Exec(str, course.Name, course.Description, course.Language, course.Author,
		course.Duration, course.Rating, course.Platform, course.Money, course.Link, course.Active, course.DateCreate)
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

func (r *psqlCourseRepo) GetCourses(searchValue string) ([]entity.CourseList, error) {
	var courses []entity.CourseList

	str := `select * from courses where name ilike $1 and active = true`

	if err := r.psql.Select(&courses, str, "%"+strings.Replace(searchValue, " ", "%", -1)+"%"); err != nil {
		return nil, newm_helper.Trace(err, str)
	}

	return courses, nil
}
