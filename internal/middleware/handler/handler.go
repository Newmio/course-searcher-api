package handler

import (
	"searcher/internal/middleware/service"
	"strings"

	"github.com/Newmio/newm_helper"
	"github.com/labstack/echo/v4"
)

type IMiddlewareHandler interface {
	ParseToken(next echo.HandlerFunc) echo.HandlerFunc
}

type Handler struct {
	s service.IMiddlewareService
}

func NewHandler(s service.IMiddlewareService) *Handler {
	return &Handler{s: s}
}

func (h *Handler) ParseToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var token []string
		header := c.Request().Header.Get("Authorization")

		if header == "" {

			cookie, err := c.Cookie("access")
			if err != nil {
				return c.JSON(401, newm_helper.ErrorResponse("token is required"))
			}

			token = strings.Split(cookie.Value, " ")
		}else{
			token = strings.Split(header, " ")
		}

		if len(token) != 2 {
			return c.JSON(401, newm_helper.ErrorResponse("invalid token"))
		}

		if token[0] != "Bearer" {
			return c.JSON(401, newm_helper.ErrorResponse("invalid token"))
		}

		id, role, err := h.s.ParseToken(token[1])
		if err != nil {
			return c.JSON(401, map[string]string{"error": "unauthorized"})
		}

		c.Set("userId", id)
		c.Set("role", role)

		return next(c)
	}
}
