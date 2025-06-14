package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"searcher/internal/file/model/dto"
	"searcher/internal/file/service"
	"strconv"
	"strings"

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
	api := e.Group("/api", middlewares["api"])
	{
		course := api.Group("/course")
		{
			get := course.Group("/get")
			{
				get.GET("/foreport", h.GetCoursesForReport)
			}

			course.GET("/genreport", h.GenerateReport)
		}

		file := api.Group("/file")
		{
			upload := file.Group("/upload")
			{
				upload.POST("/report", h.UploadReportFile)
				upload.POST("/education", h.UploadEducationFile)
			}

			get := file.Group("/get")
			{
				report := get.Group("/report")
				{
					report.GET("/info", h.GetReportFilesInfoByUserId)
					report.GET("/by_id", h.GetReportFileById)
				}

				education := get.Group("/education")
				{
					education.GET("/info", h.GetEducationFilesInfoByUserId)
					education.GET("/by_id", h.GetEducationFileById)
				}
			}

			delete := file.Group("/delete")
			{
				delete.DELETE("/report", h.DeleteReportFile)
				delete.DELETE("/education", h.DeleteEducationFile)
			}
		}
	}

	e.GET("/test", h.Test)
	// e.GET("/test2", h.Test2)
}

type CoursesForReport struct{
	Courses []service.CourseForReport
}

func (h *Handler) GenerateReport(c echo.Context) error {
	studentId, err := strconv.Atoi(c.QueryParam("student_id"))
	if err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	courseId, err := strconv.Atoi(c.QueryParam("course_id"))
	if err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	link, err := h.s.GetReport(studentId, courseId)
	if err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return c.JSON(200, map[string]string{"link": link})
}

func (h *Handler) GetCoursesForReport(c echo.Context) error {
	resp, err := h.s.GetCoursesForReport()
	if err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	courses := CoursesForReport{
		Courses: resp,
	}

	strHtml, err := newm_helper.RenderHtml("template/messages/course_reporttemplate.html", courses)
	if err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return c.HTML(200, strHtml)
}

func (h *Handler) Test2(c echo.Context) error {
	file, err := os.Open("images/image1.png")
	if err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return c.Blob(200, http.DetectContentType(b), b)
}

func (h *Handler) Test(c echo.Context) error {
	if err := h.s.TestPdf(); err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return c.JSON(200, nil)
}

func (h *Handler) DeleteEducationFile(c echo.Context) error {
	fileId, err := strconv.Atoi(c.QueryParam("file_id"))
	if err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if c.Get("role").(string) != "admin" {
		return c.JSON(200, nil)
	}

	if err := h.s.DeleteEducationFile(fileId); err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return c.JSON(200, nil)
}

func (h *Handler) DeleteReportFile(c echo.Context) error {
	fileId, err := strconv.Atoi(c.QueryParam("file_id"))
	if err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if c.Get("role").(string) != "admin" {
		return c.JSON(200, nil)
	}

	if err := h.s.DeleteReportFile(fileId); err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return c.JSON(200, nil)
}

func (h *Handler) GetEducationFileById(c echo.Context) error {
	fileId, err := strconv.Atoi(c.QueryParam("file_id"))
	if err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	getUserId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if c.Get("role").(string) != "admin" || getUserId != c.Get("userId").(int) {
		return c.JSON(200, nil)
	}

	file, err := h.s.GetEducationFileById(fileId)
	if err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return c.JSON(200, file)
}

func (h *Handler) GetReportFileById(c echo.Context) error {
	fileId, err := strconv.Atoi(c.QueryParam("file_id"))
	if err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	getUserId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if c.Get("role").(string) != "admin" || getUserId != c.Get("userId").(int) {
		return c.JSON(200, nil)
	}

	file, err := h.s.GetReportFileById(fileId)
	if err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return c.JSON(200, file)
}

func (h *Handler) GetReportFilesInfoByUserId(c echo.Context) error {
	getUserId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if c.Get("role").(string) != "admin" || getUserId != c.Get("userId").(int) {
		return c.JSON(200, nil)
	}

	files, err := h.s.GetReportFilesInfoByUserId(getUserId)
	if err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return c.JSON(200, files)
}

func (h *Handler) GetEducationFilesInfoByUserId(c echo.Context) error {
	getUserId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if c.Get("role").(string) != "admin" || getUserId != c.Get("userId").(int) {
		return c.JSON(200, nil)
	}

	files, err := h.s.GetEducationFilesInfoByUserId(getUserId)
	if err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return c.JSON(200, files)
}

func (h *Handler) UploadReportFile(c echo.Context) error {
	var file dto.CreateFileRequest

	if c.Get("role").(string) != "admin" {
		fmt.Println("role", c.Get("role"))
		return c.JSON(201, nil)
	}

	file.UserId = c.Get("userId").(int)
	file.FileType = c.QueryParam("file_type")

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}
	file.FileBytes = body

	if err := file.Validate(); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if err := h.s.CreateReportFile(file, c.QueryParam("courseLink")); err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return c.JSON(201, nil)
}

func (h *Handler) UploadEducationFile(c echo.Context) error {
	var file dto.CreateFileRequest

	file.UserId = c.Get("userId").(int)
	// file.FileType = c.QueryParam("file_type")

	// body, err := io.ReadAll(c.Request().Body)
	// if err != nil {
	// 	return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	// }
	// file.FileBytes = body

	file.FileType = strings.Split(c.FormValue("file_type"), "/")[1] // <- вот так

	// Получаем файл из multipart
	f, err := c.FormFile("file")
	if err != nil {
		return c.JSON(400, newm_helper.ErrorResponse("file not provided"))
	}

	src, err := f.Open()
	if err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}
	defer src.Close()

	body, err := io.ReadAll(src)
	if err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}
	file.FileBytes = body

	if err := file.Validate(); err != nil {
		return c.JSON(400, newm_helper.ErrorResponse(err.Error()))
	}

	if err := h.s.CreateEducationFile(file, c.QueryParam("courseLink")); err != nil {
		return c.JSON(500, newm_helper.ErrorResponse(err.Error()))
	}

	return c.JSON(201, nil)
}
