package entity

import (
	"searcher/internal/course/model/dto"
	"time"
)

type UpdateCourse struct {
	Id          int
	Name        string
	Description string
	Language    string
	Author      string
	Duration    string
	Rating      string
	Platform    string
	Money       string
	Link        string
	DateUpdate  string `db:"date_update"`
}

func NewUpdateCourse(course dto.PutUpdateCourseRequest) UpdateCourse {
	return UpdateCourse{
		Id:          course.Id,
		Name:        course.Name,
		Description: course.Description,
		Language:    course.Language,
		Author:      course.Author,
		Duration:    course.Duration,
		Rating:      course.Rating,
		Platform:    course.Platform,
		Money:       course.Money,
		Link:        course.Link,
		DateUpdate:  time.Now().Format("2006-01-02 15:04:05"),
	}
}