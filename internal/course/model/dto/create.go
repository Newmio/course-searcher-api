package dto

import "fmt"

type CreateCourseRequest struct {
	Name        string `json:"name" xml:"name"`
	Description string `json:"description" xml:"description"`
	Language    string `json:"language" xml:"language"`
	Author      string `json:"author" xml:"author"`
	Duration    string `json:"duration" xml:"duration"`
	Rating      string `json:"rating" xml:"rating"`
	Platform    string `json:"platform" xml:"platform"`
	Money       string `json:"money" xml:"money"`
	Link        string `json:"link" xml:"link"`
	IconLink    string `json:"icon_link" xml:"icon_link"`
}

func (c CreateCourseRequest) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("empty name")
	}

	if c.Link == "" {
		return fmt.Errorf("empty link")
	}

	return nil
}
