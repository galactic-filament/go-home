package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
)

func add(x int, y int) int {
	return x + y
}

func getHandler() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
	return r
}

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.WithFields(log.Fields{
			"url": req.URL,
		}).Info("Url hit")
		h.ServeHTTP(w, req)
	})
}

func main() {
	log.Info("Starting up")
	http.ListenAndServe(":80", loggingMiddleware(getHandler()))
}
