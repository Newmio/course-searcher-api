package course

import (
	"context"
	"encoding/json"
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

// func (r *redisCourseRepo) CheckExistsCourse(id int) (bool, error) {
// 	c, err := r.redis.LRange(context.Background(), "courses", 0, -1).Result()
// 	if err != nil {
// 		return false, newm_helper.Trace(err)
// 	}

// 	for _, v := range c {
// 		var courseFromRedis Course

// 		if err := json.Unmarshal([]byte(v), &courseFromRedis); err != nil {
// 			return false, newm_helper.Trace(err)
// 		}

// 		if courseFromRedis.Id == id {
// 			return true, nil
// 		}
// 	}

// 	return false, nil
// }

func (r *redisCourseRepo) UpdateCourse(course Course) error {
	var courseFromRedis Course

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

func (r *redisCourseRepo) GetCourse(searchValue string) ([]Course, error) {
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

		strName := strings.ToLower(course.Name)
		strDescription := strings.ToLower(course.Description)
		strSearchValue := strings.ToLower(searchValue)

		if strings.Contains(strName, strSearchValue) || strings.Contains(strDescription, strSearchValue) {
			courses = append(courses, course)
		}
	}

	return courses, nil
}
