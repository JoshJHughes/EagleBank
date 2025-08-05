package web

import (
	"encoding/json"
	"net/http"
)

func writeErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	resp := newErrorResponse(err)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

func writeBadRequestErrorResponse(w http.ResponseWriter, err error) {
	resp := newBadRequestErrorResponse(err)
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(resp)
}
