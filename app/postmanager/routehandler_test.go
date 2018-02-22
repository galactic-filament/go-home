package postmanager

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var db *sqlx.DB

func createPost(router *mux.Router, requestPost PostRequest) (Post, error) {
	payload, err := json.Marshal(requestPost)
	if err != nil {
		return Post{}, err
	}

	req, err := http.NewRequest("POST", "/posts", bytes.NewBuffer(payload))
	if err != nil {
		return Post{}, err
	}

	w := httptest.NewRecorder()

	if router.ServeHTTP(w, req); w.Code != http.StatusCreated {
		return Post{}, errors.New("Could not create post")
	}

	var responsePost Post
	if err = json.NewDecoder(w.Body).Decode(&responsePost); err != nil {
		return Post{}, err
	}

	return responsePost, nil
}

func TestMain(m *testing.M) {
	var err error
	db, err = sqlx.Connect(
		"postgres",
		fmt.Sprintf("postgres://postgres@%s/postgres?sslmode=disable", os.Getenv("DATABASE_HOST")),
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	os.Exit(m.Run())
}

func TestCreatePost(t *testing.T) {
	router := Init(mux.NewRouter(), db)

	requestPost := PostRequest{Body: "Hello, world!"}
	payload, err := json.Marshal(requestPost)
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", "/posts", bytes.NewBuffer(payload))
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var responsePost Post
	err = json.NewDecoder(w.Body).Decode(&responsePost)
	assert.Nil(t, err)
}

func TestGetPost(t *testing.T) {
	router := Init(mux.NewRouter(), db)

	responsePost, err := createPost(router, PostRequest{Body: "Hello, world!"})
	assert.Nil(t, err)

	req, err := http.NewRequest("GET", fmt.Sprintf("/post/%d", responsePost.ID), nil)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var receivedResponsePost Post
	err = json.NewDecoder(w.Body).Decode(&receivedResponsePost)
	assert.Nil(t, err)

	assert.Equal(t, responsePost, receivedResponsePost)
}

func TestDeletePost(t *testing.T) {
	router := Init(mux.NewRouter(), db)

	responsePost, err := createPost(router, PostRequest{Body: "Hello, world!"})
	assert.Nil(t, err)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/post/%d", responsePost.ID), nil)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdatePost(t *testing.T) {
	router := Init(mux.NewRouter(), db)

	responsePost, err := createPost(router, PostRequest{Body: "Hello, world!"})
	assert.Nil(t, err)

	newRequestPost := PostRequest{Body: "Hello, world!!"}
	payload, err := json.Marshal(newRequestPost)
	assert.Nil(t, err)

	req, err := http.NewRequest("PUT", fmt.Sprintf("/post/%d", responsePost.ID), bytes.NewBuffer(payload))
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var receivedResponsePost Post
	err = json.NewDecoder(w.Body).Decode(&receivedResponsePost)
	assert.Nil(t, err)

	assert.Equal(t, newRequestPost.Body, receivedResponsePost.Body)
}
