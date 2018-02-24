package routehandler

import (
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetHandler(t *testing.T) {
	db, err := sqlx.Connect(
		"postgres",
		fmt.Sprintf("postgres://postgres@%s/postgres?sslmode=disable", os.Getenv("DATABASE_HOST")),
	)
	assert.Nil(t, err)
	GetHandler(db)
}
