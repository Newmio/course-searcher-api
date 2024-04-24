package repository

import (
	"context"
	"encoding/json"
	"fmt"
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

func (r *redisCourseRepo) GetCacheCourses(searchValue string) ([]entity.CourseList, error) {
	var courses []entity.CourseList

	c, err := r.redis.LRange(context.Background(), fmt.Sprintf("courses_%s", searchValue), 0, -1).Result()
	if err != nil {
		return nil, newm_helper.Trace(err)
	}

	for _, v := range c {
		var course entity.CourseList

		if err := json.Unmarshal([]byte(v), &course); err != nil {
			return nil, newm_helper.Trace(err)
		}

		courses = append(courses, course)
	}

	return courses, nil
}

func (r *redisCourseRepo) CreateCacheCourses(courses []entity.CreateCourse, keyPrefix string) error {
	for _, course := range courses {
		body, err := json.Marshal(course)
		if err != nil {
			return newm_helper.Trace(err)
		}

		if err := r.redis.RPush(context.Background(), fmt.Sprintf("courses_%s", keyPrefix), string(body)).Err(); err != nil {
			return newm_helper.Trace(err)
		}
	}

	return nil
}

func (r *redisCourseRepo) CreateGlobalCourse(course entity.CreateCourse) error {
	body, err := json.Marshal(course)
	if err != nil {
		return newm_helper.Trace(err)
	}

	if err := r.redis.RPush(context.Background(), "courses_global", string(body)).Err(); err != nil {
		return newm_helper.Trace(err)
	}

	return nil
}

// Ну и гавно я написал хз как по другому
func (r *redisCourseRepo) UpdateGlobalCourseByParam(course entity.UpdateCourse) error {
	var courseFromRedis entity.UpdateCourse

	c, err := r.redis.LRange(context.Background(), "courses_global", 0, -1).Result()
	if err != nil {
		return newm_helper.Trace(err)
	}

	for i, v := range c {

		if err := json.Unmarshal([]byte(v), &courseFromRedis); err != nil {
			return newm_helper.Trace(err)
		}

		if courseFromRedis.Id == course.Id {

			if course.Name != "" {
				courseFromRedis.Name = course.Name
			}

			if course.Description != "" {
				courseFromRedis.Description = course.Description
			}

			if course.Language != "" {
				courseFromRedis.Language = course.Language
			}

			if course.Author != "" {
				courseFromRedis.Author = course.Author
			}

			if course.Duration != "" {
				courseFromRedis.Duration = course.Duration
			}

			if course.Rating != "" {
				courseFromRedis.Rating = course.Rating
			}

			if course.Platform != "" {
				courseFromRedis.Platform = course.Platform
			}

			if course.Money != "" {
				courseFromRedis.Money = course.Money
			}

			if course.Link != "" {
				courseFromRedis.Link = course.Link
			}

			if course.IconLink != "" {
				courseFromRedis.IconLink = course.IconLink
			}

			if course.DateUpdate != "" {
				courseFromRedis.DateUpdate = course.DateUpdate
			}

			jsonCourse, err := json.Marshal(courseFromRedis)
			if err != nil {
				return newm_helper.Trace(err)
			}

			if err := r.redis.LSet(context.Background(), "courses_global", int64(i), jsonCourse).Err(); err != nil {
				return newm_helper.Trace(err)
			}

			break
		}
	}

	return nil
}

func (r *redisCourseRepo) UpdateGlobalCourse(course entity.UpdateCourse) error {
	var courseFromRedis entity.UpdateCourse

	c, err := r.redis.LRange(context.Background(), "courses_global", 0, -1).Result()
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
			courseFromRedis.IconLink = course.IconLink
			courseFromRedis.DateUpdate = course.DateUpdate

			jsonCourse, err := json.Marshal(courseFromRedis)
			if err != nil {
				return newm_helper.Trace(err)
			}

			if err := r.redis.LSet(context.Background(), "courses_global", int64(i), jsonCourse).Err(); err != nil {
				return newm_helper.Trace(err)
			}

			break
		}
	}

	return nil
}

func (r *redisCourseRepo) GetGlobalCourses(searchValue string) ([]entity.CourseList, error) {
	var courses []entity.CourseList

	c, err := r.redis.LRange(context.Background(), "courses_global", 0, -1).Result()
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
