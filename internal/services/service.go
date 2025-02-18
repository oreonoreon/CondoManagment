package services

import (
	"awesomeProject/internal/entities"
	"awesomeProject/internal/erro"
	"awesomeProject/internal/models"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type Service struct {
	storageReservation StorageReservation
	storageGuest       StorageGuest
}

type StorageReservation interface {
	Create(ctx context.Context, r entities.Reservation) (*entities.Reservation, error)
	ReadWithRoomNumber(ctx context.Context, roomNumber string, checkin, checkout time.Time) ([]entities.Reservation, error)
}

type StorageGuest interface {
	CreateGuest(ctx context.Context, g entities.Guest) (*entities.Guest, error)
	FindGuestByPhoneNumber(ctx context.Context, phone string) (*entities.Guest, error)
}

func NewService(storage StorageReservation, storageGuest StorageGuest) *Service {
	return &Service{
		storageReservation: storage,
		storageGuest:       storageGuest,
	}
}

func (s *Service) CreateReservation(ctx context.Context, booking models.BookingInfo) error {
	g, err := booking.DbModelConvertGuest()
	if err != nil {
		return err
	}
	zap.L().Debug("CreateReservation", zap.Any("entities.Guest", *g))

	guest, err := s.storageGuest.FindGuestByPhoneNumber(ctx, booking.Phone)
	if err != nil {
		return err
	}
	if guest == nil {
		guest, err = s.storageGuest.CreateGuest(ctx, *g)
		if err != nil {
			return err
		}
	}

	reserv, err := booking.DbModelConvert(guest.GuestID)
	if err != nil {
		return err
	}
	zap.L().Debug("CreateReservation", zap.Any("entities.Reservation", *reserv))

	//чекнем пересечение по датам на эту комнату
	values, err := s.storageReservation.ReadWithRoomNumber(ctx, reserv.RoomNumber, reserv.CheckIn, reserv.CheckOut)
	if err != nil {
		return err
	}
	if len(values) > 0 {
		if len(values) == 1 {
			if values[0].CheckIn.UTC() == reserv.CheckIn &&
				values[0].CheckOut.UTC() == reserv.CheckOut &&
				values[0].Price == reserv.Price &&
				values[0].RoomNumber == reserv.RoomNumber {
				return fmt.Errorf("%w\nзаписываемое значение: %v;\nзначение из БД: %v;", erro.ErrFullyMatchOtherBooking, booking, values[0])
			}
		}
		return errors.New(fmt.Sprintf("Букинг: %v ; пересекается со следуюшим бронированиями: %v ;", booking, values))
	}

	//запишем новое бронирование в бд
	_, err = s.storageReservation.Create(ctx, *reserv)
	if err != nil {
		return err
	}
	return nil
}

//func (s *Service) GetBookingForPeriodByApartment(ctx context.Context, roomNumber string, start string, end string) error {
//	checkin, err := models.TimeConvert(start)
//	if err != nil {
//		return err
//	}
//
//	checkout, err := models.TimeConvert(end)
//	if err != nil {
//		return err
//	}
//
//	bookings, err := s.storageReservation.ReadWithRoomNumber(ctx, roomNumber, checkin, checkout)
//	if err != nil {
//		return err
//	}
//
//}

// Создание репорта для собственика
func (s *Service) CreateReport(ctx context.Context, roomNumber string, startPeriod string, endPeriod string) ([]entities.Reservation, error) {
	return nil, nil
}

// Оповещение собственика о предстоящем бронирование
func FutureBooking() {

}

// свободные квартиры на эти даты
func FreeApartmentForDates() {

}

// цены расчитаные по переуду, к примеру низкий сезон переходящий в высокий
func AproksimatePrice() {

}
