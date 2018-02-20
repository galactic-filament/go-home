package routehandler

import (
	"github.com/galactic-filament/go-home/app/defaultmanager"
	"github.com/galactic-filament/go-home/app/postmanager"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// GetHandler - generates this app's route handler
func GetHandler(db *sqlx.DB) *mux.Router {
	r := mux.NewRouter()

	// route handlers
	r = postmanager.Init(r, db)
	r = defaultmanager.Init(r)
	return r
}
