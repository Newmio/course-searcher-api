package entity

import "searcher/internal/course/model/dto"

type CourseList struct {
	Id int
	Name        string
	Description string
	Language    string
	Author      string
	Duration    string
	Rating      string
	Platform    string
	Money       string
	Link        string
	Active      bool
	IconLink    string `db:"icon_link"`
	DateCreate  string `db:"date_create"`
	DateUpdate  string `db:"date_update"`
}

func NewCourseList(coruses []CourseList) []dto.CourseList {
	var coursesResp []dto.CourseList

	for _, value := range coruses {
		coursesResp = append(coursesResp, dto.CourseList{
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
		})
	}

	return coursesResp
}

func NewCourseListResponse(courses []CourseList) dto.CourseListResponse {
	return dto.CourseListResponse{
		Courses: NewCourseList(courses),
		Count:   len(courses),
	}
}
