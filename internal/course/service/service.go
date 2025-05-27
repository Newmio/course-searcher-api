package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"searcher/internal/course/model/dto"
	"searcher/internal/course/model/entity"
	"searcher/internal/course/repository"
	"strconv"
	"strings"

	repoUser "searcher/internal/user/repository"

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
		MainField: "article.post-card",
		Fields: map[string]string{
			"name":        "h3.post-title",
			"description": "p.post-description",
			"language":    "div.post-footer.post-tags.post-date",
			"author":      "a.post-avatar-box span",
			"duration":    "div.post-duration",
			"rating":      "div.post-rating.post-rating-plus",
			"money":       "div.post-status",
			"link":        "picture img <>src",
			"icon-link":   "picture img <>src",
		},
	},
}

type ICourseService interface {
	GetLongCourses(searchValue dto.GetCourseRequest, userId int) (dto.CourseListResponse, error)
	CreateCourse(course dto.CreateCourseRequest) error
	UpdateCourse(course dto.PutUpdateCourseRequest) error
	UpdateCourseByParam(course dto.PutUpdateCourseRequest) error
	GetShortCourses(searchValue dto.GetCourseRequest) (dto.CourseListResponse, error)
	GetCoursesByUser(id int) (map[string]dto.CourseListResponse, error)
	CheckCourse(userId int, link string) (bool, error)
	GetCacheCoursesByUser(userId int) (dto.CourseListResponse, error)
	CreateCourseEvent(value []byte) error
	AppendEventOffset(offset int) error
	CheckExistsEventOffset(offset int) (bool, error)
	GetCacheCheckCourses() (dto.CourseListResponse, error)
	CreateApproveCourse(link string) error
	GetWaitingCheckById(userId int) (dto.CourseListResponse, error)
	SetCourseCoins(courseId, studentId, userId int, coins map[string]interface{}) error
	GetCoursesCheckStop(userId int, token string) ([]CourseCheckStop, error)
	// это худший проект в моей жизни
	// написанный на коленке без соблюдения любых архитектурных принципов
	// в этом коде нету логики
}

type courseService struct {
	r     repository.ICourseRepo
	rUser repoUser.IUserRepo
}

func NewCourseService(r repository.ICourseRepo, rUser repoUser.IUserRepo) ICourseService {
	return &courseService{r: r, rUser: rUser}
}

type CourseCheckStop struct {
	Link     string
	IconLink string
	Name     string
	DocLink  string
	Platform string
	Author   string
}

func (s *courseService) GetCoursesCheckStop(userId int, token string) ([]CourseCheckStop, error) {
	var resp []CourseCheckStop

	courses, err := s.r.GetCoursesStopCheck(userId)
	if err != nil {
		return nil, newm_helper.Trace(err)
	}

	for _, v := range courses {
		var c CourseCheckStop

		c.Link = v.Link
		c.IconLink = v.IconLink
		c.Name = v.Name
		c.Platform = v.Platform
		c.Author = v.Author

		req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8088/api/course/genreport?course_id=%d&student_id=%d", v.Id, userId), nil)
		if err != nil {
			return nil, newm_helper.Trace(err)
		}

		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}

		reqResp, err := client.Do(req)
		if err != nil {
			return nil, newm_helper.Trace(err)
		}
		defer reqResp.Body.Close()

		body, err := io.ReadAll(reqResp.Body)
		if err != nil {
			return nil, newm_helper.Trace(err)
		}

		mapResp := make(map[string]interface{})

		err = json.Unmarshal(body, &mapResp)
		if err != nil {
			return nil, newm_helper.Trace(err)
		}

		c.DocLink = mapResp["link"].(string)

		resp = append(resp, c)
	}

	return resp, nil
}

func (s *courseService) SetCourseCoins(courseId, studentId, userId int, coins map[string]interface{}) error {
	userInfo, err := s.rUser.GetUserInfo(userId)
	if err != nil {
		return newm_helper.Trace(err)
	}

	if userInfo.Proffession != "cms" {
		return nil
	}

	var stopCheck bool

	cmsUsers, err := s.rUser.GetCMSUsers()
	if err != nil {
		return newm_helper.Trace(err)
	}

	coinsMap, err := s.r.GetCourseUserCoins(courseId, studentId)
	if err != nil {
		return newm_helper.Trace(err)
	}

	if _, ok := coinsMap["res"]; !ok {
		if len(coinsMap) >= len(cmsUsers)-1 {
			stopCheck = true
		}
	}

	c, err := strconv.Atoi(coins["coins"].(string))
	if err != nil {
		return newm_helper.Trace(err)
	}

	var credits int

	if v, ok := coins["credits"]; ok{
		credits, err = strconv.Atoi(v.(string))
		if err != nil {
			return newm_helper.Trace(err)
		}
	}

	var educName string

	if v, ok := coins["educ_name"]; ok{
		educName = v.(string)
	}

	err = s.r.SetCourseCoins(courseId, studentId, userId, c, stopCheck, credits, educName)
	if err != nil {
		return newm_helper.Trace(err)
	}

	err = s.r.SumCoins(studentId, courseId)
	if err != nil {
		return newm_helper.Trace(err)
	}

	return nil
}

func (s *courseService) GetWaitingCheckById(userId int) (dto.CourseListResponse, error) {
	links, err := s.r.GetWaitingCheckById(userId)
	if err != nil {
		return dto.CourseListResponse{}, newm_helper.Trace(err)
	}

	var courses []entity.CourseList

	for _, v := range links {
		course, err := s.r.GetCacheCourseByLink(v)
		if err != nil {
			return dto.CourseListResponse{}, newm_helper.Trace(err)
		}

		courses = append(courses, course)
	}

	return entity.NewCourseListResponse(courses), nil
}

func (s *courseService) CreateApproveCourse(link string) error {
	course, err := s.r.GetCacheCourseByLink(link)
	if err != nil {
		return newm_helper.Trace(err)
	}

	createCourses := make([]entity.CourseList, 0)
	createCourses = append(createCourses, course)

	err = s.r.CreateCourse(entity.NewCreateCourses(entity.NewCourseList(createCourses))[0])
	if err != nil {
		return newm_helper.Trace(err)
	}

	err = s.r.DeleteCacheCheckCourses(link)
	if err != nil {
		return newm_helper.Trace(err)
	}

	userIds, err := s.r.GetWaitingCheck(link)
	if err != nil {
		return newm_helper.Trace(err)
	}

	course, err = s.r.GetCourseByLink(link)
	if err != nil {
		return newm_helper.Trace(err)
	}

	for _, v := range userIds {
		err = s.r.CreateCourseUser(map[string]interface{}{"id_user": v, "id_course": course.Id})
		if err != nil {
			return newm_helper.Trace(err)
		}

		err = s.r.DeleteWaitingCheck(link)
		if err != nil {
			return newm_helper.Trace(err)
		}
	}

	return nil
}

func (s *courseService) GetCacheCheckCourses() (dto.CourseListResponse, error) {
	courses, err := s.r.GetCacheCheckCourses()
	if err != nil {
		return dto.CourseListResponse{}, newm_helper.Trace(err)
	}

	return entity.NewCourseListResponse(courses), nil
}

func (s *courseService) CheckExistsEventOffset(offset int) (bool, error) {
	return s.r.CheckExistsEventOffset(offset)
}

func (s *courseService) AppendEventOffset(offset int) error {
	return s.r.AppendEventOffset(offset)
}

func (s *courseService) CreateCourseEvent(value []byte) error {
	return s.r.CreateCourseEvent(value)
}

func (s *courseService) GetCacheCoursesByUser(userId int) (dto.CourseListResponse, error) {
	cacheCourses, err := s.r.GetCacheCourses(fmt.Sprintf("courses_for_user_%d", userId))
	if err != nil {
		return dto.CourseListResponse{}, newm_helper.Trace(err)
	}

	return entity.NewCourseListResponse(cacheCourses), nil
}

func (s *courseService) CheckCourse(userId int, link string) (bool, error) {
	course, err := s.r.GetCourseByLink(link)
	if err != nil {
		return false, newm_helper.Trace(err)
	}

	if course.Link == "" {
		course, err := s.r.GetCacheCourseByLink(link)
		if err != nil {
			return false, newm_helper.Trace(err)
		}

		if err := s.r.CreateCacheCheckCourses(course); err != nil {
			return false, newm_helper.Trace(err)
		}

		if err := s.r.CreateWaitingCheck(userId, course.Link); err != nil {
			return false, newm_helper.Trace(err)
		}

		return false, nil
	}

	return true, nil
}

func (s *courseService) GetCoursesByUser(id int) (map[string]dto.CourseListResponse, error) {
	resp := make(map[string]dto.CourseListResponse)

	courses, err := s.r.GetCoursesByUser(id)
	if err != nil {
		return nil, newm_helper.Trace(err)
	}

	for key, value := range courses {
		resp[key] = entity.NewCourseListResponse(value)
	}

	return resp, nil
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

func (s *courseService) GetLongCourses(searchValue dto.GetCourseRequest, userId int) (dto.CourseListResponse, error) {
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

	entityCourses := entity.NewCreateCourses(courses)

	if err := s.r.CreateCacheCourses(entityCourses, fmt.Sprintf("%s_%d", searchValue.SearchValue, searchValue.Page)); err != nil {
		return dto.CourseListResponse{}, newm_helper.Trace(err)
	}

	if err := s.r.CreateCacheCourses(entityCourses, fmt.Sprintf("courses_for_user_%d", userId)); err != nil {
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

	case "icon-link":
		course.IconLink = strings.Join(values, ", ")
	}
}
