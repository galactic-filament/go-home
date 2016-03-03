package DefaultManager

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ihsw/go-home/app/Util"
	"net/http"
)

// GreetingRequest - reflection request body
type GreetingRequest struct {
	Greeting string `json:"greeting"`
}

// Init - route handler
func Init(r *mux.Router) *mux.Router {
	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
	r.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Pong")
	})
	r.HandleFunc("/reflection", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// decoding the request body
		var greeting GreetingRequest
		if err := json.NewDecoder(req.Body).Decode(&greeting); err != nil {
			Util.WriteJSONErrorResponse(w, err)
			return
		}

		// writing out the response
		if err := json.NewEncoder(w).Encode(greeting); err != nil {
			Util.WriteJSONErrorResponse(w, err)
			return
		}
	}).Methods("POST")

	return r
}
