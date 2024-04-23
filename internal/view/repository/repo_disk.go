package repository

import (
	"os"
	"strings"

	"github.com/Newmio/newm_helper"
)

type IDiskViewRepo interface {
	GetAllDefaultAvatarNames() ([]string, error)
}

type diskViewRepo struct {
}

func NewDiskViewRepo() IDiskViewRepo {
	return &diskViewRepo{}
}

func (r *diskViewRepo) GetAllDefaultAvatarNames() ([]string, error) {
	var names []string

	dir, err := os.Open("template/user/profile/avatars")
	if err != nil {
		return nil, newm_helper.Trace(err)
	}
	defer dir.Close()

	fileNames, err := dir.Readdirnames(-1)
	if err != nil {
		return nil, newm_helper.Trace(err)
	}

	for _, value := range fileNames {
		if strings.HasPrefix(value, "default_") {
			names = append(names, value)
		}
	}

	return names, nil
}
