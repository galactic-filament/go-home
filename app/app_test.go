package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var db *sqlx.DB

func testRequest(t *testing.T, method string, dest string, body io.Reader) *httptest.ResponseRecorder {
	// fetching the request router
	r := getHandler(db)

	// generating a request to test it
	req, err := http.NewRequest(method, dest, body)
	assert.Nil(t, err,
		fmt.Sprintf("Could not create new %s %s request", method, dest),
	)

	// serving up a single request and recording the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// asserting that it worked properly
	assert.Equal(t, http.StatusOK, w.Code, "Response code was not 200")
	return w
}

func testGetRequest(t *testing.T, dest string) *httptest.ResponseRecorder {
	return testRequest(t, "GET", dest, nil)
}

func testPostRequest(t *testing.T, dest string, payload io.Reader) *httptest.ResponseRecorder {
	w := testRequest(t, "POST", dest, payload)
	assert.Equal(t, "application/json", w.Header().Get("Content-type"), "Response content-type was not application/json")
	return w
}

func TestMain(m *testing.M) {
	var err error
	db, err = sqlx.Connect("postgres", "postgres://postgres@db/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err.Error())
	}
	os.Exit(m.Run())
}

func TestHomepage(t *testing.T) {
	w := testGetRequest(t, "/")
	assert.Equal(t, "Hello, world!", w.Body.String())
}

func TestPing(t *testing.T) {
	w := testGetRequest(t, "/ping")
	assert.Equal(t, "Pong", w.Body.String())
}

func TestReflection(t *testing.T) {
	// generating a request payload
	requestGreeting := greeting{Greeting: "Hello, world!"}
	payload, err := json.Marshal(requestGreeting)
	assert.Nil(t, err, "Could not marshal greeting")

	// requesting
	w := testPostRequest(t, "/reflection", bytes.NewBuffer(payload))

	// asserting that the request and response match
	var responseGreeting greeting
	err = json.NewDecoder(w.Body).Decode(&responseGreeting)
	assert.Nil(t, err, "Could not decode response body")
	assert.Equal(t, requestGreeting.Greeting, responseGreeting.Greeting)
}
