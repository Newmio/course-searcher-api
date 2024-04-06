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

func (h *Handler) InitCourseRoutes(e *echo.Echo, middlewares map[string]echo.MiddlewareFunc) {
	api := e.Group("/api", middlewares["api"])
	api.Use()
	{
		course := api.Group("/course")
		{
			get := course.Group("/get")
			{
				get.POST("/long", h.GetLongCourses)
				get.POST("/short", h.GetShortCourses)
			}

			course.POST("/create", h.CreateCourse)
			course.PUT("/update", h.UpdateCourse)
			course.PATCH("/update_by_param", h.UpdateCourseByParam)
		}
	}
}

func (h *Handler) GetShortCourses(c echo.Context) error {
	var course dto.GetCourseRequest
	accept := c.Request().Header.Get("Accept")

	if err := c.Bind(&course); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if accept == "" {
		c.JSON(400, newm_helper.ErrorResponse("Accept header is required"))
	}

	if err := course.Validate(); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	courses, err := h.s.GetShortCourses(course)
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

func (h *Handler) UpdateCourseByParam(c echo.Context) error {
	var course dto.PutUpdateCourseRequest

	if err := c.Bind(&course); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if err := h.s.UpdateCourseByParam(course); err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return c.JSON(200, nil)
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

	userId := c.Get("userId").(string)

	if err := c.Bind(&course); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if err := course.Validate(); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	err := h.s.CreateCourse(course, userId)
	if err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return c.JSON(201, nil)
}

func (h *Handler) GetLongCourses(c echo.Context) error {
	var course dto.GetCourseRequest
	accept := c.Request().Header.Get("Accept")

	if err := c.Bind(&course); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if accept == "" {
		c.JSON(400, newm_helper.ErrorResponse("Accept header is required"))
	}

	if err := course.Validate(); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	courses, err := h.s.GetLongCourses(course)
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
