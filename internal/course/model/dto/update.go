package dto

type PutUpdateCourseRequest struct {
	Id          int    `json:"id" xml:"id"`
	Name        string `json:"name" xml:"name"`
	Description string `json:"description" xml:"description"`
	Language    string `json:"language" xml:"language"`
	Author      string `json:"author" xml:"author"`
	Duration    string `json:"duration" xml:"duration"`
	Rating      string `json:"rating" xml:"rating"`
	Platform    string `json:"platform" xml:"platform"`
	Money       string `json:"money" xml:"money"`
	Link        string `json:"link" xml:"link"`
	Active      bool   `json:"active" xml:"active"`
}

type PatchUpdateCourseRequest struct {
	Id     int               `json:"id" xml:"id"`
	Params map[string]string `json:"params" xml:"params"`
}
