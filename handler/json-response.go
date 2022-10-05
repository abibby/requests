package handler

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	BaseResponse
	data any
}

var _ Response = &JSONResponse{}

func NewJSONResponse(data any) *JSONResponse {
	return &JSONResponse{
		data: data,
	}
}

func (r *JSONResponse) Write(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(r.data)
}
