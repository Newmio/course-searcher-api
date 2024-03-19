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
}

type courseService struct {
	r ICourseRepo
}

func NewCourseService(r ICourseRepo) ICourseService {
	return &courseService{r: r}
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
			param.Url = strings.Replace(value.Url, "{{page}}", url.QueryEscape(fmt.Sprintf("%d", value.Page)), -1)
			param.Url = strings.Replace(param.Url, "{{searchValue}}", url.QueryEscape(searchValue), -1)
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

			element.Each(func(i int, s *goquery.Selection) {
				var course Course

				for key2, value2 := range value.Fields {
					var valuesInHtml []string
					var selector, attr string

					parts := strings.Split(value2, "<>")
					if len(parts) == 1 {
						selector = parts[0]
						attr = ""
					} else {
						selector = parts[0]
						attr = parts[1]
					}

					s.Find(selector).Each(func(i int, s *goquery.Selection) {

						if attr != "" {
							str, ok := s.Attr(attr)
							if !ok {
								str = fmt.Sprintf("%s not found", key2)
							}

							valuesInHtml = append(valuesInHtml, str)
						} else {
							valuesInHtml = append(valuesInHtml, strings.TrimSpace(strings.ReplaceAll(s.Text(), "\n", "")))
						}
					})

					switch key2 {
					case "name":
						course.Name = strings.Join(valuesInHtml, ", ")
						course.Platform = key

					case "description":
						course.Description = strings.Join(valuesInHtml, ", ")

					case "language":
						course.Language = strings.Join(valuesInHtml, ", ")

					case "author":
						course.Author = strings.Join(valuesInHtml, ", ")

					case "duration":
						course.Duration = strings.Join(valuesInHtml, ", ")

					case "rating":
						course.Rating = strings.Join(valuesInHtml, ", ")

					case "money":
						course.Money = strings.Join(valuesInHtml, ", ")

					case "link":
						course.Link = strings.Join(valuesInHtml, ", ")
					}
				}

				if strings.Contains(strings.ToLower(course.Description), strings.ToLower(searchValue)) {
					courses = append(courses, course)

				} else if strings.Contains(strings.ToLower(course.Description), strings.ToLower(searchValue)) {
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
