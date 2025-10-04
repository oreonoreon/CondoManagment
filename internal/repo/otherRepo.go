package repo

import "context"

type NoOpTransactionManager struct{}

func (tm *NoOpTransactionManager) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}
