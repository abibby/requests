package handler

import (
	"log"
	"net/http"

	"github.com/abibby/validate"
)

func Handler[TRequest any](callback func(r *TRequest) (any, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req TRequest
		err := validate.Run(r, &req)
		if err, ok := err.(*validate.ValidationError); ok {
			respond(w, errorResponse(err, http.StatusUnprocessableEntity))
			return
		} else if err != nil {
			respond(w, errorResponse(err, http.StatusInternalServerError))
			return
		}
		resp, err := callback(&req)
		if err != nil {
			respond(w, errorResponse(err, http.StatusInternalServerError))
			return
		}
		if resp, ok := resp.(Response); ok {
			respond(w, resp)
			return
		}
		respond(w, NewJSONResponse(resp))
	})
}

func errorResponse(err error, status int) *JSONResponse {
	return NewJSONResponse(map[string]string{
		"error": err.Error(),
	}).SetStatus(status)
}

func respond(w http.ResponseWriter, r Response) {
	err := r.Respond(w)
	if err != nil {
		log.Print(err)
	}
}
