package service

import (
	repoUser "searcher/internal/user/repository"

	"github.com/Newmio/newm_helper"
)

type IViewService interface {
	GetUserProfile(id int)(string, error)
}

type viewService struct {
	rUser repoUser.IUserRepo
}

func NewViewService(rUser repoUser.IUserRepo) IViewService {
	return &viewService{rUser: rUser}
}

func (s *viewService) GetUserProfile(id int)(string, error){
	user, err := s.rUser.GetUserById(id)
	if err != nil {
		return "", err
	}

	if user.Phone == "" {
		user.Phone = "Номер не указан"
	}
	
	return newm_helper.RenderHtml("template/user/profile/profile.html", user)
}