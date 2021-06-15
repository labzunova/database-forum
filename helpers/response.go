package helpers

import (
	"encoding/json"
	"net/http"
)

func CreateResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)

	return
}