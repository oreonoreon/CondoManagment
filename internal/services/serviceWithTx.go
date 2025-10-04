package services

import (
	"awesomeProject/internal/entities"
	"context"
)

type TransactionalService struct {
	serviceInterface ServiceInterface
	txManager        TransactionManager
}

type ServiceInterface interface {
	UpdateBooking(ctx context.Context, booking entities.Booking) (*entities.Booking, error)
	CreateBooking(ctx context.Context, booking entities.Booking) (*entities.Booking, error)
}

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(context.Context) error) error
}

func NewTransactionalService(serviceInterface ServiceInterface, txManager TransactionManager) *TransactionalService {
	return &TransactionalService{
		serviceInterface: serviceInterface,
		txManager:        txManager,
	}
}

func (ts *TransactionalService) UpdateBooking(ctx context.Context, booking entities.Booking) (*entities.Booking, error) {
	var result *entities.Booking
	var resultErr error

	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.serviceInterface.UpdateBooking(ctx, booking)
		return resultErr
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *TransactionalService) CreateBooking(ctx context.Context, booking entities.Booking) (*entities.Booking, error) {
	var result *entities.Booking
	var resultErr error

	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.serviceInterface.CreateBooking(ctx, booking)
		return resultErr
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}
