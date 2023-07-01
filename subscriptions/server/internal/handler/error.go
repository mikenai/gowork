package handler

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func handleError(w http.ResponseWriter, err error, statusCode int, shouldLog bool) {
	if shouldLog {
		log.Error(err.Error())
	}
	w.WriteHeader(statusCode)
	errorJSON, _ := json.Marshal(struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	})

	w.Write(errorJSON)
}
