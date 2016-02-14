package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ihsw/go-home/app/PostManager"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// test handler
type testHandler struct {
	db *sqlx.DB
	t  *testing.T
}

func (th testHandler) testRequest(method string, dest string, body io.Reader) *httptest.ResponseRecorder {
	// fetching the request router
	r := getHandler(th.db)

	// generating a request to test it
	req, err := http.NewRequest(method, dest, body)
	assert.Nil(th.t, err, fmt.Sprintf("Could not create new %s %s request", method, dest))

	// serving up a single request and recording the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// checking for 500 errors
	if w.Code == http.StatusInternalServerError {
		var errResponse errorResponse
		err = json.NewDecoder(w.Body).Decode(&errResponse)
		assert.Nil(th.t, err, "Could not decode response body")
		assert.NotNil(th.t, nil, fmt.Sprintf("Response code was 500: %s", errResponse.Error))
		return w
	}

	// asserting that it worked properly
	assert.Equal(th.t, http.StatusOK, w.Code, "Response code was not 200")
	return w
}

func (th testHandler) testGetRequest(dest string) *httptest.ResponseRecorder {
	return th.testRequest("GET", dest, nil)
}

func (th testHandler) testGetJSONRequest(dest string) *httptest.ResponseRecorder {
	w := th.testRequest("GET", dest, nil)
	assert.Equal(th.t, "application/json", w.Header().Get("Content-type"), "Response content-type was not application/json")
	return w
}

func (th testHandler) testDeleteJSONRequest(dest string) *httptest.ResponseRecorder {
	w := th.testRequest("DELETE", dest, nil)
	assert.Equal(th.t, "application/json", w.Header().Get("Content-type"), "Response content-type was not application/json")
	return w
}

func (th testHandler) testPostJSONRequest(dest string, payload io.Reader) *httptest.ResponseRecorder {
	w := th.testRequest("POST", dest, payload)
	assert.Equal(th.t, "application/json", w.Header().Get("Content-type"), "Response content-type was not application/json")
	return w
}

// global test handler
var th testHandler

// main
func TestMain(m *testing.M) {
	hostname := "db"
	if os.Getenv("ENV") == "travis" {
		hostname = "localhost"
	}

	var err error
	th.db, err = sqlx.Connect(
		"postgres",
		fmt.Sprintf("postgres://postgres@%s/postgres?sslmode=disable", hostname),
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	os.Exit(m.Run())
}

// convenience methods
func createPost(th testHandler, requestPost post) (post post) {
	// generating a request payload
	payload, err := json.Marshal(requestPost)
	assert.Nil(th.t, err, "Could not marshal post")

	// requesting
	w := th.testPostJSONRequest("/posts", bytes.NewBuffer(payload))

	// asserting that the post id is returned
	err = json.NewDecoder(w.Body).Decode(&post)
	assert.Nil(th.t, err, "Could not decode response body")
	assert.NotNil(th.t, post.ID, "Post id is nil")

	return post
}

// actual tests
func TestHomepage(t *testing.T) {
	// update the test handler with the test runner
	th.t = t

	// attempt a request
	w := th.testGetRequest("/")
	assert.Equal(t, "Hello, world!", w.Body.String())
}

func TestPing(t *testing.T) {
	// update the test handler with the test runner
	th.t = t

	// attempt a request
	w := th.testGetRequest("/ping")
	assert.Equal(t, "Pong", w.Body.String())
}

func TestReflection(t *testing.T) {
	// update the test handler with the test runner
	th.t = t

	// generating a request payload
	requestGreeting := greeting{Greeting: "Hello, world!"}
	payload, err := json.Marshal(requestGreeting)
	assert.Nil(t, err, "Could not marshal greeting")

	// requesting
	w := th.testPostJSONRequest("/reflection", bytes.NewBuffer(payload))

	// asserting that the request and response match
	var responseGreeting greeting
	err = json.NewDecoder(w.Body).Decode(&responseGreeting)
	assert.Nil(t, err, "Could not decode response body")
	assert.Equal(t, requestGreeting.Greeting, responseGreeting.Greeting)
}

func TestPosts(t *testing.T) {
	// update the test handler with the test runner
	th.t = t

	// creating a post
	createPost(th, post{Body: "Hello, world!"})
}

func TestGetPost(t *testing.T) {
	// update the test handler with the test runner
	th.t = t

	// creating a post
	createPostResponse := createPost(th, post{Body: "Hello, world!"})

	// requesting
	w := th.testGetJSONRequest(fmt.Sprintf("/post/%d", createPostResponse.ID))

	// asserting that the bodies match
	var getPostResponse post
	err := json.NewDecoder(w.Body).Decode(&getPostResponse)
	assert.Nil(t, err, "Could not decode response body")
	assert.Equal(t, createPostResponse.Body, getPostResponse.Body)
}

func TestDeletePost(t *testing.T) {
	// update the test handler with the test runner
	th.t = t

	// creating a post
	createPostResponse := createPost(th, post{Body: "Hello, world!"})

	// requesting
	w := th.testDeleteJSONRequest(fmt.Sprintf("/post/%d", createPostResponse.ID))

	// asserting that the bodies match
	var deletePostResponse PostManager.DeleteResponse
	err := json.NewDecoder(w.Body).Decode(&deletePostResponse)
	assert.Nil(t, err, "Could not decode response body")
}
