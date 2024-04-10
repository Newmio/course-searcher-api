package repository

import (
	"searcher/internal/file/model/entity"

	"github.com/Newmio/newm_helper"
	"github.com/jmoiron/sqlx"
)

type psqlFileRepo struct {
	db *sqlx.DB
}

func NewPsqlFileRepo(db *sqlx.DB) IPsqlFileRepo {
	r := &psqlFileRepo{db: db}
	r.initTables()
	return r
}

func (r *psqlFileRepo) CreateReportFile(file entity.CreateFile) error {
	var reportId int

	str := `insert into report_files(name, type, directory, date_create) values($1, $2, $3, $4) returning id`

	err := r.db.QueryRow(str, file.Name, file.Type, file.Directory, file.DateCreate).Scan(&reportId)
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	str = `insert into report_file_user(id_report_file, id_user) values($1, $2)`

	_, err = r.db.Exec(str, reportId, file.UserId)
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	return nil
}

func (r *psqlFileRepo) CreateEducationFile(file entity.CreateFile) error {
	var educationId int

	str := `insert into educetion_files(name, type, directory, date_create) values($1, $2, $3, $4) returning id`

	err := r.db.QueryRow(str, file.Name, file.Type, file.Directory, file.DateCreate).Scan(&educationId)
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	str = `insert into report_file_user(id_report_file, id_user) values($1, $2)`

	_, err = r.db.Exec(str, educationId, file.UserId)
	if err != nil {
		return newm_helper.Trace(err, str)
	}

	return nil
}
