package RouteHandler

import (
	"github.com/gorilla/mux"
	"github.com/ihsw/go-home/app/DefaultManager"
	"github.com/ihsw/go-home/app/PostManager"
	"github.com/jmoiron/sqlx"
)

// GetHandler - generates this app's route handler
func GetHandler(db *sqlx.DB) *mux.Router {
	r := mux.NewRouter()

	// route handlers
	r = PostManager.Init(r, db)
	r = DefaultManager.Init(r)
	return r
}
