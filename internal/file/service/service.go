package service

import (
	"searcher/internal/file/model/dto"
	rFile "searcher/internal/file/repository"
	rUser "searcher/internal/user/repository"
)

type IFileService interface {
	CreateReportFile(file dto.CreateFileRequest) error
	CreateEducationFile(file dto.CreateFileRequest) error
	GetReportFilesInfoByUserId(userId int) ([]dto.GetFileResponse, error)
	GetEducationFilesInfoByUserId(userId int) ([]dto.GetFileResponse, error)
	GetReportFileById(fileId int) ([]byte, error)
	GetEducationFileById(fileId int) ([]byte, error)
}

type fileService struct {
	rFile rFile.IFileRepo
	rUser rUser.IUserRepo
}

func NewFileService(rFile rFile.IFileRepo, rUser rUser.IUserRepo) IFileService {
	return &fileService{rFile: rFile, rUser: rUser}
}

func (s *fileService) GetEducationFileById(fileId int) ([]byte, error) {
	return s.rFile.GetEducationFileById(fileId)
}

func (s *fileService) GetReportFileById(fileId int) ([]byte, error) {
	return s.rFile.GetReportFileById(fileId)
}

func (s *fileService) GetEducationFilesInfoByUserId(userId int) ([]dto.GetFileResponse, error) {
	files, err := s.rFile.GetEducationFilesInfoByUserId(userId)
	if err != nil {
		return nil, err
	}

	return dto.NewGetFilesResponse(files), nil
}

func (s *fileService) GetReportFilesInfoByUserId(userId int) ([]dto.GetFileResponse, error) {
	files, err := s.rFile.GetReportFilesInfoByUserId(userId)
	if err != nil {
		return nil, err
	}

	return dto.NewGetFilesResponse(files), nil
}

func (s *fileService) CreateReportFile(file dto.CreateFileRequest) error {
	return s.rFile.CreateReportFile(file.FileBytes, file.FileType)
}

func (s *fileService) CreateEducationFile(file dto.CreateFileRequest) error {
	return s.rFile.CreateEducationFile(file.FileBytes, file.FileType)
}
