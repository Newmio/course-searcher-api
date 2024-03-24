package course

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Newmio/newm_helper"
	"github.com/PuerkitoBio/goquery"
)

type ICourseRepo interface {
	GetShortCourse(valueSearch string) ([]Course, error)
	GetHtmlCourseInWeb(param newm_helper.Param) ([]byte, error)
	CreateCourse(course Course) error
	UpdateCourse(course Course) error
}

type courseService struct {
	r ICourseRepo
}

func NewCourseService(r ICourseRepo) ICourseService {
	return &courseService{r: r}
}

func (s *courseService) UpdateCourse(course Course) error {
	return s.r.UpdateCourse(course)
}

func (s *courseService) CreateCourse(course Course) error {
	return s.r.CreateCourse(course)
}

func (s *courseService) GetShortCourses(searchValue string) ([]Course, error) {
	return s.r.GetShortCourse(searchValue)
}

func (s *courseService) GetLongCourses(searchValue string) ([]Course, error) {
	var param newm_helper.Param
	var courses []Course
	fields := make(map[string]WebCourseParam)

	for key, value := range WebCourseParams {
		fields[key] = value
	}

	for key, value := range fields {

		param.Body = nil
		param.Method = "GET"
		param.Headers = map[string]interface{}{
			"Accept": "*/*",
		}
		param.CreateLog = true
		param.RequestId = newm_helper.NewRequestId()

		for {
			param.Url = strings.Replace(value.Url, PAGE, url.QueryEscape(fmt.Sprintf("%d", value.Page)), -1)
			param.Url = strings.Replace(param.Url, SEARCH_VALUE, url.QueryEscape(searchValue), -1)
			value.Page++

			body, err := s.r.GetHtmlCourseInWeb(param)
			if err != nil {
				return nil, newm_helper.Trace(err)
			}

			doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
			if err != nil {
				return nil, newm_helper.Trace(err)
			}

			element := doc.Find(value.MainField)

			element.Each(func(i int, node *goquery.Selection) {

				course := s.findCourseInHtml(node, value.Fields)
				course.Platform = key

				strName := strings.ToLower(course.Name)
				strDescription := strings.ToLower(course.Description)
				strSearchValue := strings.ToLower(searchValue)

				if strings.Contains(strName, strSearchValue) || strings.Contains(strDescription, strSearchValue) {
					courses = append(courses, course)
				}

			})

			if element.Length() == 0 {
				break
			}
		}
	}

	return courses, nil
}

func (s *courseService) findCourseInHtml(node *goquery.Selection, fields map[string]string) Course {
	var course Course

	for key, value := range fields {

		var valuesInHtml []string
		var selector, attr string

		parts := strings.Split(value, "<>")
		if len(parts) == 1 {
			selector = parts[0]
			attr = ""
		} else {
			selector = parts[0]
			attr = parts[1]
		}

		node.Find(selector).Each(func(i int, s *goquery.Selection) {

			if attr != "" {
				str, ok := s.Attr(attr)
				if !ok {
					str = fmt.Sprintf("%s not found", key)
				}

				valuesInHtml = append(valuesInHtml, str)
			} else {
				valuesInHtml = append(valuesInHtml, strings.TrimSpace(strings.ReplaceAll(s.Text(), "\n", "")))
			}
		})

		s.fillCourse(&course, valuesInHtml, key)
	}

	return course
}

func (s *courseService) fillCourse(course *Course, values []string, atribute string) {
	switch atribute {
	case "name":
		course.Name = strings.Join(values, ", ")

	case "description":
		course.Description = strings.Join(values, ", ")

	case "language":
		course.Language = strings.Join(values, ", ")

	case "author":
		course.Author = strings.Join(values, ", ")

	case "duration":
		course.Duration = strings.Join(values, ", ")

	case "rating":
		course.Rating = strings.Join(values, ", ")

	case "money":
		course.Money = strings.Join(values, ", ")

	case "link":
		course.Link = strings.Join(values, ", ")
	}
}
