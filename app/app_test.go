package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomepage(t *testing.T) {
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

func TestPing(t *testing.T) {
	// fetching the request router
	r := getHandler()

	// generating a request to test it
	req, err := http.NewRequest("GET", "/ping", nil)
	assert.Nil(t, err, "Could not create new GET /ping request")

	// serving up a single request and recording the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// asserting that it worked properly
	assert.Equal(t, http.StatusOK, w.Code, "Response code was not 200")
	assert.Equal(t, "Pong", w.Body.String())
}

func TestReflection(t *testing.T) {
	// fetching the request router
	r := getHandler()

	// generating a request payload
	requestGreeting := greeting{Greeting: "Hello, world!"}
	payload, err := json.Marshal(requestGreeting)
	assert.Nil(t, err, "Could not marshal greeting")

	// generating a request to test it
	req, err := http.NewRequest("POST", "/reflection", bytes.NewBuffer(payload))
	assert.Nil(t, err, "Could not create new POST /reflection request")

	// serving up a single request and recording the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// asserting that it worked properly
	assert.Equal(t, http.StatusOK, w.Code, "Response code was not 200")

	// aserting that the request and response match
	var responseGreeting greeting
	err = json.NewDecoder(w.Body).Decode(&responseGreeting)
	assert.Nil(t, err, "Could not decode response body")
	assert.Equal(t, requestGreeting.Greeting, responseGreeting.Greeting)
}
