package course

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Newmio/newm_helper"
	"github.com/PuerkitoBio/goquery"
)

type ICourseRepo interface {
	GetShortCourse(valueSearch string, inDescription bool) ([]Course, error)
	GetHtmlCourseInWeb(param newm_helper.Param) ([]byte, error)
}

type courseService struct {
	r ICourseRepo
}

func NewCourseService(r ICourseRepo) ICourseService {
	return &courseService{r: r}
}

func (s *courseService) GetShortCourses(searchValue string, inDescription bool) ([]Course, error) {
	return s.r.GetShortCourse(searchValue, inDescription)
}

func (s *courseService) GetLongCourses(searchValue string, inDescription bool) ([]Course, error) {
	var param newm_helper.Param
	var coueses []Course
	fields := make(map[string]WebCourseParam)

	for key, value := range WebCourseParams {
		fields[key] = value
	}

	for _, value := range fields {

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

			var names, descriptions, languages, authors, durations, moneys, links []string
			element := doc.Find(value.MainField)

			for key2, value2 := range value.Fields {

				element.Find(value2).Each(func(i int, s *goquery.Selection) {

					switch key2{
					case "name":
						names = append(names, s.Text())

					case "description":
						descriptions = append(descriptions, s.Text())

					case "language":
						languages = append(languages, s.Text())

					case "author":
						authors = append(authors, s.Text())

					case "duration":
						durations = append(durations, s.Text())

					case "money":
						moneys = append(moneys, s.Text())

					case "link":
						link, ok := s.Attr("href")
						if !ok {
							link = "link not found"
						}

						links = append(links, link)
					}
				})
			}

			for i := 0; i < len(names); i++ {
				coueses = append(coueses, Course{
					Name:        names[i],
					Description: descriptions[i],
					Language:    languages[i],
					Author:      authors[i],
					Duration:    durations[i],
					Money:       moneys[i],
					Link:        links[i],
				})
			}

			if names == nil {
				break
			}
		}
	}

	return coueses, nil
}
