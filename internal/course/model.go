package course

type Course struct {
	Id          int    `db:"id" json:"id" xml:"id"`
	Name        string `db:"name" json:"name" xml:"name"`
	Description string `db:"description" json:"description" xml:"description"`
	Language    string `db:"language" json:"language" xml:"language"`
	Author      string `db:"author" json:"author" xml:"author"`
	Duration    string `db:"duration" json:"duration" xml:"duration"`
	Rating      string `db:"rating" json:"rating" xml:"rating"`
	Platform    string `db:"platform" json:"platform" xml:"platform"`
	Money       string    `db:"money" json:"money" xml:"money"`
	Link        string `db:"link" json:"link" xml:"link"`
}

type WebCourseParam struct {
	Url        string
	MainField  string
	Pagination bool
	Page       int
	Fields     map[string]string
}

var WebCourseParams = map[string]WebCourseParam{
	"CourseHunter": {
		Url:       "https://coursehunter.net/search?q={{searchValue}}&order_by=votes_pos&order=desc&searching=true&page={{page}}",
		MainField: "article.course",
		Fields: map[string]string{
			"name":        "h3.course-primary-name",
			"description": "div.course-description",
			"language":    "div.course-lang",
			"author":      "div.course-lessons a",
			"duration":    "div.course-duration",
			"rating":      "div.course-rating-on<>data-text",
			"money":       "div.course-status",
			"link":        "div.course-details-bottom a<>href",
		},
	},
}
