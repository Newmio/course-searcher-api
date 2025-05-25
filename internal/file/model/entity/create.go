package entity

import (
	"path/filepath"
	"time"
)

type CreateFile struct {
	UserId     int
	CourseId   int
	Name       string
	Type       string
	Directory  string
	DateCreate string
}

func NewCreateFile(directory string, userId, courseId int) CreateFile {

	dir := filepath.Dir(directory)
	name := filepath.Base(directory)
	typee := filepath.Ext(directory)

	return CreateFile{
		UserId:     userId,
		Name:       name,
		CourseId:   courseId,
		Type:       typee[1:],
		Directory:  dir,
		DateCreate: time.Now().Format("2006-01-02 15:04:05"),
	}
}
