package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
)

type greeting struct {
	Greeting string `json:"greeting"`
}

type postRequest struct {
	Body string `json:"body"`
}

type post struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type customReader struct {
	*bytes.Buffer
}

type errorResponse struct {
	Error string `json:"error"`
}

func (r customReader) Close() error { return nil }

func writeJSONErrorResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	errResponse := errorResponse{Error: err.Error()}
	if err := json.NewEncoder(w).Encode(errResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not encode error response body")
		return
	}
	return
}

func getHandler(db *sqlx.DB) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
	r.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Pong")
	})
	r.HandleFunc("/reflection", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// decoding the request body
		var greeting greeting
		if err := json.NewDecoder(req.Body).Decode(&greeting); err != nil {
			writeJSONErrorResponse(w, err)
			return
		}

		// writing out the response
		if err := json.NewEncoder(w).Encode(greeting); err != nil {
			writeJSONErrorResponse(w, err)
			return
		}
	}).Methods("POST")
	r.HandleFunc("/posts", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// decoding the request body
		var postRequest postRequest
		if err := json.NewDecoder(req.Body).Decode(&postRequest); err != nil {
			writeJSONErrorResponse(w, err)
			return
		}

		// inserting the post and fetching the resulting id
		stmt, err := db.PrepareNamed("INSERT INTO posts (body) VALUES (:body) RETURNING id")
		if err != nil {
			writeJSONErrorResponse(w, err)
			return
		}
		row := stmt.QueryRow(postRequest)
		var id int
		if err := row.Scan(&id); err != nil {
			writeJSONErrorResponse(w, err)
			return
		}

		// writing out the response
		responseBody := post{ID: id, Body: postRequest.Body}
		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
			writeJSONErrorResponse(w, err)
			return
		}
	}).Methods("POST")
	r.HandleFunc("/post/{id:[0-9]+}", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// fetching the url vars
		vars := mux.Vars(req)
		id := vars["id"]

		// query for the post
		stmt, err := db.Preparex("SELECT id, body FROM posts WHERE id = $1")
		if err != nil {
			writeJSONErrorResponse(w, err)
			return
		}
		var post post
		err = stmt.Get(&post, id)
		if err != nil {
			writeJSONErrorResponse(w, err)
			return
		}

		// writing out the response
		if err := json.NewEncoder(w).Encode(post); err != nil {
			writeJSONErrorResponse(w, err)
			return
		}
	}).Methods("GET")
	return r
}

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// read the request body for logging
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.WithFields(log.Fields{
				"url": req.URL,
				"err": err.Error(),
			}).Error("Could not read request body")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Could not read request body")
			return
		}

		// re-adding the request body
		req.Body = customReader{bytes.NewBuffer(body)}

		// passing onto the next middleware
		log.WithFields(log.Fields{
			"url":  req.URL,
			"body": string(body),
		}).Info("Url hit")
		h.ServeHTTP(w, req)
	})
}

func main() {
	log.Info("Starting up")

	db, err := sqlx.Connect(
		"postgres",
		"postgres://postgres@db/postgres?sslmode=disable",
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = http.ListenAndServe(":80", loggingMiddleware(getHandler(db)))
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Could not listen on 80")
	}
}
