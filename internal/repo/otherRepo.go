package repo

import (
	"context"
	"github.com/google/uuid"
)

type NoOpTransactionManager struct{}

func (tm *NoOpTransactionManager) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
