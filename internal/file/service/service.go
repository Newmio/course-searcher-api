package service

import (
	"fmt"
	"os"
	"searcher/internal/file/model/dto"
	rFile "searcher/internal/file/repository"
	rUser "searcher/internal/user/repository"
	"time"

	"github.com/Newmio/newm_helper"
)

type IFileService interface {
}

type fileService struct {
	rFile rFile.IFileRepo
	rUser rUser.IUserRepo
}

func NewFileService(rFile rFile.IFileRepo, rUser rUser.IUserRepo) IFileService {
	return &fileService{rFile: rFile, rUser: rUser}
}

func (s *fileService) CreateEducationFile(file dto.CreateFileRequest) error {

	if err := createFile(file.FileBytes, file.FileType); err != nil{
		return newm_helper.Trace(err)
	}

	return nil
}

func createFile(bodyBytes []byte, fileType string) error {
	for {
		name := fmt.Sprint(time.Now().UnixNano())

		if checkExistsFile(name) {
			continue
		}

		file, err := os.Create(name + "." + fileType)
		if err != nil {
			return newm_helper.Trace(err)
		}

		_, err = file.Write(bodyBytes)
		if err != nil {
			return newm_helper.Trace(err)
		}

		return nil
	}
}

func checkExistsFile(directory string) bool {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return false
	}
	return true
}
