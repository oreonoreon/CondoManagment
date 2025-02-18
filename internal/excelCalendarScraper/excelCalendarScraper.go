package excelCalendarScraper

import (
	"awesomeProject/internal/models"
	"github.com/xuri/excelize/v2"
	"strings"
	"time"
	"unicode"
)

type SheetConfig struct {
	Path               string `yaml:"path"`
	Name               string `yaml:"name"`
	StartColumNum      int    `yaml:"startColumNum"`
	YearAndMonthRowNum int    `yaml:"yearAndMonthRowNum"`
	DaysRowNum         int    `yaml:"daysRowNum"`
	StartRow           int    `yaml:"startRow"`
	ApartmentColumNum  int    `yaml:"apartmentColumNum"`
}

func NewSheetConfig(path string, name string, startColumNum int, yearAndMonthRowNum int, daysRowNum int, startRow int, apartmentColumNum int) SheetConfig {
	return SheetConfig{
		Path:               path,
		Name:               name,
		StartColumNum:      startColumNum,
		YearAndMonthRowNum: yearAndMonthRowNum,
		DaysRowNum:         daysRowNum,
		StartRow:           startRow,
		ApartmentColumNum:  apartmentColumNum,
	}
}

func NewSearchPeriod(startDate string, endDate string) SearchPeriod {
	return SearchPeriod{startDate: startDate, endDate: endDate}
}

type SearchPeriod struct {
	startDate string
	endDate   string
}

func datePlusOneDay(date string) (string, error) {
	timeDate, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return "", err
	}
	return timeDate.Add(time.Hour * 24).Format(time.DateOnly), nil
}

func parseValue(cellValue string, roomNumber string) (models.BookingInfo, error) {
	s := strings.FieldsFunc(cellValue, func(c rune) bool {
		return c == ':' || c == ','
	})
	//удалим пробелы
	for i, v := range s {
		s[i] = strings.TrimFunc(v, func(c rune) bool {
			return unicode.IsSpace(c)
		})
	}

	bookingInfo := new(models.BookingInfo)
	bookingInfo.RoomNumber = roomNumber
	bookingInfo.ParseOutBookingInfo(s)
	return *bookingInfo, nil
}

func getCellValueByName(file *excelize.File, cellName string, sheetname string) (string, error) {
	str, err := file.GetCellValue(sheetname, cellName)
	if err != nil {
		return "", err
	}
	return str, nil
}
