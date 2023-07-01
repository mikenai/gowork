package handler

import "net/http"

// Healthz is used to check if everything is working properly.
//
//	GET /_healthz
//	Responds: 200
func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
