package service

import (
	"fmt"
	"os"
	rCourse "searcher/internal/course/repository"
	"searcher/internal/file/model/dto"
	rFile "searcher/internal/file/repository"
	rUser "searcher/internal/user/repository"

	"github.com/Newmio/newm_helper"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type IFileService interface {
	CreateReportFile(file dto.CreateFileRequest, link string) error
	CreateEducationFile(file dto.CreateFileRequest, link string) error
	GetReportFilesInfoByUserId(userId int) (dto.GetFileResponse, error)
	GetEducationFilesInfoByUserId(userId int) (dto.GetFileResponse, error)
	GetReportFileById(fileId int) ([]byte, error)
	GetEducationFileById(fileId int) ([]byte, error)
	DeleteReportFile(fileId int) error
	DeleteEducationFile(fileId int) error
	GetCoursesForReport() ([]CourseForReport, error)

	TestPdf() error
}

type fileService struct {
	rFile   rFile.IFileRepo
	rUser   rUser.IUserRepo
	rCourse rCourse.ICourseRepo
}

func NewFileService(rFile rFile.IFileRepo, rUser rUser.IUserRepo, rCourse rCourse.ICourseRepo) IFileService {
	return &fileService{rFile: rFile, rUser: rUser, rCourse: rCourse}
}

type CourseForReport struct {
	Link        string
	Name        string
	IconLink string
	Platform    string
	Author      string
	FileLinks   []string
	StudentInfo string
	UserId      int
}

func (s *fileService) GetCoursesForReport() ([]CourseForReport, error) {
	var resp []CourseForReport

	coursesMap, err := s.rCourse.GetCoursesForReport()
	if err != nil {
		return nil, newm_helper.Trace(err)
	}

	for userId, courses := range coursesMap {
		var c CourseForReport

		userInfo, err := s.rUser.GetUserInfo(userId)
		if err != nil {
			return nil, newm_helper.Trace(err)
		}

		for _, v := range courses {
			c.IconLink = v.IconLink
			c.Link = v.Link
			c.Name = v.Name
			c.Platform = v.Platform
			c.Author = v.Author
			c.UserId = userId
			c.StudentInfo = fmt.Sprintf("%s %s %s", userInfo.MiddleName, userInfo.Name, userInfo.LastName)

			files, err := s.rFile.GetEducationFilesByCourseId(v.Id, userId)
			if err != nil {
				return nil, newm_helper.Trace(err)
			}

			for _, v := range files {
				c.FileLinks = append(c.FileLinks, v.Directory+"/"+v.Name)
			}
		}

		resp = append(resp, c)
	}

	return resp, nil
}


/*
группа +
фио студента +
кредиты -
баллы -
дисциплина -
имя ресурса +
дата подачи +
подпись студента -
специальность и номер +
Контактні особи, які можуть підтвердити факт навчання -
дата прохождения +
члены коммисии и их подписи -
*/

func (s *fileService) TestPdf() error {

	file, err := os.Open("test.html")
	if err != nil {
		return err
	}
	defer file.Close()

	// Создаем экземпляр PDF конвертера
	converter, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return err
	}

	// Создаем новую страницу для конвертации
	page := wkhtmltopdf.NewPageReader(file)
	converter.AddPage(page)

	// Генерируем PDF
	err = converter.Create()
	if err != nil {
		return err
	}

	// Сохраняем PDF файл
	err = converter.WriteFile("output.pdf")
	if err != nil {
		return err
	}

	return nil
}

func (s *fileService) DeleteEducationFile(fileId int) error {
	return s.rFile.DeleteEducationFile(fileId)
}

func (s *fileService) DeleteReportFile(fileId int) error {
	return s.rFile.DeleteReportFile(fileId)
}

func (s *fileService) GetEducationFileById(fileId int) ([]byte, error) {
	return s.rFile.GetEducationFileById(fileId)
}

func (s *fileService) GetReportFileById(fileId int) ([]byte, error) {
	return s.rFile.GetReportFileById(fileId)
}

func (s *fileService) GetEducationFilesInfoByUserId(userId int) (dto.GetFileResponse, error) {
	files, err := s.rFile.GetEducationFilesInfoByUserId(userId)
	if err != nil {
		return dto.GetFileResponse{}, newm_helper.Trace(err)
	}

	return dto.NewGetFilesResponse(files), nil
}

func (s *fileService) GetReportFilesInfoByUserId(userId int) (dto.GetFileResponse, error) {
	files, err := s.rFile.GetReportFilesInfoByUserId(userId)
	if err != nil {
		return dto.GetFileResponse{}, newm_helper.Trace(err)
	}

	return dto.NewGetFilesResponse(files), nil
}

func (s *fileService) CreateReportFile(file dto.CreateFileRequest, link string) error {
	course, err := s.rCourse.GetCourseByLink(link)
	if err != nil {
		return newm_helper.Trace(err)
	}

	return s.rFile.CreateReportFile(file.FileBytes, file.FileType, file.UserId, course.Id)
}

func (s *fileService) CreateEducationFile(file dto.CreateFileRequest, link string) error {
	course, err := s.rCourse.GetCourseByLink(link)
	if err != nil {
		return newm_helper.Trace(err)
	}

	err = s.rFile.CreateEducationFile(file.FileBytes, file.FileType, file.UserId, course.Id)
	if err != nil {
		return newm_helper.Trace(err)
	}

	err = s.rCourse.SetCheckCourseUser(course.Id, file.UserId)
	if err != nil {
		return newm_helper.Trace(err)
	}

	return nil
}
