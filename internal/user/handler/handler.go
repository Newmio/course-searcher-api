package handler

import (
	"searcher/internal/user/model/dto"
	"searcher/internal/user/service"

	"github.com/Newmio/newm_helper"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	s service.IUserService
}

func NewHandler(s service.IUserService) *Handler {
	return &Handler{s: s}
}

func (h *Handler) InitUserRoutes(e *echo.Echo) {

	e.POST("/register", h.Register)

}

func (h *Handler) Register(c echo.Context) error {
	var user dto.RegisterUserRequest

	if err := c.Bind(&user); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if err := user.Validate(); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if err := h.s.CreateUser(user); err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return c.JSON(201, nil)
}
