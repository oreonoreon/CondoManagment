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
	DeleteReservation(ctx context.Context, id int) (*entities.Reservation, error)
	CreateReservation(ctx context.Context, reservation entities.Reservation) (*entities.Reservation, error)
	CreateGuest(ctx context.Context, g entities.Guest) (*entities.Guest, error)
	CreateReport(ctx context.Context, roomNumber string, startPeriod string, endPeriod string) ([]byte, error)
	GetBookingALLForApartmentALL(ctx context.Context, roomNumbers []string) ([]entities.Booking, error)
	GetBookingALLForApartment(ctx context.Context, roomNumber string) ([]entities.Booking, error)
	GetReservationALLForApartment(ctx context.Context, roomNumber string) ([]entities.Reservation, error)
	GetBooking(ctx context.Context, roomNumber string, start string, end string) ([]entities.Booking, error)
	GetReservationForPeriodByApartment(ctx context.Context, roomNumber string, start string, end string) ([]entities.Reservation, error)
	GetReservationByPhoneNumber(ctx context.Context, phone string) ([]entities.Reservation, error)
	FindTotalPriceForPeriodReport(ctx context.Context, apartments []entities.Apartment, startPeriod, endPeriod string) (map[string]int, error)
	FindMiddlePriceForPeriodReport(ctx context.Context, apartments []entities.Apartment, startPeriod, endPeriod string) (map[string]int, error)
	FindMiddlePriceForPeriod(ctx context.Context, roomNumber string, startPeriod, endPeriod string) (int, error)
	FindTotalPriceForPeriod(ctx context.Context, roomNumber, startPeriod, endPeriod string) (int, int, error)
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

func (ts *TransactionalService) DeleteReservation(ctx context.Context, id int) (*entities.Reservation, error) {
	var result *entities.Reservation
	var resultErr error

	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.serviceInterface.DeleteReservation(ctx, id)
		return resultErr
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *TransactionalService) CreateReservation(ctx context.Context, reservation entities.Reservation) (*entities.Reservation, error) {
	var result *entities.Reservation
	var resultErr error

	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.serviceInterface.CreateReservation(ctx, reservation)
		return resultErr
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *TransactionalService) CreateGuest(ctx context.Context, g entities.Guest) (*entities.Guest, error) {
	var result *entities.Guest
	var resultErr error

	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.serviceInterface.CreateGuest(ctx, g)
		return resultErr
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *TransactionalService) CreateReport(ctx context.Context, roomNumber string, startPeriod string, endPeriod string) ([]byte, error) {
	var result []byte
	var resultErr error

	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.serviceInterface.CreateReport(ctx, roomNumber, startPeriod, endPeriod)
		return resultErr
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *TransactionalService) GetBookingALLForApartmentALL(ctx context.Context, roomNumbers []string) ([]entities.Booking, error) {
	var result []entities.Booking
	var resultErr error

	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.serviceInterface.GetBookingALLForApartmentALL(ctx, roomNumbers)
		return resultErr
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *TransactionalService) GetBookingALLForApartment(ctx context.Context, roomNumber string) ([]entities.Booking, error) {
	var result []entities.Booking
	var resultErr error

	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.serviceInterface.GetBookingALLForApartment(ctx, roomNumber)
		return resultErr
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *TransactionalService) GetReservationALLForApartment(ctx context.Context, roomNumber string) ([]entities.Reservation, error) {
	var result []entities.Reservation
	var resultErr error

	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.serviceInterface.GetReservationALLForApartment(ctx, roomNumber)
		return resultErr
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *TransactionalService) GetBooking(ctx context.Context, roomNumber string, start string, end string) ([]entities.Booking, error) {
	var result []entities.Booking
	var resultErr error

	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.serviceInterface.GetBooking(ctx, roomNumber, start, end)
		return resultErr
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *TransactionalService) GetReservationForPeriodByApartment(ctx context.Context, roomNumber string, start string, end string) ([]entities.Reservation, error) {
	var result []entities.Reservation
	var resultErr error

	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.serviceInterface.GetReservationForPeriodByApartment(ctx, roomNumber, start, end)
		return resultErr
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *TransactionalService) GetReservationByPhoneNumber(ctx context.Context, phone string) ([]entities.Reservation, error) {
	var result []entities.Reservation
	var resultErr error

	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.serviceInterface.GetReservationByPhoneNumber(ctx, phone)
		return resultErr
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *TransactionalService) FindTotalPriceForPeriodReport(ctx context.Context, apartments []entities.Apartment, startPeriod, endPeriod string) (map[string]int, error) {
	var result map[string]int
	var resultErr error

	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.serviceInterface.FindTotalPriceForPeriodReport(ctx, apartments, startPeriod, endPeriod)
		return resultErr
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *TransactionalService) FindMiddlePriceForPeriodReport(ctx context.Context, apartments []entities.Apartment, startPeriod, endPeriod string) (map[string]int, error) {
	var result map[string]int
	var resultErr error

	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.serviceInterface.FindMiddlePriceForPeriodReport(ctx, apartments, startPeriod, endPeriod)
		return resultErr
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *TransactionalService) FindMiddlePriceForPeriod(ctx context.Context, roomNumber string, startPeriod, endPeriod string) (int, error) {
	var result int
	var resultErr error

	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.serviceInterface.FindMiddlePriceForPeriod(ctx, roomNumber, startPeriod, endPeriod)
		return resultErr
	})

	if err != nil {
		return 0, err
	}
	return result, nil
}

func (ts *TransactionalService) FindTotalPriceForPeriod(ctx context.Context, roomNumber, startPeriod, endPeriod string) (int, int, error) {
	var result1 int
	var result2 int
	var resultErr error

	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result1, result2, resultErr = ts.serviceInterface.FindTotalPriceForPeriod(ctx, roomNumber, startPeriod, endPeriod)
		return resultErr
	})

	if err != nil {
		return 0, 0, err
	}
	return result1, result2, nil
}
