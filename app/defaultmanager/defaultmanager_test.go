package defaultmanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHello(t *testing.T) {
	router := Init(mux.NewRouter())

	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err, fmt.Sprintf("Could not create new POST / request"))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Response was not 200 OK")
}

func TestPing(t *testing.T) {
	router := Init(mux.NewRouter())

	req, err := http.NewRequest("POST", "/ping", nil)
	assert.Nil(t, err, fmt.Sprintf("Could not create new POST / request"))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Response was not 200 OK")
}

func TestReflection(t *testing.T) {
	router := Init(mux.NewRouter())

	requestGreeting := GreetingRequest{Greeting: "Hello, world!"}
	payload, err := json.Marshal(requestGreeting)
	assert.Nil(t, err, "Could not marshal greeting")

	req, err := http.NewRequest("POST", "/reflection", bytes.NewBuffer(payload))
	assert.Nil(t, err, fmt.Sprintf("Could not create new POST / request"))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Response was not 200 OK")

	var responseGreeting GreetingRequest
	err = json.NewDecoder(w.Body).Decode(&responseGreeting)
	assert.Nil(t, err, "Could not decode response body")

	assert.Equal(t, requestGreeting, responseGreeting)
}
