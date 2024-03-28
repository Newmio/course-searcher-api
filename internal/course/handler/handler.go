package handler

import (
	"fmt"
	"searcher/internal/course/model/dto"
	"searcher/internal/course/service"

	"github.com/Newmio/newm_helper"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	s service.ICourseService
}

func NewHandler(s service.ICourseService) *Handler {
	return &Handler{s: s}
}

func (h *Handler) InitCourseRoutes(e *echo.Echo) {
	api := e.Group("/api")
	{
		course := api.Group("/course")
		{
			course.GET("/get/long", h.GetLongCourses)
			course.POST("/create", h.CreateCourse)
			course.PUT("/update", h.UpdateCourse)
		}
	}
}

func (h *Handler) UpdateCourse(c echo.Context) error {
	var course dto.PutUpdateCourseRequest

	if err := c.Bind(&course); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if err := h.s.UpdateCourse(course); err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return c.JSON(200, nil)
}

func (h *Handler) CreateCourse(c echo.Context) error {
	var course dto.CreateCourseRequest

	if err := c.Bind(&course); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if err := course.Validate(); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
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
		return c.XML(200, courses)

	} else if accept == "text/html" {
		strHtml, err := newm_helper.RenderHtml("static/course/course_template.html", courses)
		if err != nil {
			fmt.Println(err)
			return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
		}
		return c.HTML(200, strHtml)
	}

	return c.JSON(200, courses)
}
