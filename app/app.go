package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/ihsw/go-home/app/PostManager"
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

	// misc route endpoints
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

	// post handler
	r = PostManager.Init(r, db)
	return r
}

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// read the request body for logging
		var (
			body []byte
			err  error
		)
		if body, err = ioutil.ReadAll(req.Body); err != nil {
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

		// logging the request body
		log.WithFields(log.Fields{
			"url":  req.URL,
			"body": string(body),
		}).Info("Url hit")

		// passing onto the next middleware
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
