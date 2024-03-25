package user

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type IUserRepo interface {
	CreateUser(user User) error
	GetUser(login, password string) (User, error)
}

type userService struct {
	r IUserRepo
}

func NewUserService(r IUserRepo) IUserService {
	return &userService{r: r}
}

func (s *userService) CreateUser(user User) error {
	user.Password = generatePasswordHash(user.Password)

	return s.r.CreateUser(user)
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
