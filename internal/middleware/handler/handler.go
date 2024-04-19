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

// func (h *Handler) ParseToken(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		var token []string
// 		header := c.Request().Header.Get("Authorization")

// 		if header == "" {

// 			cookie, err := c.Cookie("access")
// 			if err != nil {
// 				c.Response().Header().Set("HX-Redirect", "http://localhost:8088/login_form")
// 				c.Response().Header().Set("Cache-Control", "no-cache, no-store")
// 				return c.JSON(401, newm_helper.ErrorResponse("unauthorized"))
// 			}

// 			token = strings.Split(cookie.Value, " ")
// 		} else {
// 			token = strings.Split(header, " ")
// 		}

// 		if len(token) != 2 {
// 			c.Response().Header().Set("HX-Redirect", "http://localhost:8088/login_form")
// 			c.Response().Header().Set("Cache-Control", "no-cache, no-store")
// 			return c.JSON(401, newm_helper.ErrorResponse("unauthorized"))
// 		}

// 		if token[0] != "Bearer" {
// 			c.Response().Header().Set("HX-Redirect", "http://localhost:8088/login_form")
// 			c.Response().Header().Set("Cache-Control", "no-cache, no-store")
// 			return c.JSON(401, newm_helper.ErrorResponse("unauthorized"))
// 		}

// 		id, role, err := h.s.ParseToken(token[1])
// 		if err != nil {
// 			c.Response().Header().Set("HX-Redirect", "http://localhost:8088/login_form")
// 			c.Response().Header().Set("Cache-Control", "no-cache, no-store")
// 			return c.JSON(401, newm_helper.ErrorResponse("unauthorized"))
// 		}

// 		c.Set("userId", id)
// 		c.Set("role", role)

// 		return next(c)
// 	}
// }

func AuthError(c echo.Context) error {
	if c.Request().Header.Get("User-Agent") != "" {
		c.Redirect(301, "/login_form")
	}

	return c.JSON(401, newm_helper.ErrorResponse("unauthorized"))
}

func (h *Handler) ParseToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var token []string
		header := c.Request().Header.Get("Authorization")

		if header == "" {

			cookie, err := c.Cookie("access")
			if err != nil {
				return AuthError(c)
			}

			token = strings.Split(cookie.Value, " ")
		} else {
			token = strings.Split(header, " ")
		}

		if len(token) != 2 {
			return AuthError(c)
		}

		if token[0] != "Bearer" {
			return AuthError(c)
		}

		id, role, err := h.s.ParseToken(token[1])
		if err != nil {
			return AuthError(c)
		}

		c.Set("userId", id)
		c.Set("role", role)

		return next(c)
	}
}
