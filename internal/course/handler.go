package course

import "github.com/labstack/echo/v4"

type Handler struct {
	s ICourseService
}

func NewHandler(s ICourseService) *Handler {
	return &Handler{s: s}
}

func (h *Handler) InitCourseRoutes(e *echo.Echo) *echo.Echo {
	return e
}

func (h *Handler) Test(c echo.Context) error {
	return c.String(200, "Hello, World!")
}
