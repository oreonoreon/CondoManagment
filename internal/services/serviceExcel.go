package services

import (
	"awesomeProject/internal/excelCalendarScraper"
	"awesomeProject/internal/models"
	"context"
)

type ServiceExcel struct {
	storageSettings StorageSettings
}

func NewServiceExcel(storageSettings StorageSettings) *ServiceExcel {
	return &ServiceExcel{storageSettings: storageSettings}
}

func (e *ServiceExcel) GetBookingForPeriod(ctx context.Context, sheetName string, roomNumber string, start string, end string) ([]models.BookingInfo, error) {
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

func (e *ServiceExcel) GetAllBookingForPeriod(ctx context.Context, sheetName string, start string, end string) ([]models.BookingInfo, error) {
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
