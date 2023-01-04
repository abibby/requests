package handler

import (
	"log"
	"net/http"
	"reflect"

	"github.com/abibby/validate"
)

func Handler[TRequest any](callback func(r *TRequest) (any, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req TRequest
		err := validate.Run(r, &req)
		if err, ok := err.(*validate.ValidationError); ok {
			respond(w, ErrorResponse(err, http.StatusUnprocessableEntity))
			return
		} else if err != nil {
			respond(w, ErrorResponse(err, http.StatusInternalServerError))
			return
		}

		injectRequest(&req, r)
		injectResponseWriter(&req, w)

		resp, err := callback(&req)
		if err != nil {
			respond(w, ErrorResponse(err, http.StatusInternalServerError))
			return
		}
		if resp, ok := resp.(Response); ok {
			respond(w, resp)
			return
		}
		respond(w, NewJSONResponse(resp))
	})
}

func ErrorResponse(err error, status int) *JSONResponse {
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

func injectRequest[TRequest any](req TRequest, httpR *http.Request) {
	v := reflect.ValueOf(req).Elem()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if (f.Type() != reflect.TypeOf(&http.Request{})) {
			continue
		}
		f.Set(reflect.ValueOf(httpR))
	}

}

func injectResponseWriter[TRequest any](req TRequest, httpRW http.ResponseWriter) {
	v := reflect.ValueOf(req).Elem()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		var rw http.ResponseWriter
		if !f.Type().Implements(reflect.TypeOf(&rw).Elem()) {
			continue
		}
		f.Set(reflect.ValueOf(httpRW))
	}

}
