package repository

import (
	"searcher/internal/file/model/entity"

	"github.com/jmoiron/sqlx"
)

type IFileRepo interface {
	CreateReportFile(file entity.CreateFile) error
	CreateEducationFile(file entity.CreateFile) error
}

type IPsqlFileRepo interface {
	CreateReportFile(file entity.CreateFile) error
	CreateEducationFile(file entity.CreateFile) error
}

type managerFileRepo struct {
	psql IPsqlFileRepo
}

func NewManagerFileRepo(psql *sqlx.DB) IFileRepo {
	psqlRepo := NewPsqlFileRepo(psql)
	return &managerFileRepo{psql: psqlRepo}
}

func (r *managerFileRepo) CreateReportFile(file entity.CreateFile) error {
	return r.psql.CreateReportFile(file)
}

func (r *managerFileRepo) CreateEducationFile(file entity.CreateFile) error {
	return r.psql.CreateEducationFile(file)
}