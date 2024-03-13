package user

type IUserService interface {
	
}

type userService struct {
	r IPsqlUserRepo
}

func NewUserService(r IPsqlUserRepo) IUserService {
	return &userService{r: r}
}