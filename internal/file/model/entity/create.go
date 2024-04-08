package entity

import (
	"time"
)

type CreateFile struct {
	Name       string
	Type       string
	Directory  string
	DateCreate string
}

func NewCreateFile(directory string) CreateFile {

	//parts := strings.Split(directory, "/")

	//d := strings.Split(parts[len(parts)-1], ".")[]

	return CreateFile{
		Name:       "",
		Type:       "",
		Directory:  "",
		DateCreate: time.Now().Format("2006-01-02 15:04:05"),
	}
}
