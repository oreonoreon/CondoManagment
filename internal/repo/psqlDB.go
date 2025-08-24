package repo

import (
	"database/sql"
	"github.com/gopsql/standard"
	_ "github.com/lib/pq"
)

type Repository struct {
	PostgreSQL *standard.DB
}

func NewRepository(db *standard.DB) *Repository {
	return &Repository{PostgreSQL: db}
}

var DataSourceName = "postgres://oreonoreon:12345@postgres:5432/postgres?sslmode=disable" // "postgres://oreonoreon:12345@localhost:5432/postgres?sslmode=disable"

// ConnectionPostgreSQl не забываем закрыть соединение с бд, defer conn.Close()
func ConnectionPostgreSQl() (*standard.DB, error) {
	c, err := sql.Open("postgres", DataSourceName)
	if err != nil {
		return nil, err
	}
	db := standard.NewDB("github.com/lib/pq", c)

	return db, nil
}
