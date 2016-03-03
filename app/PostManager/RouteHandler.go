package PostManager

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/ihsw/go-home/app/Util"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
)

// PostRequest - post request body
type PostRequest struct {
	Body string `json:"body"`
}

// DeleteResponse - delete post response body
type DeleteResponse struct{}

// Init - route handler
func Init(r *mux.Router, db *sqlx.DB) *mux.Router {
	postManager := newManager(db)

	r.HandleFunc("/posts", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// decoding the request body
		var request PostRequest
		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil {
			Util.WriteJSONErrorResponse(w, err)
			return
		}

		// creating the post
		post, err := postManager.create(request)
		if err != nil {
			Util.WriteJSONErrorResponse(w, err)
			return
		}

		// writing out the response
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(post)
		if err != nil {
			Util.WriteJSONErrorResponse(w, err)
			return
		}
	}).Methods("POST")

	r.HandleFunc("/post/{id:[0-9]+}", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// fetching the url vars
		id, err := strconv.Atoi(mux.Vars(req)["id"])
		if err != nil {
			Util.WriteJSONErrorResponse(w, err)
			return
		}

		// getting the post
		post, err := postManager.get(id)
		if err != nil {
			Util.WriteJSONErrorResponse(w, err)
			return
		}

		// writing out the response
		err = json.NewEncoder(w).Encode(post)
		if err != nil {
			Util.WriteJSONErrorResponse(w, err)
			return
		}
	}).Methods("GET")
	r.HandleFunc("/post/{id:[0-9]+}", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// fetching the url vars
		id, err := strconv.Atoi(mux.Vars(req)["id"])
		if err != nil {
			Util.WriteJSONErrorResponse(w, err)
			return
		}

		// getting the post
		post, err := postManager.get(id)
		if err != nil {
			Util.WriteJSONErrorResponse(w, err)
			return
		}

		// deleting the post
		err = postManager.delete(post)
		if err != nil {
			Util.WriteJSONErrorResponse(w, err)
		}

		// writing out the response
		err = json.NewEncoder(w).Encode(DeleteResponse{})
		if err != nil {
			Util.WriteJSONErrorResponse(w, err)
			return
		}
	}).Methods("DELETE")
	r.HandleFunc("/post/{id:[0-9]+}", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// fetching the url vars
		id, err := strconv.Atoi(mux.Vars(req)["id"])
		if err != nil {
			Util.WriteJSONErrorResponse(w, err)
			return
		}

		// getting the post
		post, err := postManager.get(id)
		if err != nil {
			Util.WriteJSONErrorResponse(w, err)
			return
		}

		// decoding the request body
		var request PostRequest
		err = json.NewDecoder(req.Body).Decode(&request)
		if err != nil {
			Util.WriteJSONErrorResponse(w, err)
			return
		}

		// updating the post
		post, err = postManager.update(post, request)
		if err != nil {
			Util.WriteJSONErrorResponse(w, err)
		}

		// writing out the response
		err = json.NewEncoder(w).Encode(post)
		if err != nil {
			Util.WriteJSONErrorResponse(w, err)
			return
		}
	}).Methods("PUT")

	return r
}
