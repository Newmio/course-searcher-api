package repository

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Newmio/newm_helper"
)

type diskFileRepo struct {
}

func NewDiskFileRepo() IDiskFileRepo {
	return &diskFileRepo{}
}

func (r *diskFileRepo) DeleteFile(directory string) error {
	return os.Remove(directory)
}

func (r *diskFileRepo) GetFile(directory string) ([]byte, error) {
	file, err := os.Open(directory)
	if err != nil {
		return nil, newm_helper.Trace(err)
	}

	return io.ReadAll(file)
}

func (r *diskFileRepo) CreateFile(bodyBytes []byte, fileType string) (string, error) {
	for {
		name := fmt.Sprint(time.Now().UnixNano())

		if checkExistsFile(name) {
			continue
		}

		t := time.Now()

		dir := fmt.Sprintf("media/%d-%d-%d/", t.Day(), t.Month(), t.Year())

		if err := os.MkdirAll(dir, 0777); err != nil {
			return "", newm_helper.Trace(err)
		}

		dir = dir + fmt.Sprintf("%s.%s", name, fileType)

		file, err := os.Create(dir)
		if err != nil {
			return "", newm_helper.Trace(err)
		}

		_, err = file.Write(bodyBytes)
		if err != nil {
			return "", newm_helper.Trace(err)
		}

		return dir, nil
	}
}
func checkExistsFile(directory string) bool {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return false
	}
	return true
}
