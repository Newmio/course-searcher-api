package user

import (
	"github.com/labstack/echo/v4"
)

type Handler struct {
	s IUserService
}

func NewHandler(s IUserService) *Handler {
	return &Handler{s: s}
}

func (h *Handler) InitUserRoutes(e *echo.Echo) *echo.Echo {
	return e
}
