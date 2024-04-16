package dto

import "searcher/internal/file/model/entity"

type FileList struct {
	Id   int    `json:"id" xml:"id"`
	Name string `json:"name" xml:"name"`
}

type GetFileResponse struct {
	Files []FileList `json:"files" xml:"files"`
	Count int        `json:"count" xml:"count"`
}

func NewGetFilesResponse(file []entity.GetFile) GetFileResponse {
	var files []FileList

	for _, v := range file {
		files = append(files, FileList{Id: v.Id, Name: v.Name})
	}

	return GetFileResponse{
		Files: files,
		Count: len(files),
	}
}
