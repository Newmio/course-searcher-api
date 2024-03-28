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
	Active      bool
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
		Active:      course.Active,
		DateUpdate:  time.Now().Format("2006-01-02 15:04:05"),
	}
}

type UpdateCourseByParam struct {
	Id     int               `json:"id" xml:"id"`
	Params map[string]string `json:"params" xml:"params"`
}

func NewUpdateCourseByParam(course dto.PatchUpdateCourseRequest) UpdateCourseByParam {
	var courseByParam UpdateCourseByParam

	for k, v := range course.Params {
		courseByParam.Id = course.Id
		courseByParam.Params["date_update"] = time.Now().Format("2006-01-02 15:04:05")

		switch k {

		case "name":
			courseByParam.Params["name"] = v

		case "description":
			courseByParam.Params["description"] = v

		case "language":
			courseByParam.Params["language"] = v

		case "author":
			courseByParam.Params["author"] = v

		case "duration":
			courseByParam.Params["duration"] = v

		case "rating":
			courseByParam.Params["rating"] = v

		case "platform":
			courseByParam.Params["platform"] = v

		case "money":
			courseByParam.Params["money"] = v

		case "link":
			courseByParam.Params["link"] = v

		case "active":
			courseByParam.Params["active"] = v
		}
	}

	return courseByParam
}
