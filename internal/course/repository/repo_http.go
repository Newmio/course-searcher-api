package repository

import (
	"fmt"

	"github.com/Newmio/newm_helper"
)

type httpCourseRepo struct{}

func NewHttpCourseRepo() IHttpCourseRepo {
	return &httpCourseRepo{}
}

func (r *httpCourseRepo) GetHtmlCourseInWeb(param newm_helper.Param) ([]byte, error) {
	status, body, err := newm_helper.RequestHTTP(param)
	if err != nil {
		return nil, newm_helper.Trace(err)
	}

	if status == 404 {
		return nil, nil
	}

	if status > 299 {
		return nil, newm_helper.Trace(fmt.Errorf("status code %d\n\n%s", status, string(body)))
	}

	return body, nil
}
