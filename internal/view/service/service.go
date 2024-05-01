package service

import (
	repoUser "searcher/internal/user/repository"
	repoCourse "searcher/internal/course/repository"
	"searcher/internal/view/repository"

	"github.com/Newmio/newm_helper"
)

type IViewService interface {
	GetUserProfile(id int, directory string) (string, error)
	GetAllDefaultAvatarNames() (string, error)
}

type viewService struct {
	r       repository.IDiskViewRepo
	rUser   repoUser.IUserRepo
	rCourse repoCourse.ICourseRepo
}

func NewViewService(r repository.IDiskViewRepo, rUser repoUser.IUserRepo, rCourse repoCourse.ICourseRepo) IViewService {
	return &viewService{r: r, rUser: rUser, rCourse: rCourse}
}

func (s *viewService) GetUserProfile(id int, directory string) (string, error) {
	user, err := s.rUser.GetUserById(id)
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	return newm_helper.RenderHtml(directory, user)
}

func (s *viewService) GetAllDefaultAvatarNames() (string, error) {
	var data struct {
		Names []string
	}

	names, err := s.r.GetAllDefaultAvatarNames()
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	data.Names = names

	return newm_helper.RenderHtml("template/user/profile/update/update.html", data)
}
