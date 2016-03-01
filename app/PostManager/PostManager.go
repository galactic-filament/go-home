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

func (m manager) delete(p Post) (err error) {
	stmt, err := m.db.Prepare("DELETE FROM posts WHERE id = $1")
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(p.ID); err != nil {
		return err
	}

	return nil
}

func (m manager) update(p Post, pr postRequest) (Post, error) {
	stmt, err := m.db.Prepare("UPDATE posts SET body = $1 WHERE id = $2")
	if err != nil {
		return Post{}, err
	}

	if _, err := stmt.Exec(pr.Body, p.ID); err != nil {
		return Post{}, err
	}

	p.Body = pr.Body

	return p, nil
}

type postRequest struct {
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
		var request postRequest
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
		var request postRequest
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
