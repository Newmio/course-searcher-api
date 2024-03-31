package repository

import (
	"searcher/internal/course/model/entity"

	"github.com/Newmio/newm_helper"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type ICourseRepo interface {
	GetShortCourses(searchValue string) ([]entity.CourseList, error)
	GetHtmlCourseInWeb(param newm_helper.Param) ([]byte, error)
	CreateCourse(course entity.CreateCourse) error
	UpdateCourse(course entity.UpdateCourse) error
	UpdateCourseByParam(course entity.UpdateCourse) error
}

type IPsqlCourseRepo interface {
	GetCourses(searchValue string) ([]entity.CourseList, error)
	CreateCourse(course entity.CreateCourse) error
	UpdateCourse(course entity.UpdateCourse) error
	UpdateCourseByParam(course entity.UpdateCourse) error
}

type IRedisCourseRepo interface {
	GetCourses(searchValue string) ([]entity.CourseList, error)
	UpdateCourse(course entity.UpdateCourse) error
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

func (r *managerCourseRepo) UpdateCourseByParam(course entity.UpdateCourse) error {
	return r.psql.UpdateCourseByParam(course)
}

func (r *managerCourseRepo) UpdateCourse(course entity.UpdateCourse) error {
	if err := r.psql.UpdateCourse(course); err != nil {
		return newm_helper.Trace(err)
	}

	return r.redis.UpdateCourse(course)
}

func (r *managerCourseRepo) CreateCourse(course entity.CreateCourse) error {
	return r.psql.CreateCourse(course)
}

func (r *managerCourseRepo) GetShortCourses(searchValue string) ([]entity.CourseList, error) {

	courses, err := r.redis.GetCourses(searchValue)
	if err != nil {
		return nil, newm_helper.Trace(err)
	}

	if courses != nil {
		return courses, nil
	}

	return r.psql.GetCourses(searchValue)
}

func (r *managerCourseRepo) GetHtmlCourseInWeb(param newm_helper.Param) ([]byte, error) {
	return r.http.GetHtmlCourseInWeb(param)
}
