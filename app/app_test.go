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
	// fetching the request router
	r := getHandler()

	// generating a request to test it
	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err, "Could not create new GET / request")

	// serving up a single request and recording the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// asserting that it worked properly
	assert.Equal(t, http.StatusOK, w.Code, "Response code was not 200")
	assert.Equal(t, "Hello, world!", w.Body.String())
}
