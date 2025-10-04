package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"os"
)

func init() {
	DataSourceName = buildDataSourceName()
	fmt.Println(DataSourceName)
}

var defaultLink = "postgres://oreonoreon:12345@localhost:5432/postgres?sslmode=disable" // "postgres://oreonoreon:12345@postgres:5432/postgres?sslmode=disable"
var DataSourceName string

//type Repository struct {
//	PostgreSQL *standard.DB
//}
//
//func NewRepository(db *standard.DB) *Repository {
//	return &Repository{PostgreSQL: db}
//}
//
//ConnectionPostgreSQl не забываем закрыть соединение с бд, defer conn.Close()
//func ConnectionPostgreSQl() (*standard.DB, error) {
//	c, err := sql.Open("postgres", DataSourceName)
//	if err != nil {
//		return nil, err
//	}
//	db := standard.NewDB("github.com/lib/pq", c)
//
//	return db, nil
//}

type Repository struct {
	PostgreSQL *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{PostgreSQL: db}
}

func ConnectionPostgreSQl() (*sql.DB, error) {
	c, err := sql.Open("postgres", DataSourceName)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func buildDataSourceName() string {
	link, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		DB_HOST := os.Getenv("DB_HOST")
		DB_USER := os.Getenv("DB_USER")
		DB_PASS := os.Getenv("DB_PASS")
		DB_NAME := os.Getenv("DB_NAME")

		if DB_HOST != "" && DB_USER != "" && DB_PASS != "" && DB_NAME != "" {
			return "postgres://" + DB_USER + ":" + DB_PASS + "@" + DB_HOST + ":5432/" + DB_NAME + "?sslmode=disable"
		}

		return defaultLink
	} else {
		return link
	}
}

type txKey struct{}

func WithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func From(ctx context.Context) *sql.Tx {
	v := ctx.Value(txKey{})
	if v == nil {
		return nil
	}
	if tx, ok := v.(*sql.Tx); ok {
		return tx
	}
	return nil
}

// BeginTx привязывает транзакцию к контексту и возвращает новый контекст и саму транзакцию.
func (db *Repository) BeginTx(ctx context.Context, opts *sql.TxOptions) (context.Context, *sql.Tx, error) {
	tx, err := db.PostgreSQL.BeginTx(ctx, opts)
	if err != nil {
		return ctx, nil, err
	}
	return WithTx(ctx, tx), tx, nil
}

func (db *Repository) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	ctx, tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fn(ctx); err != nil {
		return err
	}

	return tx.Commit()
}

type queryer interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func getRunner(ctx context.Context, db *sql.DB) queryer {
	if tx := From(ctx); tx != nil {
		return tx
	}
	zap.L().Debug("getRunner", zap.Error(errors.New("context doesn't contain transaction")))
	return db
}
