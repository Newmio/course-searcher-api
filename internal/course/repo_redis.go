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
