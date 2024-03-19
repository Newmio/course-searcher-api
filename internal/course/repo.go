package course

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Newmio/newm_helper"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type courseRepo struct {
	psql  *sqlx.DB
	redis *redis.Client
}

func NewCourseRepo(psql *sqlx.DB, redis *redis.Client) ICourseRepo {
	r := &courseRepo{psql: psql, redis: redis}
	r.initTables()
	return r
}

func (r *courseRepo) GetShortCourse(searchValue string) ([]Course, error) {

	courses, err := r.getCourseInRedis(searchValue)
	if err != nil {
		return nil, newm_helper.Trace(err)
	}

	if courses != nil {
		return courses, nil
	}

	return r.getCourseInPostgres(searchValue)
}

func (r *courseRepo) GetHtmlCourseInWeb(param newm_helper.Param) ([]byte, error) {
	status, body, err := newm_helper.RequestHTTP(param)
	if err != nil {
		return nil, newm_helper.Trace(err)
	}

	if status == 404 {
		return nil, nil
	}

	if status > 299 {
		return nil, newm_helper.Trace(fmt.Errorf("status code %d\n\n%s", status, string(body)))
	}

	return body, nil
}

func (r *courseRepo) getCourseInPostgres(searchValue string) ([]Course, error) {
	var courses []Course

	str := `select * from courses where name ilike $1`

	if err := r.psql.Select(&courses, str, "%"+strings.Replace(searchValue, " ", "%", -1)+"%"); err != nil {
		return nil, newm_helper.Trace(err)
	}

	return courses, nil
}

func (r *courseRepo) getCourseInRedis(searchValue string) ([]Course, error) {
	var courses []Course

	c, err := r.redis.LRange(context.Background(), "courses", 0, -1).Result()
	if err != nil {
		return nil, newm_helper.Trace(err)
	}

	for _, v := range c {
		var course Course

		if err := json.Unmarshal([]byte(v), &course); err != nil {
			return nil, newm_helper.Trace(err)
		}

		if strings.Contains(course.Description, searchValue) {
			courses = append(courses, course)
		}
	}

	return courses, nil
}
