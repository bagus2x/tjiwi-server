package db

import (
	"database/sql"

	"github.com/bagus2x/tjiwi/config"
	_ "github.com/lib/pq"
)

func NewPostgresDatabase(c *config.Config) *sql.DB {
	db, err := sql.Open("postgres", c.DatabaseConnection())
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}
