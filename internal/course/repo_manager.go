package course

import (
	"github.com/Newmio/newm_helper"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type IPsqlCourseRepo interface {
	GetCourse(searchValue string) ([]Course, error)
	CreateCourse(course Course) error
}

type IRedisCourseRepo interface {
	GetCourse(searchValue string) ([]Course, error)
}

type IHttpCourseRepo interface {
	GetHtmlCourseInWeb(param newm_helper.Param) ([]byte, error)
}

type managerCourseRepo struct {
	psql  IPsqlCourseRepo
	redis IRedisCourseRepo
	http  IHttpCourseRepo
}

func NewManagerCourseRepo(psql *sqlx.DB, redis *redis.Client) ICourseRepo {
	psqlRepo := NewPsqlCourseRepo(psql)
	redisRepo := NewRedisCourseRepo(redis)
	httpRepo := NewHttpCourseRepo()
	return &managerCourseRepo{psql: psqlRepo, redis: redisRepo, http: httpRepo}
}

func (r *managerCourseRepo) CreateCourse(course Course) error {
	return r.psql.CreateCourse(course)
}

func (r *managerCourseRepo) GetShortCourse(searchValue string) ([]Course, error) {

	courses, err := r.redis.GetCourse(searchValue)
	if err != nil {
		return nil, newm_helper.Trace(err)
	}

	if courses != nil {
		return courses, nil
	}

	return r.psql.GetCourse(searchValue)
}

func (r *managerCourseRepo) GetHtmlCourseInWeb(param newm_helper.Param) ([]byte, error){
	return r.http.GetHtmlCourseInWeb(param)
}