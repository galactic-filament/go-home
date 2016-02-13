package PostManager

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/ihsw/go-home/app/Util"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
)

// Post - entity
type Post struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type manager struct {
	db *sqlx.DB
}

func newManager(db *sqlx.DB) manager {
	return manager{db}
}

func (m manager) create(pr postRequest) (p Post, err error) {
	stmt, err := m.db.PrepareNamed("INSERT INTO posts (body) VALUES (:body) RETURNING id")
	if err != nil {
		return Post{}, err
	}

	row := stmt.QueryRow(pr)
	var id int
	if err := row.Scan(&id); err != nil {
		return Post{}, err
	}

	return Post{ID: id, Body: pr.Body}, nil
}

func (m manager) get(id int) (p Post, err error) {
	stmt, err := m.db.Preparex("SELECT id, body FROM posts WHERE id = $1")
	if err != nil {
		return Post{}, err
	}

	err = stmt.Get(&p, id)
	if err != nil {
		return Post{}, err
	}

	return p, nil
}

type postRequest struct {
	Body string `json:"body"`
}

// Init - route handler
func Init(r *mux.Router, db *sqlx.DB) *mux.Router {
	m := newManager(db)

	r.HandleFunc("/posts", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// decoding the request body
		var (
			pr  postRequest
			err error
		)
		if err = json.NewDecoder(req.Body).Decode(&pr); err != nil {
			Util.WriteJSONErrorResponse(w, err)
			return
		}

		// creating the post
		var p Post
		if p, err = m.create(pr); err != nil {
			Util.WriteJSONErrorResponse(w, err)
			return
		}

		// writing out the response
		if err := json.NewEncoder(w).Encode(p); err != nil {
			Util.WriteJSONErrorResponse(w, err)
			return
		}
	}).Methods("POST")

	r.HandleFunc("/post/{id:[0-9]+}", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-type", "application/json")

		// fetching the url vars
		vars := mux.Vars(req)
		var (
			id  int
			err error
		)
		if id, err = strconv.Atoi(vars["id"]); err != nil {
			Util.WriteJSONErrorResponse(w, err)
			return
		}

		// getting the post
		var p Post
		if p, err = m.get(id); err != nil {
			Util.WriteJSONErrorResponse(w, err)
			return
		}

		// writing out the response
		if err = json.NewEncoder(w).Encode(p); err != nil {
			Util.WriteJSONErrorResponse(w, err)
			return
		}
	}).Methods("GET")

	return r
}