package course

import (
	"github.com/Newmio/newm_helper"
	"github.com/labstack/echo/v4"
)

type ICourseService interface {
	GetLongCourses(searchValue string, inDescription bool) ([]Course, error)
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
		//api.GET("/short_course", h.GetShortCourse)
		api.GET("/long_course", h.GetLongCourse)
	}

	return e
}

func (h *Handler) GetLongCourse(c echo.Context) error {
	
	if c.QueryParam("searchValue") == "" {
		return c.JSON(400, newm_helper.ErrorResponse("searchValue is required"))
	}
}
