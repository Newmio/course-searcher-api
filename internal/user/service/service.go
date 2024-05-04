package service

import (
	"crypto/sha1"
	"fmt"
	"searcher/internal/user/model/dto"
	"searcher/internal/user/model/entity"
	"searcher/internal/user/repository"
	"time"

	"github.com/Newmio/newm_helper"
	"github.com/golang-jwt/jwt"
)

const (
	SALT       = "fwfjsndfwdqwdqwuriiotncna23219nsncjancasncuenfen834832u0423u094239jdjsanjsiqepee33425e1rqwftdyvghsuqw78e6trgdhbsuw3e7ref"
	ACCESSKEY    = "ncaeuwbcewr43943qfb8340hdq4t93q48ugmx9bgbfbydsbufxy6g37b2qg6fxbg67b4gfbxq6xf7x349q6gf76gew7gqf67xg4qf76g437fggf6gwefg"
	REFRESHKEY = "fjdhsjdkfjcnfjsoeorowmamxnswyfjvkjxkisporognfhsuwjeoosjshgdivifiebejhwefjweooeiwoirbvcbnnsuweybfwfbdyhbwybfwueyyw" +
		"weybfiwbfhbveywbfiewbfniewhnfdiuqwenfiewbvhwtrngiqubnrifybieryqbywwqebgcmoquyxmouyrgfxoqurguyxmuy4ghxueyaugxmaeygmfxaeymuyeggyu" +
		"fuebfuyerbavuearybfuyebrfuberauygboreyabguvyabhsuhfbuvyaewrbfouyaerbgvuobzbmyuaweogmfoeruyagouygbaeuyrbguyearbgxemryaguybagbgyu" +
		"yergxureygnsxuyemrshyghieuhrgyuxmaberygmxeauyoghmxuyearhguyfrbegbveauhgfurheag7heargiufvhnraiyhgfyreghaerhgurehgiurhgurahgfurhg" +
		"zmyerwbvzeroivrteiuwhbmouvxyhxmbeuywhgxvneryagnuxearuybveiruhfuihaweyfhaeruygfbuyhbuywqbufyberuybofnuerogbyeurbgyerbaguyebruagb" +
		"nvxertwyfuviywerbmfxgerwuygbuywerbgfuyeruyfyrsuegbuvydfbaguyraeuygfuyaergfuyoaegrhufyhgaeruygfyuaerghfuyhaeruygfyuaegfyuegrayuf" +
		"fanbeuyrbfueyroabgsrbxekrighxmhuysehrmxyherysughyesurbguyershgiyhearigyuheuoryghfuyerahgfyuiahsdfiygheryaghfuioydshgfyrhyrhgyrh"
	TOKENTIME = 7200
	REFRESHTIME = 86400 * 14
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int    `json:"id_user"`
	Role   string `json:"role"`
}

type IUserService interface {
	CreateUser(user dto.RegisterUserRequest) error
	Login(userReq dto.LoginUserRequest) (dto.LoginUserResponse, error)
	UpdateUser(user dto.UpdateUserRequest) error
	UpdatePassword(user dto.UpdateUserPasswordRequest) error
	UpdateUserInfo(info dto.CreatUserInfoRequest) error
}

type userService struct {
	r repository.IUserRepo
}

func NewUserService(r repository.IUserRepo) IUserService {
	return &userService{r: r}
}

func (s *userService) GetUserInfo(userId int) (dto.GetUserInfoResponse, error){
	info, err := s.r.GetUserInfo(userId)
	if err != nil {
		return dto.GetUserInfoResponse{}, err
	}

	return entity.NewGetUserInfoResponse(info), nil
}

func (s *userService) UpdateUserInfo(info dto.CreatUserInfoRequest) error{
	return s.r.UpdateUserInfo(entity.NewCreateUserInfo(info))
}

func (s *userService) UpdatePassword(user dto.UpdateUserPasswordRequest) error {
	user.Password = generatePasswordHash(user.Password)
	return s.r.UpdatePassword(entity.NewUpdateUserPassword(user))
}

func (s *userService) UpdateUser(user dto.UpdateUserRequest) error {
	return s.r.UpdateUser(entity.NewUpdateUser(user))
}

func (s *userService) Login(userReq dto.LoginUserRequest) (dto.LoginUserResponse, error) {
	user, err := s.r.GetUser(userReq.Login, generatePasswordHash(userReq.Password))
	if err != nil {
		return dto.LoginUserResponse{}, err
	}

	if user.Id == 0 {
		return dto.LoginUserResponse{}, fmt.Errorf("user not found")
	}

	access, refresh, err := s.GenerateToken(user.Id, user.Role)
	if err != nil {
		return dto.LoginUserResponse{}, err
	}

	return dto.NewLoginUserResponse(access, refresh, TOKENTIME, REFRESHTIME), nil
}

func (s *userService) CreateUser(user dto.RegisterUserRequest) error {
	user.Password = generatePasswordHash(user.Password)

	return s.r.CreateUser(entity.NewCreateUser(user))
}

func (s *userService) GenerateToken(id int, role string) (string, string, error) {
	accessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TOKENTIME * time.Second).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		id,
		role,
	})

	refreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(REFRESHTIME * time.Second).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		id,
		role,
	})

	accessToken, err := accessClaims.SignedString([]byte(ACCESSKEY))
	if err != nil {
		return "", "", newm_helper.Trace(err)
	}

	refreshToken, err := refreshClaims.SignedString([]byte(REFRESHKEY))
	if err != nil {
		return "", "", newm_helper.Trace(err)
	}

	return accessToken, refreshToken, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(SALT)))
}
