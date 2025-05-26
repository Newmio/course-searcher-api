package repository

import (
	"fmt"
	"searcher/internal/file/model/entity"

	"github.com/Newmio/newm_helper"
	"github.com/jmoiron/sqlx"
)

type IFileRepo interface {
	CreateReportFile(fileBytes []byte, fileType string, userId, courseId int) error
	CreateEducationFile(fileBytes []byte, fileType string, userId, courseId int) error
	GetReportFilesInfoByUserId(userId int) ([]entity.GetFile, error)
	GetEducationFilesInfoByUserId(userId int) ([]entity.GetFile, error)
	GetReportFileById(fileId int) ([]byte, error)
	GetEducationFileById(fileId int) ([]byte, error)
	DeleteReportFile(fileId int) error
	DeleteEducationFile(fileId int) error
	GetEducationFilesByCourseId(courseId, userId int) ([]entity.GetFile, error)
}

type IDiskFileRepo interface {
	CreateFile(bodyBytes []byte, fileType string) (string, error)
	GetFile(directory string) ([]byte, error)
	DeleteFile(directory string) error
}

type IPsqlFileRepo interface {
	CreateReportFile(file entity.CreateFile) error
	CreateEducationFile(file entity.CreateFile) error
	GetReportFilesByUserId(userId int) ([]entity.GetFile, error)
	GetEducationFilesByUserId(userId int) ([]entity.GetFile, error)
	GetReportFileById(id int) (entity.GetFile, error)
	GetEducationFileById(id int) (entity.GetFile, error)
	DeleteReportFile(id int) (string, error)
	DeleteEducationFile(id int) (string, error)
	GetEducationFilesByCourseId(courseId, userId int) ([]entity.GetFile, error)
}

type managerFileRepo struct {
	disk IDiskFileRepo
	psql IPsqlFileRepo
}

func NewManagerFileRepo(psql *sqlx.DB) IFileRepo {
	psqlRepo := NewPsqlFileRepo(psql)
	diskFileRepo := NewDiskFileRepo()
	return &managerFileRepo{psql: psqlRepo, disk: diskFileRepo}
}

func (r *managerFileRepo) GetEducationFilesByCourseId(courseId, userId int) ([]entity.GetFile, error) {
	return r.psql.GetEducationFilesByCourseId(courseId, userId)
}

func (r *managerFileRepo) DeleteEducationFile(fileId int) error {
	dir, err := r.psql.DeleteEducationFile(fileId)
	if err != nil {
		return newm_helper.Trace(err)
	}

	return r.disk.DeleteFile(dir)
}

func (r *managerFileRepo) DeleteReportFile(fileId int) error {
	dir, err := r.psql.DeleteReportFile(fileId)
	if err != nil {
		return newm_helper.Trace(err)
	}

	return r.disk.DeleteFile(dir)
}

func (r *managerFileRepo) GetEducationFileById(fileId int) ([]byte, error) {
	file, err := r.psql.GetEducationFileById(fileId)
	if err != nil {
		return nil, newm_helper.Trace(err)
	}

	if file.Name == "" {
		return nil, nil
	}

	return r.disk.GetFile(fmt.Sprintf("%s/%s", file.Directory, file.Name))
}

func (r *managerFileRepo) GetReportFileById(fileId int) ([]byte, error) {
	file, err := r.psql.GetReportFileById(fileId)
	if err != nil {
		return nil, newm_helper.Trace(err)
	}

	if file.Name == "" {
		return nil, nil
	}

	return r.disk.GetFile(fmt.Sprintf("%s/%s", file.Directory, file.Name))
}

func (r *managerFileRepo) GetReportFilesInfoByUserId(userId int) ([]entity.GetFile, error) {
	return r.psql.GetReportFilesByUserId(userId)
}

func (r *managerFileRepo) GetEducationFilesInfoByUserId(userId int) ([]entity.GetFile, error) {
	return r.psql.GetEducationFilesByUserId(userId)
}

func (r *managerFileRepo) CreateReportFile(fileBytes []byte, fileType string, userId, courseId int) error {
	dir, err := r.disk.CreateFile(fileBytes, fileType)
	if err != nil {
		return newm_helper.Trace(err)
	}

	return r.psql.CreateReportFile(entity.NewCreateFile(dir, userId, courseId))
}

func (r *managerFileRepo) CreateEducationFile(fileBytes []byte, fileType string, userId, courseId int) error {
	dir, err := r.disk.CreateFile(fileBytes, fileType)
	if err != nil {
		return newm_helper.Trace(err)
	}

	return r.psql.CreateEducationFile(entity.NewCreateFile(dir, userId, courseId))
}
