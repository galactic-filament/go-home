package PostManager

import (
	"github.com/jmoiron/sqlx"
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

func (m manager) create(pr PostRequest) (p Post, err error) {
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

func (m manager) update(p Post, pr PostRequest) (Post, error) {
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
