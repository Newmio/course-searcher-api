package dto

import "fmt"

type GetCourseRequest struct {
	Page        int    `json:"page" xml:"page"`
	SearchValue string `json:"search_value" xml:"search_value"`
}

func (c GetCourseRequest) Validate() error {
	if c.SearchValue == "" {
		return fmt.Errorf("empty search value")
	}
	return nil
}
