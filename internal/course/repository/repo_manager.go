package repository

import (
	"searcher/internal/course/model/entity"

	"github.com/IBM/sarama"
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
	CreateCacheCourses(courses []entity.CreateCourse, keyPrefix string) error
	GetCacheCourses(searchValue string) ([]entity.CourseList, error)
	GetCoursesByUser(id int) (map[string][]entity.CourseList, error)
	GetCourseByLink(link string) (entity.CourseList, error)
	CreateCourseEvent(value []byte) error
}

type IKafkaCourseRepo interface{
	CreateCourseEvent(value []byte) error
}

type IPsqlCourseRepo interface {
	GetCourses(searchValue string) ([]entity.CourseList, error)
	CreateCourse(course entity.CreateCourse) error
	UpdateCourse(course entity.UpdateCourse) error
	UpdateCourseByParam(course entity.UpdateCourse) error
	GetCoursesByUser(id int) ([]entity.CourseList, error)
	GetCourseByLink(link string) (entity.CourseList, error)
}

type IRedisCourseRepo interface {
	GetGlobalCourses(searchValue string) ([]entity.CourseList, error)
	UpdateGlobalCourse(course entity.UpdateCourse) error
	CreateGlobalCourse(course entity.CreateCourse) error
	UpdateGlobalCourseByParam(course entity.UpdateCourse) error
	CreateCacheCourses(courses []entity.CreateCourse, keyPrefix string) error
	GetCacheCourses(searchValue string) ([]entity.CourseList, error)
	GetCoursesByUser(id int) ([]entity.CourseList, error)
}

type IHttpCourseRepo interface {
	GetHtmlCourseInWeb(param newm_helper.Param) ([]byte, error)
}

type managerCourseRepo struct {
	psql  IPsqlCourseRepo
	redis IRedisCourseRepo
	http  IHttpCourseRepo
	kafka IKafkaCourseRepo
}

func NewManagerCourseRepo(psql *sqlx.DB, redis *redis.Client, kafka sarama.Client) ICourseRepo {
	return &managerCourseRepo{
		psql: NewPsqlCourseRepo(psql), 
		redis: NewRedisCourseRepo(redis), 
		http: NewHttpCourseRepo(),
		kafka: NewKafkaCourseRepo(kafka),
	}
}

func (r *managerCourseRepo) CreateCourseEvent(value []byte) error{
	return r.kafka.CreateCourseEvent(value)
}

func (r *managerCourseRepo) GetCoursesByUser(id int) (map[string][]entity.CourseList, error){
	psqlCourses, err := r.psql.GetCoursesByUser(id)
	if err != nil {
		return nil, newm_helper.Trace(err)
	}

	redisCourses, err := r.redis.GetCoursesByUser(id)
	if err != nil {
		return nil, newm_helper.Trace(err)
	}

	return map[string][]entity.CourseList{"psql": psqlCourses, "redis": redisCourses}, nil
}

func (r *managerCourseRepo) GetCourseByLink(link string) (entity.CourseList, error){
	return r.psql.GetCourseByLink(link)
}

func (r *managerCourseRepo) GetCacheCourses(searchValue string) ([]entity.CourseList, error) {
	return r.redis.GetCacheCourses(searchValue)
}

func (r *managerCourseRepo) CreateCacheCourses(courses []entity.CreateCourse, keyPrefix string) error {
	return r.redis.CreateCacheCourses(courses, keyPrefix)
}

func (r *managerCourseRepo) UpdateCourseByParam(course entity.UpdateCourse) error {
	if err := r.psql.UpdateCourseByParam(course); err != nil {
		return newm_helper.Trace(err)
	}

	return r.redis.UpdateGlobalCourseByParam(course)
}

func (r *managerCourseRepo) UpdateCourse(course entity.UpdateCourse) error {
	if err := r.psql.UpdateCourse(course); err != nil {
		return newm_helper.Trace(err)
	}

	return r.redis.UpdateGlobalCourse(course)
}

func (r *managerCourseRepo) CreateCourse(course entity.CreateCourse) error {
	if err := r.psql.CreateCourse(course); err != nil {
		if err.Error() == "created" {
			return nil
		}
		return newm_helper.Trace(err)
	}

	if err := r.redis.CreateGlobalCourse(course); err != nil {
		return newm_helper.Trace(err)
	}

	return nil
}

func (r *managerCourseRepo) GetShortCourses(searchValue string) ([]entity.CourseList, error) {

	courses, err := r.redis.GetGlobalCourses(searchValue)
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
