package repository

import (
	"context"
	"encoding/json"
	"searcher/internal/course/model/entity"
	"strings"

	"github.com/Newmio/newm_helper"
	"github.com/go-redis/redis/v8"
)

type redisCourseRepo struct {
	redis *redis.Client
}

func NewRedisCourseRepo(redis *redis.Client) IRedisCourseRepo {
	return &redisCourseRepo{redis: redis}
}

func (r *redisCourseRepo) UpdateCourseByParam(course entity.UpdateCourseByParam) error {
	// TODO доделать апдейт по параметрам (хер знает как это сделать в редисе)
	return nil
}

func (r *redisCourseRepo) UpdateCourse(course entity.UpdateCourse) error {
	var courseFromRedis entity.UpdateCourse

	c, err := r.redis.LRange(context.Background(), "courses", 0, -1).Result()
	if err != nil {
		return newm_helper.Trace(err)
	}

	for i, v := range c {
		
		if err := json.Unmarshal([]byte(v), &courseFromRedis); err != nil {
			return newm_helper.Trace(err)
		}

		if courseFromRedis.Id == course.Id {
			courseFromRedis.Name = course.Name
			courseFromRedis.Description = course.Description
			courseFromRedis.Language = course.Language
			courseFromRedis.Author = course.Author
			courseFromRedis.Duration = course.Duration
			courseFromRedis.Rating = course.Rating
			courseFromRedis.Platform = course.Platform
			courseFromRedis.Money = course.Money
			courseFromRedis.Link = course.Link
			courseFromRedis.Active = course.Active
			courseFromRedis.DateUpdate = course.DateUpdate

			jsonCourse, err := json.Marshal(courseFromRedis)
			if err != nil {
				return newm_helper.Trace(err)
			}

			if err := r.redis.LSet(context.Background(), "courses", int64(i), jsonCourse).Err(); err != nil {
				return newm_helper.Trace(err)
			}

			break
		}
	}

	return nil
}

func (r *redisCourseRepo) GetCourses(searchValue string) ([]entity.CourseList, error) {
	var courses []entity.CourseList

	c, err := r.redis.LRange(context.Background(), "courses", 0, -1).Result()
	if err != nil {
		return nil, newm_helper.Trace(err)
	}

	for _, v := range c {
		var course entity.CourseList

		if err := json.Unmarshal([]byte(v), &course); err != nil {
			return nil, newm_helper.Trace(err)
		}

		strName := strings.ToLower(course.Name)
		strDescription := strings.ToLower(course.Description)
		strSearchValue := strings.ToLower(searchValue)

		if strings.Contains(strName, strSearchValue) || strings.Contains(strDescription, strSearchValue) {
			courses = append(courses, course)
		}
	}

	return courses, nil
}
