package course

import "github.com/labstack/echo/v4"

type ICourseService interface {
	
}

type Handler struct {
	s ICourseService
}

func NewHandler(s ICourseService) *Handler {
	return &Handler{s: s}
}

func (h *Handler) InitCourseRoutes(e *echo.Echo) *echo.Echo {
	course := e.Group("/course")
	{
		api := course.Group("/api")
		{
			api.GET("search_course", h.SearchCourse)
		}
	}

	return e
}

func (h *Handler) SearchCourse(c echo.Context) error {
	return nil
}
