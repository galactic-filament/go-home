package main

import (
	"fmt"
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

func main() {
	http.ListenAndServe(":80", getHandler())
}
