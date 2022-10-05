package handler

import "net/http"

type Response interface {
	Write(w http.ResponseWriter) error
	Status() int
	Headers() map[string]string
}

type BaseResponse struct {
	status  int
	headers map[string]string
}

func (r *BaseResponse) SetStatus(status int) *BaseResponse {
	r.status = status
	return r
}

func (r *BaseResponse) AddHeader(key, value string) *BaseResponse {
	if r.headers == nil {
		r.headers = map[string]string{}
	}
	r.headers[key] = value
	return r
}

func (r *BaseResponse) Status() int {
	return r.status
}

func (r *BaseResponse) Headers() map[string]string {
	return r.headers
}
