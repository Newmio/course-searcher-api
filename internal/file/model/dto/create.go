package dto

import "fmt"

type CreateFileRequest struct {
	UserId    int
	FileBytes []byte
	FileType  string
}

func (c CreateFileRequest) Validate() error {
	if c.FileType == "" {
		return fmt.Errorf("empty file type")
	}

	if len(c.FileBytes) == 0 {
		return fmt.Errorf("empty file")
	}

	return nil
}
