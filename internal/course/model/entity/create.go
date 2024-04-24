package entity

import (
	"searcher/internal/course/model/dto"
	"time"
)

type CreateCourse struct {
	Name        string
	Description string
	Language    string
	Author      string
	Duration    string
	Rating      string
	Platform    string
	Money       string
	Link        string
	IconLink    string
	Active      bool
	DateCreate  string
}

func NewCreateCourse(course dto.CreateCourseRequest) CreateCourse {
	return CreateCourse{
		Name:        course.Name,
		Description: course.Description,
		Language:    course.Language,
		Author:      course.Author,
		Duration:    course.Duration,
		Rating:      course.Rating,
		Platform:    course.Platform,
		Money:       course.Money,
		Link:        course.Link,
		IconLink:    course.IconLink,
		Active:      true,
		DateCreate:  time.Now().Format("2006-01-02 15:04:05"),
	}
}

func NewCreateCourses(courses []dto.CourseList) []CreateCourse {
	var coursesResp []CreateCourse

	for _, value := range courses {
		coursesResp = append(coursesResp, CreateCourse{
			Name:        value.Name,
			Description: value.Description,
			Language:    value.Language,
			Author:      value.Author,
			Duration:    value.Duration,
			Rating:      value.Rating,
			Platform:    value.Platform,
			Money:       value.Money,
			Link:        value.Link,
			IconLink:    value.IconLink,
			Active:      true,
			DateCreate:  time.Now().Format("2006-01-02 15:04:05"),
		})
	}

	return coursesResp
}
