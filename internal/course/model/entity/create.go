package entity

import "searcher/internal/course/model/dto"

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
	DateCreate  string
	Active      bool
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
		Active:      true,
	}
}
