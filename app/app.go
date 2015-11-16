package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type greeting struct {
	Greeting string `json:"greeting"`
}

type customReader struct {
	*bytes.Buffer
}

func (r customReader) Close() error { return nil }

func getHandler() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
	r.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Pong")
	})
	r.HandleFunc("/reflection", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-type", "application/json")

		var greeting greeting
		if err := json.NewDecoder(req.Body).Decode(&greeting); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Could not decode request body")
			return
		}

		if err := json.NewEncoder(w).Encode(greeting); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Could not encode response body")
			return
		}
	})
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

	err := http.ListenAndServe(":80", loggingMiddleware(getHandler()))
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Could not listen on 80")
	}
}
