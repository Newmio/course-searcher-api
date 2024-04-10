package entity

import (
	"path/filepath"
	"time"
)

type CreateFile struct {
	UserId     int
	Name       string
	Type       string
	Directory  string
	DateCreate string
}

func NewCreateFile(directory string) CreateFile {

	dir := filepath.Dir(directory)
	name := filepath.Base(directory)
	typee := filepath.Ext(directory)

	return CreateFile{
		Name:       name,
		Type:       typee,
		Directory:  dir,
		DateCreate: time.Now().Format("2006-01-02 15:04:05"),
	}
}
