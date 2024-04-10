package service

import (
	"fmt"
	"os"
	"searcher/internal/file/model/dto"
	"searcher/internal/file/model/entity"
	rFile "searcher/internal/file/repository"
	rUser "searcher/internal/user/repository"
	"time"

	"github.com/Newmio/newm_helper"
)

type IFileService interface {
	CreateReportFile(file dto.CreateFileRequest) error
	CreateEducationFile(file dto.CreateFileRequest) error
}

type fileService struct {
	rFile rFile.IFileRepo
	rUser rUser.IUserRepo
}

func NewFileService(rFile rFile.IFileRepo, rUser rUser.IUserRepo) IFileService {
	return &fileService{rFile: rFile, rUser: rUser}
}

func (s *fileService) CreateReportFile(file dto.CreateFileRequest) error {

	dir, err := createFile(file.FileBytes, file.FileType)
	if err != nil {
		return newm_helper.Trace(err)
	}

	return s.rFile.CreateReportFile(entity.NewCreateFile(dir))
}

func (s *fileService) CreateEducationFile(file dto.CreateFileRequest) error {

	dir, err := createFile(file.FileBytes, file.FileType)
	if err != nil {
		return newm_helper.Trace(err)
	}

	return s.rFile.CreateEducationFile(entity.NewCreateFile(dir))
}

func createFile(bodyBytes []byte, fileType string) (string, error) {
	for {
		name := fmt.Sprint(time.Now().UnixNano())

		if checkExistsFile(name) {
			continue
		}

		dir := fmt.Sprintf("media/%s.%s", name, fileType)

		file, err := os.Create(dir)
		if err != nil {
			return "", newm_helper.Trace(err)
		}

		_, err = file.Write(bodyBytes)
		if err != nil {
			return "", newm_helper.Trace(err)
		}

		return dir, nil
	}
}

func checkExistsFile(directory string) bool {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return false
	}
	return true
}
