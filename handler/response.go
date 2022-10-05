package handler

import "net/http"

type Response interface {
	Respond(w http.ResponseWriter) error
}
