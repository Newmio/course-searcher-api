package service

import (
	"fmt"

	"github.com/Newmio/newm_helper"
	"github.com/golang-jwt/jwt"
)


const (
	ACCESSKEY    = "ncaeuwbcewr43943qfb8340hdq4t93q48ugmx9bgbfbydsbufxy6g37b2qg6fxbg67b4gfbxq6xf7x349q6gf76gew7gqf67xg4qf76g437fggf6gwefg"
	REFRESHKEY = "fjdhsjdkfjcnfjsoeorowmamxnswyfjvkjxkisporognfhsuwjeoosjshgdivifiebejhwefjweooeiwoirbvcbnnsuweybfwfbdyhbwybfwueyyw" +
		"weybfiwbfhbveywbfiewbfniewhnfdiuqwenfiewbvhwtrngiqubnrifybieryqbywwqebgcmoquyxmouyrgfxoqurguyxmuy4ghxueyaugxmaeygmfxaeymuyeggyu" +
		"fuebfuyerbavuearybfuyebrfuberauygboreyabguvyabhsuhfbuvyaewrbfouyaerbgvuobzbmyuaweogmfoeruyagouygbaeuyrbguyearbgxemryaguybagbgyu" +
		"yergxureygnsxuyemrshyghieuhrgyuxmaberygmxeauyoghmxuyearhguyfrbegbveauhgfurheag7heargiufvhnraiyhgfyreghaerhgurehgiurhgurahgfurhg" +
		"zmyerwbvzeroivrteiuwhbmouvxyhxmbeuywhgxvneryagnuxearuybveiruhfuihaweyfhaeruygfbuyhbuywqbufyberuybofnuerogbyeurbgyerbaguyebruagb" +
		"nvxertwyfuviywerbmfxgerwuygbuywerbgfuyeruyfyrsuegbuvydfbaguyraeuygfuyaergfuyoaegrhufyhgaeruygfyuaerghfuyhaeruygfyuaegfyuegrayuf" +
		"fanbeuyrbfueyroabgsrbxekrighxmhuysehrmxyherysughyesurbguyershgiyhearigyuheuoryghfuyerahgfyuiahsdfiygheryaghfuioydshgfyrhyrhgyrh"
	TOKENTIME = 7200
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int    `json:"id_user"`
	Role   string `json:"role"`
}

type IMiddlewareService interface {
	ParseToken(accessToken string) (int, string, error)
}

type middlewareService struct{}

func NewMiddlewareService() IMiddlewareService {
	return &middlewareService{}
}

func (s *middlewareService) ParseToken(accessToken string) (int, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(ACCESSKEY), nil
	})
	if err != nil {
		return 0, "", newm_helper.Trace(err)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, "", fmt.Errorf("token claims are not of type *tokenClaims")
	}

	return claims.UserId, claims.Role, nil
}
