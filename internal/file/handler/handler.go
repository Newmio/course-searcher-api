package handler

import (
	"io"
	"searcher/internal/file/model/dto"
	"searcher/internal/file/service"

	"github.com/Newmio/newm_helper"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	s service.IFileService
}

func NewHandler(s service.IFileService) *Handler {
	return &Handler{s: s}
}

func (h *Handler) InitFileRoutes(e *echo.Echo, middlewares map[string]echo.MiddlewareFunc) {

}

func (h *Handler) UploadEducationFile(c echo.Context) error {
	var file dto.CreateFileRequest

	file.UserId = c.Get("userId").(int)
	file.FileType = c.QueryParam("fileType")

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}
	file.FileBytes = body

	if err := file.Validate(); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	return c.JSON(201, nil)
}
