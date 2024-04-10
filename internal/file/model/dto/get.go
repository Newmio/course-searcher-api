package dto

import "searcher/internal/file/model/entity"

type GetFileResponse struct {
	Id   int    `json:"id" xml:"id"`
	Name string `json:"name" xml:"name"`
}

func NewGetFilesResponse(file []entity.GetFile) []GetFileResponse {
	var files []GetFileResponse

	for _, v := range file {
		files = append(files, GetFileResponse{Id: v.Id, Name: v.Name})
	}

	return files
}
