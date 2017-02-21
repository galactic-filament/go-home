package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/ihsw/go-home/app/RouteHandler"
	"github.com/ihsw/go-home/app/Util"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type customReader struct {
	*bytes.Buffer
}

func (r customReader) Close() error { return nil }

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
	envVars, errorMessages := Util.ValidateEnvironment()
	if errorMessages != nil {
		for _, message := range errorMessages {
			fmt.Println(message)
		}

		os.Exit(1)
		return
	}

	log.Info("Starting up")
	db, err := sqlx.Connect(
		"postgres",
		fmt.Sprintf("postgres://postgres@%s/postgres?sslmode=disable", envVars["DATABASE_HOST"]),
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = http.ListenAndServe(":80", loggingMiddleware(corsMiddleware(RouteHandler.GetHandler(db))))
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Could not listen on 80")
	}
}
