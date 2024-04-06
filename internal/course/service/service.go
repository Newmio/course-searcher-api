package service

import (
	"fmt"
	"searcher/internal/course/model/dto"
	"searcher/internal/course/model/entity"
	"searcher/internal/course/repository"
	"strings"

	"github.com/Newmio/newm_helper"
	"github.com/PuerkitoBio/goquery"
)

type WebCourseParam struct {
	Url        string
	MainField  string
	Pagination bool
	Fields     map[string]string
}

const (
	SEARCH_VALUE = "{{searchValue}}"
	PAGE         = "{{page}}"
)

var WebCourseParams = map[string]WebCourseParam{
	"CourseHunter": {
		Url:       fmt.Sprintf("https://coursehunter.net/search?q=%s&order_by=votes_pos&order=desc&searching=true&page=%s", SEARCH_VALUE, PAGE),
		MainField: "article.course",
		Fields: map[string]string{
			"name":        "h3.course-primary-name",
			"description": "div.course-description",
			"language":    "div.course-lang",
			"author":      "div.course-lessons a",
			"duration":    "div.course-duration",
			"rating":      "div.course-rating-on<>data-text",
			"money":       "div.course-status",
			"link":        "div.course-details-bottom a<>href",
		},
	},
}

type ICourseService interface {
	GetLongCourses(searchValue dto.GetCourseRequest) (dto.CourseListResponse, error)
	CreateCourse(course dto.CreateCourseRequest) error
	UpdateCourse(course dto.PutUpdateCourseRequest) error
	UpdateCourseByParam(course dto.PutUpdateCourseRequest) error
	GetShortCourses(searchValue dto.GetCourseRequest) (dto.CourseListResponse, error)
}

type courseService struct {
	r repository.ICourseRepo
}

func NewCourseService(r repository.ICourseRepo) ICourseService {
	return &courseService{r: r}
}

func (s *courseService) UpdateCourseByParam(course dto.PutUpdateCourseRequest) error {
	return s.r.UpdateCourseByParam(entity.NewUpdateCourse(course))
}

func (s *courseService) UpdateCourse(course dto.PutUpdateCourseRequest) error {
	return s.r.UpdateCourse(entity.NewUpdateCourse(course))
}

func (s *courseService) CreateCourse(course dto.CreateCourseRequest) error {
	return s.r.CreateCourse(entity.NewCreateCourse(course))
}

func (s *courseService) GetShortCourses(searchValue dto.GetCourseRequest) (dto.CourseListResponse, error) {
	courses, err := s.r.GetShortCourses(searchValue.SearchValue)
	if err != nil {
		return dto.CourseListResponse{}, newm_helper.Trace(err)
	}

	return entity.NewCourseListResponse(courses), nil
}

func (s *courseService) GetLongCourses(searchValue dto.GetCourseRequest) (dto.CourseListResponse, error) {
	var param newm_helper.Param
	var courses []dto.CourseList

	cacheCourses, err := s.r.GetCacheCourses(fmt.Sprintf("%s_%d", searchValue.SearchValue, searchValue.Page))
	if err != nil {
		return dto.CourseListResponse{}, newm_helper.Trace(err)
	}

	if len(cacheCourses) > 0 {
		return entity.NewCourseListResponse(cacheCourses), nil
	}

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
		param.Url = strings.Replace(value.Url, PAGE, fmt.Sprintf("%d", searchValue.Page), -1)
		param.Url = strings.Replace(param.Url, SEARCH_VALUE, searchValue.SearchValue, -1)

		body, err := s.r.GetHtmlCourseInWeb(param)
		if err != nil {
			return dto.CourseListResponse{}, newm_helper.Trace(err)
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
		if err != nil {
			return dto.CourseListResponse{}, newm_helper.Trace(err)
		}

		element := doc.Find(value.MainField)

		element.Each(func(i int, node *goquery.Selection) {

			course := s.findCourseInHtml(node, value.Fields)
			course.Platform = key

			strName := strings.Replace(strings.ToLower(course.Name), " ", "", -1)
			strDescription := strings.Replace(strings.ToLower(course.Description), " ", "", -1)
			strSearchValue := strings.Replace(strings.ToLower(searchValue.SearchValue), " ", "", -1)

			if strings.Contains(strName, strSearchValue) || strings.Contains(strDescription, strSearchValue) {
				courses = append(courses, course)
			}
		})

	}

	if err := s.r.CreateCacheCourses(entity.NewCreateCourses(courses), fmt.Sprintf("%s_%d", searchValue.SearchValue, searchValue.Page)); err != nil {
		return dto.CourseListResponse{}, newm_helper.Trace(err)
	}

	return dto.NewCourseListResponse(courses), nil
}

func (s *courseService) findCourseInHtml(node *goquery.Selection, fields map[string]string) dto.CourseList {
	var course dto.CourseList

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

func (s *courseService) fillCourse(course *dto.CourseList, values []string, atribute string) {
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
