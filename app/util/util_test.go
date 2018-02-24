package util

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestWriteJSONErrorResponse(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		WriteJSONErrorResponse(w, errors.New("Test error"))
	}).Methods("GET")

	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
