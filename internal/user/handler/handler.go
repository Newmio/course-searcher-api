package handler

import (
	"fmt"
	"net/http"
	"searcher/internal/user/model/dto"
	"searcher/internal/user/service"
	"strconv"

	"github.com/Newmio/newm_helper"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	s service.IUserService
}

func NewHandler(s service.IUserService) *Handler {
	return &Handler{s: s}
}

func (h *Handler) InitUserRoutes(e *echo.Echo, middlewares map[string]echo.MiddlewareFunc) {

	e.POST("/register", h.Register)
	e.POST("/login", h.Login) // TODO: add refresh token

	api := e.Group("/api", middlewares["api"])
	{
		user := api.Group("/user")
		{
			update := user.Group("/update")
			{
				update.PUT("", h.UpdateUser)
				update.PATCH("/password", h.UpdatePassword)
				update.PUT("/info", h.UpdateUserInfo)
				update.PUT("/avatar", h.UpdateUserAvatar)
			}
		}
	}
}

func (h *Handler) UpdateUserAvatar(c echo.Context) error {
	var user dto.UpdateUserAvatarRequest

	if err := c.Bind(&user); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if c.Get("role").(string) != "admin" {
		user.Id = strconv.Itoa(c.Get("userId").(int))
	}

	if err := h.s.UpdateUserAvatar(user); err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	c.Response().Header().Set("HX-Redirect", "/profile")

	return c.JSON(200, nil)
}

func (h *Handler) UpdateUserInfo(c echo.Context) error {
	var info dto.CreatUserInfoRequest

	if err := c.Bind(&info); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if c.Get("role").(string) != "admin" {
		info.IdUser = strconv.Itoa(c.Get("userId").(int))
	}

	if err := h.s.UpdateUserInfo(info); err != nil {
		fmt.Println(err)
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	c.Response().Header().Set("HX-Redirect", "/profile")

	return c.JSON(200, nil)
}

func (h *Handler) UpdatePassword(c echo.Context) error {
	var user dto.UpdateUserPasswordRequest

	if err := c.Bind(&user); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if c.Get("role").(string) != "admin" {
		user.Id = c.Get("userId").(int)
	}

	if err := user.Validate(); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if err := h.s.UpdatePassword(user); err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return c.JSON(200, nil)
}

func (h *Handler) UpdateUser(c echo.Context) error {
	var user dto.UpdateUserRequest

	if err := c.Bind(&user); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if c.Get("role").(string) != "admin" {
		user.Id = strconv.Itoa(c.Get("userId").(int))
	}

	if err := user.Validate(); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if err := h.s.UpdateUser(user); err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	c.Response().Header().Set("HX-Redirect", "/profile")

	return c.JSON(200, nil)
}

func (h *Handler) Login(c echo.Context) error {
	var user dto.LoginUserRequest

	cookie, err := strconv.ParseBool(c.QueryParam("cookie"))
	if err != nil {
		cookie = false
	}

	if err := c.Bind(&user); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if err := user.Validate(); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	tokens, err := h.s.Login(user)
	if err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	if cookie {
		c.SetCookie(&http.Cookie{
			Name:     "access",
			Value:    "Bearer " + tokens.Access,
			Path:     "/",
			MaxAge:   tokens.Exp,
			Secure:   false,
			HttpOnly: true,
		})
		c.SetCookie(&http.Cookie{
			Name:     "refresh",
			Value:    "Bearer " + tokens.Refresh,
			Path:     "/",
			MaxAge:   tokens.ExpRefresh,
			Secure:   false,
			HttpOnly: true,
		})
	}

	c.Response().Header().Set("HX-Redirect", "/")

	return c.JSON(200, tokens)
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

	c.Response().Header().Set("HX-Redirect", "/login_form")

	return c.JSON(201, nil)
}
