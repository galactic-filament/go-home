package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestCoresMiddleware(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
	handler := corsMiddleware(router)

	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, "*", w.HeaderMap["Access-Control-Allow-Origin"][0])
	assert.Equal(t, "POST, GET, OPTIONS, PUT, DELETE", w.HeaderMap["Access-Control-Allow-Methods"][0])
	assert.Equal(t, "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token", w.HeaderMap["Access-Control-Allow-Headers"][0])
	assert.Equal(t, "true", w.HeaderMap["Access-Control-Allow-Credentials"][0])
}

func TestCoresMiddlewareOptionsRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
	handler := corsMiddleware(router)

	req, err := http.NewRequest("OPTIONS", "/", nil)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, "*", w.HeaderMap["Access-Control-Allow-Origin"][0])
	assert.Equal(t, "POST, GET, OPTIONS, PUT, DELETE", w.HeaderMap["Access-Control-Allow-Methods"][0])
	assert.Equal(t, "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token", w.HeaderMap["Access-Control-Allow-Headers"][0])
	assert.Equal(t, "true", w.HeaderMap["Access-Control-Allow-Credentials"][0])
}

func TestLoggingMiddleware(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
	handler := loggingMiddleware(router)

	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoggingMiddlewareWithRequestBody(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
	handler := loggingMiddleware(router)

	req, err := http.NewRequest("GET", "/", strings.NewReader("Hello, world!"))
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
