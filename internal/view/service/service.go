package service

import (
	repoCourse "searcher/internal/course/repository"
	repoUser "searcher/internal/user/repository"
	"searcher/internal/view/repository"

	"github.com/Newmio/newm_helper"
)

type IViewService interface {
	GetUserProfile(id int, directory string) (string, error)
	GetAllDefaultAvatarNames(userId int) (string, error)
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

	info, err := s.rUser.GetUserInfo(user.Id)
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	data := struct {
		Id                int
		Role              string
		Email             string
		Phone             string
		Avatar            string
		Name              string
		MiddleName        string
		LastName          string
		CourseNumber      int
		GroupName         string
		Proffession       string
		ProffessionNumber string
	}{
		Id:                user.Id,
		Role:              user.Role,
		Email:             user.Email,
		Phone:             user.Phone,
		Avatar:            user.Avatar,
		Name:              info.Name,
		MiddleName:        info.MiddleName,
		LastName:          info.LastName,
		CourseNumber:      info.CourseNumber,
		GroupName:         info.GroupName,
		Proffession:       info.Proffession,
		ProffessionNumber: info.ProffessionNumber,
	}

	return newm_helper.RenderHtml(directory, data)
}

func (s *viewService) GetAllDefaultAvatarNames(userId int) (string, error) {
	user, err := s.rUser.GetUserById(userId)
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	names, err := s.r.GetAllDefaultAvatarNames()
	if err != nil {
		return "", newm_helper.Trace(err)
	}

	data := struct {
		Id    int
		Names []string
	}{
		Id:    user.Id,
		Names: names,
	}

	return newm_helper.RenderHtml("template/user/profile/update/avatar.html", data)
}
