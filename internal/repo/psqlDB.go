package repo

import (
	"database/sql"
	"github.com/gopsql/standard"
	_ "github.com/lib/pq"
)

type DBPostgreSQl struct {
	PostgreSQL *standard.DB
}

// DataBasePostgreSQl не забываем закрыть соединение с бд, defer conn.Close()
func DataBasePostgreSQl() (*DBPostgreSQl, error) {
	c, err := sql.Open("postgres", "postgres://oreonoreon:12345@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		return nil, err
	}
	conn := standard.NewDB("github.com/lib/pq", c)

	return &DBPostgreSQl{PostgreSQL: conn}, nil
}
