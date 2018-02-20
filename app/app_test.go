package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/galactic-filament/go-home/app/defaultmanager"
	"github.com/galactic-filament/go-home/app/postmanager"
	"github.com/galactic-filament/go-home/app/testhandler"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// global test handler
var th testhandler.TestHandler

// main
func TestMain(m *testing.M) {
	var err error
	th.Db, err = sqlx.Connect(
		"postgres",
		fmt.Sprintf("postgres://postgres@%s/postgres?sslmode=disable", os.Getenv("DATABASE_HOST")),
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	os.Exit(m.Run())
}

// convenience methods
func createPost(th testhandler.TestHandler, requestPost postmanager.PostRequest) (post postmanager.Post) {
	// generating a request payload
	payload, err := json.Marshal(requestPost)
	assert.Nil(th.T, err, "Could not marshal post")

	// requesting
	w := th.TestPostJSONRequest("/posts", bytes.NewBuffer(payload), http.StatusCreated)

	// asserting that the post id is returned
	err = json.NewDecoder(w.Body).Decode(&post)
	assert.Nil(th.T, err, "Could not decode response body")
	assert.NotNil(th.T, post.ID, "Post id is nil")

	return post
}

// actual tests
func TestHomepage(t *testing.T) {
	// update the test handler with the test runner
	th.T = t

	// attempt a request
	w := th.TestGetRequest("/", http.StatusOK)
	assert.Equal(t, "Hello, world!", w.Body.String())
}

func TestPing(t *testing.T) {
	// update the test handler with the test runner
	th.T = t

	// attempt a request
	w := th.TestGetRequest("/ping", http.StatusOK)
	assert.Equal(t, "Pong", w.Body.String())
}

func TestReflection(t *testing.T) {
	// update the test handler with the test runner
	th.T = t

	// generating a request payload
	requestGreeting := defaultmanager.GreetingRequest{Greeting: "Hello, world!"}
	payload, err := json.Marshal(requestGreeting)
	assert.Nil(t, err, "Could not marshal greeting")

	// requesting
	w := th.TestPostJSONRequest("/reflection", bytes.NewBuffer(payload), http.StatusOK)

	// asserting that the request and response match
	var responseGreeting defaultmanager.GreetingRequest
	err = json.NewDecoder(w.Body).Decode(&responseGreeting)
	assert.Nil(t, err, "Could not decode response body")
	assert.Equal(t, requestGreeting.Greeting, responseGreeting.Greeting)
}

func TestPosts(t *testing.T) {
	// update the test handler with the test runner
	th.T = t

	// creating a post
	createPost(th, postmanager.PostRequest{Body: "Hello, world!"})
}

func TestGetPost(t *testing.T) {
	// update the test handler with the test runner
	th.T = t

	// creating a post
	createPostResponse := createPost(th, postmanager.PostRequest{Body: "Hello, world!"})

	// requesting
	w := th.TestGetJSONRequest(fmt.Sprintf("/post/%d", createPostResponse.ID), http.StatusOK)

	// asserting that the bodies match
	var getPostResponse postmanager.Post
	err := json.NewDecoder(w.Body).Decode(&getPostResponse)
	assert.Nil(t, err, "Could not decode response body")
	assert.Equal(t, createPostResponse.Body, getPostResponse.Body)
}

func TestDeletePost(t *testing.T) {
	// update the test handler with the test runner
	th.T = t

	// creating a post
	createPostResponse := createPost(th, postmanager.PostRequest{Body: "Hello, world!"})

	// requesting
	w := th.TestDeleteJSONRequest(fmt.Sprintf("/post/%d", createPostResponse.ID), http.StatusOK)

	// asserting that the bodies match
	var deletePostResponse postmanager.DeleteResponse
	err := json.NewDecoder(w.Body).Decode(&deletePostResponse)
	assert.Nil(t, err, "Could not decode response body")
}

func TestPutPost(t *testing.T) {
	// update the test handler with the test runner
	th.T = t

	// creating a post
	postRequest := postmanager.PostRequest{Body: "Hello, world!"}
	createPostResponse := createPost(th, postRequest)

	// generating a request payload
	putPost := postmanager.Post{Body: "Jello, world!"}
	payload, err := json.Marshal(putPost)
	assert.Nil(th.T, err, "Could not marshal post")

	// requesting
	w := th.TestPutJSONRequest(fmt.Sprintf("/post/%d", createPostResponse.ID), bytes.NewBuffer(payload), http.StatusOK)

	// asserting that the bodies match
	var post postmanager.Post
	err = json.NewDecoder(w.Body).Decode(&post)
	assert.Nil(t, err, "Could not decode response body")
	assert.Equal(t, putPost.Body, post.Body)
}
