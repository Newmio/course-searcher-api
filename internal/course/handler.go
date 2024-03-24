package course

import (
	"fmt"

	"github.com/Newmio/newm_helper"
	"github.com/labstack/echo/v4"
)

type ICourseService interface {
	GetLongCourses(searchValue string) ([]Course, error)
	CreateCourse(course Course) error
	UpdateCourse(course Course) error
}

type Handler struct {
	s ICourseService
}

func NewHandler(s ICourseService) *Handler {
	return &Handler{s: s}
}

func (h *Handler) InitCourseRoutes(e *echo.Echo) *echo.Echo {

	api := e.Group("/api")
	{
		//api.GET("/short_courses", h.GetShortCourse)
		api.GET("/long_courses", h.GetLongCourses)
	}

	return e
}

func (h *Handler) CreateCourse(c echo.Context) error {
	var course Course

	if err := c.Bind(&course); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if course.Link != "" {
		return c.JSON(400, newm_helper.ErrorResponse("bad request"))
	}

	err := h.s.CreateCourse(course)
	if err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return c.JSON(201, nil)
}

func (h *Handler) GetLongCourses(c echo.Context) error {
	searchValue := c.QueryParam("search_value")
	accept := c.Request().Header.Get("Accept")

	if accept == "" {
		c.JSON(400, newm_helper.ErrorResponse("Accept header is required"))
	}

	if searchValue == "" {
		return c.JSON(400, newm_helper.ErrorResponse("searchValue is required"))
	}

	courses, err := h.s.GetLongCourses(searchValue)
	if err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	if accept == "application/xml" {
		return c.XML(200, map[string]interface{}{
			"courses": courses,
			"count":   len(courses),
		})

	} else if accept == "text/html" {
		strHtml, err := newm_helper.RenderHtml("static/course/course_template.html", courses)
		if err != nil {
			fmt.Println(err)
			return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
		}
		return c.HTML(200, strHtml)
	}

	return c.JSON(200, map[string]interface{}{
		"courses": courses,
		"count":   len(courses),
	})
}
