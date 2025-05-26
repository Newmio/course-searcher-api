package entity

type GetFile struct {
	Id         int
	Name       string
	Type       string
	Directory  string
	DateCreate string `db:"date_create"`
}

type GetReportFileUser struct {
	Id           int
	IdReportFile int `db:"id_report_file"`
	IdUser       int `db:"id_user"`
	IdCourse     int `db:"id_course"`
}

type GetEducationFileUser struct {
	Id              int
	IdEducationFile int `db:"id_education_file"`
	IdUser          int `db:"id_user"`
	IdCourse        int `db:"id_course"`
}
