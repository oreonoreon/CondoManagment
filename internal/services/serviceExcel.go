package services

import (
	"awesomeProject/internal/entities"
	"awesomeProject/internal/erro"
	"awesomeProject/internal/excelCalendarScraper"
	"awesomeProject/internal/excelCalendarScraper/models"
	"context"
	"fmt"
	"go.uber.org/zap"
)

type ServiceExcel struct {
	storageSettings StorageSettings
	service         *Service
}

func NewServiceExcel(storageSettings StorageSettings, service *Service) *ServiceExcel {
	return &ServiceExcel{
		storageSettings: storageSettings,
		service:         service,
	}
}

func (e *ServiceExcel) Sync(ctx context.Context, sheetName string, roomNumber string, start string, end string) ([]entities.Booking, []entities.Booking, error) {
	bookings, err := e.GetBookingForPeriod(ctx, sheetName, roomNumber, start, end)
	if err != nil {
		zap.L().Error("ServiceExcel.GetBookingForPeriod", zap.Error(err))
		return nil, nil, err
	}
	if len(bookings) == 0 {
		return nil, nil, erro.ErrNoFoundBookingsFromExcel
	}

	bookingsFromDB, err := e.service.GetBooking(ctx, roomNumber, start, end)
	if err != nil {
		zap.L().Error("Service.GetBooking", zap.Error(err))
		return nil, nil, err
	}

	for i := 0; i < len(bookingsFromDB); i++ {
		var equal bool
		for j := 0; j < len(bookings); j++ {
			if BookingDbEquilToBooking(bookings[j], bookingsFromDB[i]) {
				bookings = removeElement(bookings, j)
				equal = true
				j--
			}
		}
		if !equal {
			_, err := e.service.storageReservation.Delete(ctx, bookingsFromDB[i].Oid)
			if err != nil {
				zap.L().Debug("Sync/Delete", zap.Any("e.service.storageReservation.Delete(ctx, bookingsFromDB[i].Oid)", bookingsFromDB[i]))
				return nil, nil, err
			}
			bookingsFromDB = removeElement(bookingsFromDB, i)
			i--
		}
	}

	for _, booking := range bookings {
		err = e.service.CreateReservation1(ctx, booking)
		if err != nil {
			zap.L().Error("Sync", zap.Error(err))
			return nil, nil, err
		}
	}

	return bookings, bookingsFromDB, nil
}

//---------------------------------------------------------------------------------------------------

// SynchronizeSliceFromDBandExcel Можно заменить часть кода в Sync на этот метод. Изначально метод был создан для простоты тестирования.
func (e *ServiceExcel) SynchronizeSliceFromDBandExcel(ctx context.Context, bookingsFromDB, bookingsFromExcel []entities.Booking) ([]entities.Booking, []entities.Booking, error) {
	for i := 0; i < len(bookingsFromDB); i++ {
		var equal bool
		for j := 0; j < len(bookingsFromExcel); j++ {
			if BookingDbEquilToBooking(bookingsFromExcel[j], bookingsFromDB[i]) {
				bookingsFromExcel = removeElement(bookingsFromExcel, j)
				equal = true
				j--
			}
		}
		if !equal {
			_, err := e.service.storageReservation.Delete(ctx, bookingsFromDB[i].Oid)
			if err != nil {
				zap.L().Debug("Sync/Delete", zap.Any("e.service.storageReservation.Delete(ctx, bookingsFromDB[i].Oid)", bookingsFromDB[i]))
				return nil, nil, err
			}
			bookingsFromDB = removeElement(bookingsFromDB, i)
			i--
		}
	}
	return bookingsFromExcel, bookingsFromDB, nil
}

//-------------------------------------------------------------------------------------------------

func removeElement(bookings []entities.Booking, j int) []entities.Booking {
	// 1. Копировать последний элемент в индекс j.
	bookings[j] = bookings[len(bookings)-1]

	// 2. Удалить последний элемент (записать нулевое значение).
	bookings[len(bookings)-1] = entities.Booking{}

	if len(bookings) <= 1 {
		return nil
	}
	// 3. Усечь срез.
	bookings = bookings[:len(bookings)-1]
	return bookings
}

func BookingDbEquilToBooking(booking, bookingFromDB entities.Booking) bool {
	return bookingFromDB.Guest.Phone == booking.Guest.Phone &&
		bookingFromDB.Reservation.RoomNumber == booking.Reservation.RoomNumber &&
		bookingFromDB.Reservation.CheckIn.UTC() == booking.Reservation.CheckIn &&
		bookingFromDB.Reservation.CheckOut.UTC() == booking.Reservation.CheckOut &&
		bookingFromDB.Reservation.Price == booking.Reservation.Price
}

func (e *ServiceExcel) GetBookingInfoForPeriod(ctx context.Context, sheetName string, roomNumber string, start string, end string) ([]models.BookingInfo, error) {
	settings, err := e.storageSettings.Get(ctx, sheetName)
	if err != nil {
		return nil, err
	}

	sheetConfig := models.DbConvertToModel(*settings)

	sheet, err := excelCalendarScraper.NewSheet(sheetConfig)
	if err != nil {
		return nil, err
	}

	searchPeriod := excelCalendarScraper.NewSearchPeriod(start, end)

	bookings, err := sheet.GetBookingForPeriod(roomNumber, searchPeriod)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (e *ServiceExcel) GetAllBookingInfoForPeriod(ctx context.Context, sheetName string, start string, end string) ([]models.BookingInfo, error) {
	settings, err := e.storageSettings.Get(ctx, sheetName)
	if err != nil {
		return nil, err
	}

	sheetConfig := models.DbConvertToModel(*settings)

	sheet, err := excelCalendarScraper.NewSheet(sheetConfig)
	if err != nil {
		return nil, err
	}

	searchPeriod := excelCalendarScraper.NewSearchPeriod(start, end)

	allBookings := make([]models.BookingInfo, 0)
	for room := range sheet.GetApartMap() {
		bookings, err := sheet.GetBookingForPeriod(room, searchPeriod)
		if err != nil {
			return nil, err
		}
		allBookings = append(allBookings, bookings...)
	}

	return allBookings, nil
}

func (e *ServiceExcel) GetBookingForPeriod(ctx context.Context, sheetName string, roomNumber string, start string, end string) ([]entities.Booking, error) {
	settings, err := e.storageSettings.Get(ctx, sheetName)
	if err != nil {
		return nil, err
	}

	sheetConfig := models.DbConvertToModel(*settings)

	sheet, err := excelCalendarScraper.NewSheet(sheetConfig)
	if err != nil {
		return nil, err
	}

	searchPeriod := excelCalendarScraper.NewSearchPeriod(start, end)

	bs, err := sheet.GetBookingForPeriod(roomNumber, searchPeriod)
	if err != nil {
		return nil, err
	}

	bookings := make([]entities.Booking, 0, 5)
	for _, b := range bs {
		zap.L().Debug("GetBookingForPeriod", zap.Any("b", b))
		guest, err := b.DbModelConvertGuest()
		if err != nil {
			return nil, err
		}
		reservation, err := b.DbModelConvert(guest.GuestID)
		if err != nil {
			return nil, fmt.Errorf("%w\nmodels.BookingInfo: %v", err, b)
		}
		booking := entities.Booking{
			Guest:       *guest,
			Reservation: *reservation,
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}
