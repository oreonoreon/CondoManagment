package repo

import (
	"database/sql"
	"fmt"
	"github.com/gopsql/standard"
	_ "github.com/lib/pq"
	"os"
)

func init() {
	buildDataSourceName()
	fmt.Println(DataSourceName)
}

type Repository struct {
	PostgreSQL *standard.DB
}

func NewRepository(db *standard.DB) *Repository {
	return &Repository{PostgreSQL: db}
}

var DataSourceName = "postgres://oreonoreon:12345@localhost:5432/postgres?sslmode=disable" // "postgres://oreonoreon:12345@postgres:5432/postgres?sslmode=disable"

// ConnectionPostgreSQl не забываем закрыть соединение с бд, defer conn.Close()
func ConnectionPostgreSQl() (*standard.DB, error) {
	c, err := sql.Open("postgres", DataSourceName)
	if err != nil {
		return nil, err
	}
	db := standard.NewDB("github.com/lib/pq", c)

	return db, nil
}

func buildDataSourceName() string {
	DB_HOST := os.Getenv("DB_HOST")
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASS")
	DB_NAME := os.Getenv("DB_NAME")
	if DB_HOST != "" && DB_USER != "" && DB_PASS != "" && DB_NAME != "" {
		DataSourceName = "postgres://" + DB_USER + ":" + DB_PASS + "@" + DB_HOST + ":5432/" + DB_NAME + "?sslmode=disable"
	}
	return DataSourceName
}
