package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/abibby/validate"
)

func Handler[TRequest any, TResponse Response](callback func(r *TRequest) TResponse) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req TRequest
		err := validate.Run(r, &req)
		if err, ok := err.(*validate.ValidationError); ok {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(map[string]*validate.ValidationError{
				"error": err,
			})
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}
		resp := callback(&req)
		w.WriteHeader(resp.Status())
		for k, v := range resp.Headers() {
			w.Header().Add(k, v)
		}
		err = resp.Write(w)
		if err != nil {
			log.Print(err)
		}
	})
}
