package services

import (
	"awesomeProject/internal/entities"
	"context"
	"time"
)

type TransactionalService struct {
	serviceInterface            ServiceInterface
	cleaningServiceInterface    CleaningServiceInterface
	reservationInfoSvcInterface ReservationInfoServiceInterface
	txManager                   TransactionManager
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
	FindTotalPriceForPeriodReport(ctx context.Context, apartments []entities.Apartment, start, end time.Time) (map[string]int, error)
	FindMiddlePriceForPeriodReport(ctx context.Context, apartments []entities.Apartment, start, end time.Time) (map[string]int, error)
	FindMiddlePriceForPeriod(ctx context.Context, roomNumber string, start, end time.Time) (int, error)
	FindTotalPriceForPeriod(ctx context.Context, roomNumber string, start, end time.Time) (int, int, error)
	GetBookingByCheckIn(ctx context.Context, date time.Time) ([]entities.Booking, error)
	GetBookingByCheckOut(ctx context.Context, date time.Time) ([]entities.Booking, error)
	GetBookingsForRooms(ctx context.Context, roomNumbers []string) ([]entities.Booking, error)
}

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(context.Context) error) error
}

type ReservationInfoServiceInterface interface {
	UpdateReservationInfoByID(ctx context.Context, ri entities.ReservationInfo) (*entities.ReservationInfo, error)
}

type CleaningServiceInterface interface {
	CreateCleaningManual(ctx context.Context, c entities.Cleaning) (*entities.Cleaning, error)
	UpdateCleaning(ctx context.Context, c entities.Cleaning) (*entities.Cleaning, error)
	DeleteCleaning(ctx context.Context, id int) (*entities.Cleaning, error)
	GetCleaningByID(ctx context.Context, id int) (*entities.Cleaning, error)
	GetAllCleaning(ctx context.Context) ([]entities.Cleaning, error)
	GetCleaningByDate(ctx context.Context, date time.Time) ([]entities.Cleaning, error)
}

func NewTransactionalService(serviceInterface ServiceInterface, cleaningServiceInterface CleaningServiceInterface, reservationInfoSvc ReservationInfoServiceInterface, txManager TransactionManager) *TransactionalService {
	return &TransactionalService{
		serviceInterface:            serviceInterface,
		cleaningServiceInterface:    cleaningServiceInterface,
		reservationInfoSvcInterface: reservationInfoSvc,
		txManager:                   txManager,
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

func (ts *TransactionalService) FindTotalPriceForPeriodReport(ctx context.Context, apartments []entities.Apartment, start, end time.Time) (map[string]int, error) {
	var result map[string]int
	var resultErr error

	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.serviceInterface.FindTotalPriceForPeriodReport(ctx, apartments, start, end)
		return resultErr
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *TransactionalService) FindMiddlePriceForPeriodReport(ctx context.Context, apartments []entities.Apartment, start, end time.Time) (map[string]int, error) {
	var result map[string]int
	var resultErr error

	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.serviceInterface.FindMiddlePriceForPeriodReport(ctx, apartments, start, end)
		return resultErr
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *TransactionalService) FindMiddlePriceForPeriod(ctx context.Context, roomNumber string, start, end time.Time) (int, error) {
	var result int
	var resultErr error

	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.serviceInterface.FindMiddlePriceForPeriod(ctx, roomNumber, start, end)
		return resultErr
	})

	if err != nil {
		return 0, err
	}
	return result, nil
}

func (ts *TransactionalService) FindTotalPriceForPeriod(ctx context.Context, roomNumber string, start, end time.Time) (int, int, error) {
	var result1 int
	var result2 int
	var resultErr error

	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result1, result2, resultErr = ts.serviceInterface.FindTotalPriceForPeriod(ctx, roomNumber, start, end)
		return resultErr
	})

	if err != nil {
		return 0, 0, err
	}
	return result1, result2, nil
}

func (ts *TransactionalService) GetBookingByCheckIn(ctx context.Context, date time.Time) ([]entities.Booking, error) {
	var result []entities.Booking
	var resultErr error
	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.serviceInterface.GetBookingByCheckIn(ctx, date)
		return resultErr
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *TransactionalService) GetBookingByCheckOut(ctx context.Context, date time.Time) ([]entities.Booking, error) {
	var result []entities.Booking
	var resultErr error
	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.serviceInterface.GetBookingByCheckOut(ctx, date)
		return resultErr
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *TransactionalService) GetBookingsForRooms(ctx context.Context, roomNumbers []string) ([]entities.Booking, error) {
	var result []entities.Booking
	var resultErr error
	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.serviceInterface.GetBookingsForRooms(ctx, roomNumbers)
		return resultErr
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ------- Cleaning (через транзакцию) -------

func (ts *TransactionalService) CreateCleaningManual(ctx context.Context, c entities.Cleaning) (*entities.Cleaning, error) {
	var result *entities.Cleaning
	var resultErr error
	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.cleaningServiceInterface.CreateCleaningManual(ctx, c)
		return resultErr
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *TransactionalService) UpdateCleaning(ctx context.Context, c entities.Cleaning) (*entities.Cleaning, error) {
	var result *entities.Cleaning
	var resultErr error
	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.cleaningServiceInterface.UpdateCleaning(ctx, c)
		return resultErr
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *TransactionalService) DeleteCleaning(ctx context.Context, id int) (*entities.Cleaning, error) {
	var result *entities.Cleaning
	var resultErr error
	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.cleaningServiceInterface.DeleteCleaning(ctx, id)
		return resultErr
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *TransactionalService) GetCleaningByID(ctx context.Context, id int) (*entities.Cleaning, error) {
	var result *entities.Cleaning
	var resultErr error
	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.cleaningServiceInterface.GetCleaningByID(ctx, id)
		return resultErr
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *TransactionalService) GetAllCleaning(ctx context.Context) ([]entities.Cleaning, error) {
	var result []entities.Cleaning
	var resultErr error
	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.cleaningServiceInterface.GetAllCleaning(ctx)
		return resultErr
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *TransactionalService) GetCleaningByDate(ctx context.Context, date time.Time) ([]entities.Cleaning, error) {
	var result []entities.Cleaning
	var resultErr error
	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.cleaningServiceInterface.GetCleaningByDate(ctx, date)
		return resultErr
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ------- ReservationInfo (через транзакцию) -------

func (ts *TransactionalService) UpdateReservationInfoByID(ctx context.Context, ri entities.ReservationInfo) (*entities.ReservationInfo, error) {
	var result *entities.ReservationInfo
	var resultErr error
	err := ts.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		result, resultErr = ts.reservationInfoSvcInterface.UpdateReservationInfoByID(ctx, ri)
		return resultErr
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
