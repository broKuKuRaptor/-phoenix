package common

import (
	"encoding/json"
	"net/http"
)

// RespondJSON serializes v as JSON and writes it to the response
// with the given HTTP status code and Content-Type: application/json.
func RespondJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// RespondError writes a JSON error response of the form {"error": "..."}.
//
// Typical usage:
//
//	common.RespondError(w, http.StatusNotFound, err)
//	common.RespondError(w, http.StatusInternalServerError, err)
func RespondError(w http.ResponseWriter, status int, err error) error {
	return RespondJSON(w, status, map[string]string{"error": err.Error()})
}
