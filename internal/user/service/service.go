package service

import (
	"crypto/sha1"
	"fmt"
	"searcher/internal/user/model/dto"
	"searcher/internal/user/model/entity"
	"searcher/internal/user/repository"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	SALT    = "fwfjsndfwdqwdqwuriiotncna23219nsncjancasncuenfen834832u0423u094239jdjsanjsiqepee33425e1rqwftdyvghsuqw78e6trgdhbsuw3e7ref"
	SIGNKEY = "ncaeuwbcewr43943qfb8340hdq4t93q48ugmx9bgbfbydsbufxy6g37b2qg6fxbg67b4gfbxq6xf7x349q6gf76gew7gqf67xg4qf76g437fggf6gwefg"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int    `json:"id_user"`
	Role   string `json:"role"`
}

type IUserService interface {
	CreateUser(user dto.RegisterUserRequest) error
}

type userService struct {
	r repository.IUserRepo
}

func NewUserService(r repository.IUserRepo) IUserService {
	return &userService{r: r}
}

func (s *userService) CreateUser(user dto.RegisterUserRequest) error {
	user.Password = generatePasswordHash(user.Password)

	return s.r.CreateUser(entity.NewCreateUser(user))
}

func (s *userService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(SIGNKEY), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, fmt.Errorf("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func (s *userService) GenerateToken(username, password string) (string, error) {
	user, err := s.r.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	if user.Id == 0 {
		return "", fmt.Errorf("user not found")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
		user.Role,
	})

	return token.SignedString([]byte(SIGNKEY))
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(SALT)))
}
