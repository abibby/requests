package handler

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	data    any
	status  int
	headers map[string]string
}

var _ Response = &JSONResponse{}

func NewJSONResponse(data any) *JSONResponse {
	return &JSONResponse{
		data: data,
		headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func (r *JSONResponse) SetStatus(status int) *JSONResponse {
	r.status = status
	return r
}

func (r *JSONResponse) AddHeader(key, value string) *JSONResponse {
	r.headers[key] = value
	return r
}

func (r *JSONResponse) Respond(w http.ResponseWriter) error {
	w.WriteHeader(r.status)
	for k, v := range r.headers {
		w.Header().Set(k, v)
	}

	return json.NewEncoder(w).Encode(r.data)
}
