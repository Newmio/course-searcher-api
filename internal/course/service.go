package course

type ICourseRepo interface{
	ShortSearchCourse(valueSearch string, inDescription bool)([]Course, error)
}

type courseService struct {
	r ICourseRepo
}

func NewCourseService(r ICourseRepo,) ICourseService {
	return &courseService{r: r}
}

func (s *courseService) ShortSearchCourse(searchValue string, inDescription bool)([]Course, error){
	return s.r.ShortSearchCourse(searchValue, inDescription)
}