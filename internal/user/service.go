package user

type IUserRepo interface {
	
}

type userService struct {
	r IUserRepo
}

func NewUserService(r IUserRepo) IUserService {
	return &userService{r: r}
}