package service

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	rCourse "searcher/internal/course/repository"
	"searcher/internal/file/model/dto"
	rFile "searcher/internal/file/repository"
	rUser "searcher/internal/user/repository"
	"strings"
	"time"

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
	GetReport(studentId, courseId int) (string, error)

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
	CourseId    int
	Link        string
	Name        string
	IconLink    string
	Platform    string
	Author      string
	FileLinks   []string
	StudentInfo string
	UserId      int
	EducName    bool
	Credits     bool
}

type ReportInfo struct {
	GroupName         string //
	FIO               string //
	Coins             int    //
	Credits           int
	CourseName        string //
	Platform          string //
	Proffession       string //
	EducationName     string
	CourseDescription string //
	CourseProofs      string //
	DateStart         string //
	DateStop          string //
	AllCmsFIO         string //
	DateGeneration    string //
	Name              string //
	MiddleName        string //
	ProffessionNumber string //
	CourseThemes      string
	FileLinks         []string
}

func (s *fileService) GetReport(studentId, courseId int) (string, error) {
	var repInfo ReportInfo
	var link string

	reports, err := s.rFile.GetReportFilesInfoByUserId(studentId)
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	userReport, err := s.rFile.GetReportUser(studentId, courseId)
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	for _, v := range reports {
		if v.Id == userReport {
			link = v.Directory + "/" + v.Name
		}
	}

	if link != "" {
		return link, nil
	}

	studentInfo, err := s.rUser.GetUserInfo(studentId)
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	course, err := s.rCourse.GetCourseById(courseId)
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	courseUser, err := s.rCourse.GetCourseUser(courseId, studentId)
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	userCoins := make(map[string]interface{})

	err = json.Unmarshal([]byte(courseUser["name"].(string)), &userCoins)
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	// dateStart, err := time.Parse("2006-01-02 15:04:05.000", courseUser["date_start"].(string))
	// if err != nil {
	// 	return "", newm_helper.Trace(err)
	// }

	// dateEnd, err := time.Parse("2006-01-02 15:04:05.000", courseUser["date_end"].(string))
	// if err != nil {
	// 	return "", newm_helper.Trace(err)
	// }

	cmsUsers, err := s.rUser.GetCMSUsers()
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	var cmsFios []string

	for _, v := range cmsUsers {
		cmsUserInfo, err := s.rUser.GetUserInfo(v.Id)
		if err != nil {
			return "", newm_helper.Trace(err)
		}

		cmsFios = append(cmsFios, fmt.Sprintf("%s %s %s", cmsUserInfo.MiddleName, cmsUserInfo.Name, cmsUserInfo.LastName))
	}

	files, err := s.rFile.GetEducationFilesByCourseId(courseId, studentId)
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	for _, v := range files {
		repInfo.FileLinks = append(repInfo.FileLinks, "http://localhost:8088/" + v.Directory+"/"+v.Name)
	}

	repInfo.Coins = int(userCoins["res"].(float64))
	repInfo.Credits = int(courseUser["credits"].(int64))
	repInfo.EducationName = courseUser["educ_name"].(string)
	repInfo.GroupName = studentInfo.GroupName
	repInfo.FIO = fmt.Sprintf("%s %s %s", studentInfo.MiddleName, studentInfo.Name, studentInfo.LastName)
	repInfo.Name = studentInfo.Name
	repInfo.MiddleName = studentInfo.MiddleName
	repInfo.CourseName = course.Name
	repInfo.Platform = course.Platform
	repInfo.Proffession = studentInfo.Proffession
	repInfo.ProffessionNumber = studentInfo.ProffessionNumber
	// repInfo.DateStart = dateStart.Format("02.01.2006")
	// repInfo.DateStop = dateEnd.Format("02.01.2006")
	repInfo.DateStart = courseUser["date_start"].(time.Time).Format("02.01.2006")
	repInfo.DateStop = courseUser["date_end"].(time.Time).Format("02.01.2006")
	repInfo.AllCmsFIO = strings.Join(cmsFios, ", ")
	repInfo.DateGeneration = time.Now().Format("02.01.2006")

	strHtml, err := newm_helper.RenderHtml("templatereport.html", repInfo)
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	converter, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	tempFileName := fmt.Sprintf("%d.html", time.Now().UnixNano())

	file, err := os.Create(tempFileName)
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	_, err = file.WriteString(strHtml)
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	err = file.Close()
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	fileHtml, err := os.Open(tempFileName)
	if err != nil {
		return "", err
	}

	page := wkhtmltopdf.NewPageReader(fileHtml)
	converter.AddPage(page)

	err = converter.Create()
	if err != nil {
		return "", err
	}

	tempReportFileName := fmt.Sprintf("%d.html", time.Now().UnixNano())

	err = converter.WriteFile(tempReportFileName)
	if err != nil {
		return "", err
	}

	filePdf, err := os.Open(tempReportFileName)
	if err != nil {
		return "", err
	}

	var fileCreator dto.CreateFileRequest

	fileCreator.UserId = studentId
	fileCreator.FileType = "pdf"

	fileCreator.FileBytes, err = io.ReadAll(filePdf)
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	err = s.CreateReportFile(fileCreator, course.Link)
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	err = filePdf.Close()
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	err = fileHtml.Close()
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	err = os.Remove(tempFileName)
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	err = os.Remove(tempReportFileName)
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	reports, err = s.rFile.GetReportFilesInfoByUserId(studentId)
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	userReport, err = s.rFile.GetReportUser(studentId, courseId)
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	for _, v := range reports {
		if v.Id == userReport {
			link = v.Directory + "/" + v.Name
		}
	}

	return link, nil
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
			courseUser, err := s.rCourse.GetCourseUser(v.Id, userId)
			if err != nil {
				return nil, newm_helper.Trace(err)
			}

			if v, ok := courseUser["credits"].(int64); !ok || v == 0 {
				c.Credits = true
			}

			if v, ok := courseUser["educ_name"]; !ok || v == "" {
				c.EducName = true
			}

			c.CourseId = v.Id
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
