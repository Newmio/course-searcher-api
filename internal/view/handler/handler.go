package handler

import (
	"searcher/internal/view/service"

	"github.com/Newmio/newm_helper"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	s service.IViewService
}

func NewHandler(s service.IViewService) *Handler {
	return &Handler{s: s}
}

func (h *Handler) InitViewRoutes(e *echo.Echo, middlewares map[string]echo.MiddlewareFunc) {
	e.GET("/", h.Profile, middlewares["api"])

	e.GET("/login_form", func(c echo.Context) error {
		return c.File("template/user/login/login.html")
	})
}

func (h *Handler) Profile(c echo.Context) error {
	id := c.Get("userId").(int)

	html, err := h.s.GetUserProfile(id)
	if err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return c.HTML(200, html)
}
