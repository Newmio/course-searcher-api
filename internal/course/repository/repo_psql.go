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

func (r *psqlCourseRepo) UpdateCourseByParam(course entity.UpdateCourse) error {
	str := "update courses set"

	if course.Name != "" {
		if !newm_helper.СontainsSQLInjection(course.Name) {
			str += fmt.Sprintf(" name = '%s',", course.Name)
		}
	}

	if course.Description != "" {
		if !newm_helper.СontainsSQLInjection(course.Description) {
			str += fmt.Sprintf(" description = '%s',", course.Description)
		}
	}

	if course.Language != "" {
		if !newm_helper.СontainsSQLInjection(course.Language) {
			str += fmt.Sprintf(" language = '%s',", course.Language)
		}
	}

	if course.Author != "" {
		if !newm_helper.СontainsSQLInjection(course.Author) {
			str += fmt.Sprintf(" author = '%s',", course.Author)
		}
	}

	if course.Duration != "" {
		if !newm_helper.СontainsSQLInjection(course.Duration) {
			str += fmt.Sprintf(" duration = '%s',", course.Duration)
		}
	}

	if course.Rating != "" {
		if !newm_helper.СontainsSQLInjection(course.Rating) {
			str += fmt.Sprintf(" rating = '%s',", course.Rating)
		}
	}

	if course.Platform != "" {
		if !newm_helper.СontainsSQLInjection(course.Platform) {
			str += fmt.Sprintf(" platform = '%s',", course.Platform)
		}
	}

	if course.Money != "" {
		if !newm_helper.СontainsSQLInjection(course.Money) {
			str += fmt.Sprintf(" money = '%s',", course.Money)
		}
	}

	if course.Link != "" {
		if !newm_helper.СontainsSQLInjection(course.Link) {
			str += fmt.Sprintf(" link = '%s',", course.Link)
		}
	}

	if course.DateUpdate != "" {
		if !newm_helper.СontainsSQLInjection(course.DateUpdate) {
			str += fmt.Sprintf(" date_update = '%s',", course.DateUpdate)
		}
	}

	str = strings.TrimRight(str, ",") + " where id = $1"

	_, err := r.psql.DB.Exec(str, course.Id)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return newm_helper.Trace(fmt.Errorf("link already exists"))
		}

		return newm_helper.Trace(err, str)
	}

	return nil
}

func (r *psqlCourseRepo) UpdateCourse(course entity.UpdateCourse) error {
	str := `update courses set name = $1, description = $2, language = $3, author = $4, 
	duration = $5, rating = $6, platform = $7, money = $8, link = $9, date_update = $10 where id = $11`

	_, err := r.psql.Exec(str, course.Name, course.Description, course.Language, course.Author,
		course.Duration, course.Rating, course.Platform, course.Money, course.Link, course.DateUpdate, course.Id)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return newm_helper.Trace(fmt.Errorf("link already exists"))
		}

		return newm_helper.Trace(err, str)
	}

	return nil
}

func (r *psqlCourseRepo) CreateCourse(course entity.CreateCourse) error {
	str := `select * from courses where link = $1`

	_, err := r.psql.Exec(str, course.Link)
	if err != nil {
		if err != sql.ErrNoRows {
			return newm_helper.Trace(err, str)
		}
		return nil
	}

	str = `insert into courses(name, description, language, author, duration, rating, platform, money, link, active, date_create) 
	values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

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
