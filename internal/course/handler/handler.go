package handler

import (
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
				get.GET("/by_user", h.GetCoursesByUser)
				get.GET("/history", h.GetCoursesHistory)
			}

			course.POST("/create", h.CreateCourse)
			course.PUT("/update", h.UpdateCourse)
			course.PATCH("/update_by_param", h.UpdateCourseByParam)
			course.GET("/check", h.CheckCourse)
		}
	}
	e.GET("/test", h.T)
}

func (h *Handler) T(c echo.Context) error {
	if err := h.s.CreateCourseEvent([]byte("test")); err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}
	return c.JSON(200, nil)
}

func (h *Handler) GetCoursesHistory(c echo.Context) error {
	courses, err := h.s.GetCacheCoursesByUser(c.Get("userId").(int))
	if err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return courseResponse(c, courses)
}

func (h *Handler) CheckCourse(c echo.Context) error {
	accept := c.Request().Header.Get("Accept")

	flag, err := h.s.CheckCourse(c.QueryParam("link"))
	if err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	resp := map[string]bool{"check": flag}

	switch accept {
	case "application/xml":
		return c.XML(200, resp)

	case "text/html":
		strHtml := `<h5 style="color: red;">Запрос на подтверждение курса<br>уже отправлен на проверку</h5>`

		if flag {
			strHtml = `<h5 style="color: green;">Курс уже доступен!<br>Найдите его у себя в профиле</h5>`
		}

		return c.HTML(200, strHtml)

	default:
		return c.JSON(200, resp)
	}
}

func (h *Handler) GetCoursesByUser(c echo.Context) error {
	id := c.Get("userId").(int)

	courses, err := h.s.GetCoursesByUser(id)
	if err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return c.JSON(200, courses)
}

func (h *Handler) GetShortCourses(c echo.Context) error {
	var course dto.GetCourseRequest

	if err := c.Bind(&course); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if err := course.Validate(); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	courses, err := h.s.GetShortCourses(course)
	if err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return courseResponse(c, courses)
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
	var course dto.GetCourseRequest

	if err := c.Bind(&course); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if err := course.Validate(); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	courses, err := h.s.GetLongCourses(course, c.Get("userId").(int))
	if err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return courseResponse(c, courses)
}

func courseResponse(c echo.Context, courses dto.CourseListResponse) error {
	switch c.Request().Header.Get("Accept") {
	case "application/xml":
		return c.XML(200, courses)

	case "text/html":
		strHtml, err := newm_helper.RenderHtml("template/course/course_template.html", courses)
		if err != nil {
			return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
		}
		return c.HTML(200, strHtml)

	default:
		return c.JSON(200, courses)
	}
}
