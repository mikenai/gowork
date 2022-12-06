package response

import (
	"encoding/json"
	"net/http"
)

func NotFound(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func BadRequest(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func InternalError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func JSON(w http.ResponseWriter, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}

func JSONWithStatus(w http.ResponseWriter, status int, data interface{}) error {
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
