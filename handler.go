package validate

import "net/http"

func Handler[Request any](callback func(w http.ResponseWriter, r *Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req Request
		err := Run(r, &req)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		callback(w, &req)
	})
}
