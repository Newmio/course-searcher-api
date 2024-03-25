package user

import (
	"github.com/Newmio/newm_helper"
	"github.com/labstack/echo/v4"
)

type IUserService interface {
	CreateUser(user User) error
}

type Handler struct {
	s IUserService
}

func NewHandler(s IUserService) *Handler {
	return &Handler{s: s}
}

func (h *Handler) InitUserRoutes(e *echo.Echo) *echo.Echo {

	e.POST("/register", h.Register)

	return e
}

func (h *Handler) Register(c echo.Context) error {
	var user User

	if err := c.Bind(&user); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if user.Login == "" || user.Password == ""{
		return c.JSON(400, newm_helper.ErrorResponse("bad request"))
	}

	if err := h.s.CreateUser(user); err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return c.JSON(201, nil)
}
