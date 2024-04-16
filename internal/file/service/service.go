package service

import (
	"os"
	"searcher/internal/file/model/dto"
	rFile "searcher/internal/file/repository"
	rUser "searcher/internal/user/repository"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type IFileService interface {
	CreateReportFile(file dto.CreateFileRequest) error
	CreateEducationFile(file dto.CreateFileRequest) error
	GetReportFilesInfoByUserId(userId int) (dto.GetFileResponse, error)
	GetEducationFilesInfoByUserId(userId int) (dto.GetFileResponse, error)
	GetReportFileById(fileId int) ([]byte, error)
	GetEducationFileById(fileId int) ([]byte, error)
	DeleteReportFile(fileId int) error
	DeleteEducationFile(fileId int) error

	TestPdf() error
}

type fileService struct {
	rFile rFile.IFileRepo
	rUser rUser.IUserRepo
}

func NewFileService(rFile rFile.IFileRepo, rUser rUser.IUserRepo) IFileService {
	return &fileService{rFile: rFile, rUser: rUser}
}

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
		return dto.GetFileResponse{}, err
	}

	return dto.NewGetFilesResponse(files), nil
}

func (s *fileService) GetReportFilesInfoByUserId(userId int) (dto.GetFileResponse, error) {
	files, err := s.rFile.GetReportFilesInfoByUserId(userId)
	if err != nil {
		return dto.GetFileResponse{}, err
	}

	return dto.NewGetFilesResponse(files), nil
}

func (s *fileService) CreateReportFile(file dto.CreateFileRequest) error {
	return s.rFile.CreateReportFile(file.FileBytes, file.FileType, file.UserId)
}

func (s *fileService) CreateEducationFile(file dto.CreateFileRequest) error {
	return s.rFile.CreateEducationFile(file.FileBytes, file.FileType, file.UserId)
}
