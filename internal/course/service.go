package course

type ICourseService interface {
	
}

type courseService struct {
	r ICourseRepo
}

func NewCourseService(r ICourseRepo) ICourseService {
	return &courseService{r: r}
}