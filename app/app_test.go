package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAdd(t *testing.T) {
	assert.Equal(t, 4, add(2, 2), "Add failed!")
}

func TestGetServer(t *testing.T) {
	r := getHandler()

	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err, "Could not create new GET / request")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code, "Response code was not 200")
	assert.Equal(t, "Hello, world!", w.Body.String())
}
