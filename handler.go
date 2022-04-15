package validate

import (
	"encoding/json"
	"net/http"
)

func Handler[Request any](callback func(w http.ResponseWriter, r *Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req Request
		err := Run(r, &req)
		if err, ok := err.(*ValidationError); ok {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(err)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		callback(w, &req)
	})
}
