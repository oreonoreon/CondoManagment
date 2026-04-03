package repo

import (
	"awesomeProject/internal/config"
	"context"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Repository struct {
	PostgreSQL *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{PostgreSQL: db}
}

func ConnectionPostgreSQl(env *config.ConfigEnv) (*sql.DB, error) {
	c, err := sql.Open("postgres", env.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	return c, nil
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
