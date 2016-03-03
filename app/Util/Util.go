package Util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorResponse - error response body in json
type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteJSONErrorResponse - writes error out in json
func WriteJSONErrorResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	errResponse := ErrorResponse{Error: err.Error()}
	if err := json.NewEncoder(w).Encode(errResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not encode error response body")
		return
	}
	return
}
