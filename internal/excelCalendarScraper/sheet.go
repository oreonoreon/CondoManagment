package excelCalendarScraper

import (
	"awesomeProject/internal/excelCalendarScraper/models"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"strconv"
)

type Sheet struct {
	SheetConfig
	SheetMap
}

func NewSheet(config models.SheetConfig) (Sheet, error) {
	con := NewSheetConfig(config.Path, config.Name, config.StartColumNum, config.YearAndMonthRowNum, config.DaysRowNum, config.StartRow, config.ApartmentColumNum)

	sheetMap, err := NewSheetMap(con)
	if err != nil {
		return Sheet{}, err
	}
	return Sheet{con, *sheetMap}, nil
}

type SheetMap struct {
	dateMap  map[string]string
	apartMap map[string]int
}

func (s SheetMap) GetDateMap() map[string]string {
	return s.dateMap
}

func (s SheetMap) GetApartMap() map[string]int {
	return s.apartMap
}

func NewSheetMap(config SheetConfig) (*SheetMap, error) {
	file, err := ExcelFileOpen(config.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dateMap, err := config.writeInMapDatesAndColumesNames(file)
	if err != nil {
		return nil, err
	}
	apartMap, err := config.apartmentsMap(file)
	if err != nil {
		return nil, err
	}
	return &SheetMap{dateMap: dateMap, apartMap: apartMap}, nil
}

func ExcelFileOpen(path string) (*excelize.File, error) {
	file, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (e Sheet) GetBookingForPeriod(roomNumber string, searchPeriod SearchPeriod) ([]models.BookingInfo, error) {
	file, err := ExcelFileOpen(e.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var value string
	bookingInfoSlice := make([]models.BookingInfo, 0, 15)

	for date := searchPeriod.startDate; date <= searchPeriod.endDate; {
		colum, ok := e.dateMap[date]
		if !ok {
			err := errors.New("date was not found in Map of date")
			zap.L().Error("getBookingForPeriod", zap.String("date", date), zap.Error(err))
			zap.L().Debug("getBookingForPeriod", zap.Any("dateMap", e.dateMap))
			break
		}

		v, err := getCellValueByName(file, colum+strconv.Itoa(e.apartMap[roomNumber]), e.Name)
		if err != nil {
			return nil, err
		}

		switch {
		case value == "":
			value = v

		case value != "" && v == value:
			if date == searchPeriod.endDate {
				bookingInfo, err := parseValue(value, roomNumber)
				if err != nil {
					zap.L().Error("getAllBookingForApartment", zap.Error(err), zap.Any("bookingInfo", bookingInfo))
					return nil, fmt.Errorf("%w:\n%v", err, bookingInfo)
				}
				bookingInfoSlice = append(bookingInfoSlice, bookingInfo)
			}
			value = v
		case value != v:
			bookingInfo, err := parseValue(value, roomNumber)
			if err != nil {
				zap.L().Error("getAllBookingForApartment", zap.Error(err), zap.Any("bookingInfo", bookingInfo))
				return nil, fmt.Errorf("%w:\n%v", err, bookingInfo)
			}
			bookingInfoSlice = append(bookingInfoSlice, bookingInfo)
			value = v
		}

		// add one day to date
		date, err = datePlusOneDay(date)
		if err != nil {
			return nil, err
		}
	}

	return bookingInfoSlice, nil
}
