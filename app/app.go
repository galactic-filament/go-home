package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/galactic-filament/go-home/app/routehandler"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func readBody(req *http.Request) ([]byte, error) {
	if req.Body == nil {
		return []byte{}, nil
	}

	return ioutil.ReadAll(req.Body)
}

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// read the request body for logging
		body, err := readBody(req)
		if err != nil {
			log.WithFields(log.Fields{
				"url": req.URL,
				"err": err.Error(),
			}).Error("Could not read request body")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Could not read request body")
			return
		}

		// optionally re-adding the request body
		if len(body) > 0 {
			req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}

		// logging the request body
		log.WithFields(log.Fields{
			"url":  req.URL,
			"body": string(body),
		}).Info("Url hit")

		// passing onto the next middleware
		h.ServeHTTP(w, req)
	})
}

func corsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if req.Method == "OPTIONS" {
			return
		}

		// passing onto the next middleware
		h.ServeHTTP(w, req)
	})
}

func main() {
	log.Info("Starting up")
	db, err := sqlx.Connect(
		"postgres",
		fmt.Sprintf("postgres://postgres@%s/postgres?sslmode=disable", os.Getenv("DATABASE_HOST")),
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = http.ListenAndServe(":80", loggingMiddleware(corsMiddleware(routehandler.GetHandler(db))))
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Could not listen on 80")
	}
}
