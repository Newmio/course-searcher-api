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
	if err := r.initTables(); err != nil {
		panic(err)
	}
	return r
}

func (r *psqlCourseRepo) GetCoursesForReport() (map[int][]entity.CourseList, error) {
	rows, err := r.psql.Queryx(`select * from course_user where topic = 'check'`)
	if err != nil {
		return nil, newm_helper.Trace(err)
	}
	defer rows.Close()

	resp := make(map[int][]entity.CourseList)

	for rows.Next() {
		row := make(map[string]interface{})

		if err := rows.MapScan(row); err != nil {
			return nil, newm_helper.Trace(err)
		}

		course, err := r.GetCourseById(int(row["id_course"].(int64)))
		if err != nil {
			return nil, newm_helper.Trace(err)
		}

		key := int(row["id_user"].(int64))

		resp[key] = append(resp[key], course)
	}

	return resp, nil
}

func (r *psqlCourseRepo) SetCheckCourseUser(courseId, userId int) error {
	str := `update course_user set topic = 'check', date_end = now() where id_user = $1 and id_course = $2`

	_, err := r.psql.Exec(str, userId, courseId)
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	return nil
}

func (r *psqlCourseRepo) CreateCourseUser(val map[string]interface{}) error {

	str := `insert into course_user(id_user, id_course)
	values($1, $2)
	on conflict (id_user, id_course) do nothing`

	_, err := r.psql.Exec(str, val["id_user"], val["id_course"])
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	return nil
}

func (r *psqlCourseRepo) GetCourseById(id int) (entity.CourseList, error) {
	var course entity.CourseList

	str := `select * from courses where id = $1`

	if err := r.psql.Get(&course, str, id); err != nil {
		if err != sql.ErrNoRows {
			return course, newm_helper.Trace(err, str)
		}
		return entity.CourseList{}, nil
	}

	return course, nil
}

func (r *psqlCourseRepo) GetCourseByLink(link string) (entity.CourseList, error) {
	var course entity.CourseList

	str := `select * from courses where link = $1`

	if err := r.psql.Get(&course, str, link); err != nil {
		if err != sql.ErrNoRows {
			return course, newm_helper.Trace(err, str)
		}
		return entity.CourseList{}, nil
	}

	return course, nil
}

func (r *psqlCourseRepo) GetCoursesByUser(id int) ([]entity.CourseList, error) {
	var courses []entity.CourseList

	var id_courses []int

	str := `select id_course from course_user where id_user = $1`

	if err := r.psql.Select(&id_courses, str, id); err != nil {
		return nil, newm_helper.Trace(err, str)
	}

	str = "select * from courses where id = $1"

	stmt, err := r.psql.Preparex(str)
	if err != nil {
		return nil, newm_helper.Trace(err)
	}
	defer stmt.Close()

	for _, v := range id_courses {
		var course entity.CourseList

		if err := stmt.Get(&course, v); err != nil {
			return nil, newm_helper.Trace(err)
		}

		courses = append(courses, course)
	}

	return courses, nil
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

	if course.Link != "" {
		if !newm_helper.СontainsSQLInjection(course.Link) {
			str += fmt.Sprintf(" icon_link = '%s',", course.Link)
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
	duration = $5, rating = $6, platform = $7, money = $8, link = $9, icon_link = $10, date_update = $11 where id = $12`

	_, err := r.psql.Exec(str, course.Name, course.Description, course.Language, course.Author,
		course.Duration, course.Rating, course.Platform, course.Money, course.Link, course.IconLink, course.DateUpdate, course.Id)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return newm_helper.Trace(fmt.Errorf("link already exists"))
		}

		return newm_helper.Trace(err, str)
	}

	return nil
}

func (r *psqlCourseRepo) CreateCourse(course entity.CreateCourse) error {
	var id int

	str := `select id from courses where link = $1`

	err := r.psql.QueryRow(str, course.Link).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			return newm_helper.Trace(err, str)
		}
	}

	if id != 0 {
		return fmt.Errorf("created")
	}

	str = `insert into courses(name, description, language, author, duration, rating, platform, money, link, icon_link, active, date_create) 
			values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	result, err := r.psql.Exec(str, course.Name, course.Description, course.Language, course.Author,
		course.Duration, course.Rating, course.Platform, course.Money, course.Link, course.IconLink, course.Active, course.DateCreate)
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

	searchStr := "%" + strings.Replace(searchValue, " ", "%", -1) + "%"

	str := `select * from courses where name ilike $1 or description ilike $2 and active = true`

	if err := r.psql.Select(&courses, str, searchStr, searchStr); err != nil {
		return nil, newm_helper.Trace(err, str)
	}

	return courses, nil
}
