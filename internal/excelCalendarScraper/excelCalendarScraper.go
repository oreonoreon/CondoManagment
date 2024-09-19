package excelCalendarScraper

import (
	"awesomeProject/internal/myLogger"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func ExcelCalendarScraper() {
	file, err := excelize.OpenFile("CalendarCheckInCheckOut.xlsx")
	if err != nil {
		myLogger.Logger.Println(err)
		return
	}
	defer func(file *excelize.File) {
		err := file.Close()
		if err != nil {
			myLogger.Logger.Println(err)
		}
	}(file)

	datesMap, err := writeInMapDatesAndColumesNames(file, 7)
	if err != nil {
		myLogger.Logger.Println(err)
		return
	}

	excelSettings := excelParserSettings{file, "2023-11-01", 9, datesMap}

	bookings, err := excelSettings.getAllBookingForApartment("A602")
	if err != nil {
		myLogger.Logger.Println(err)
		return
	}

	for _, booking := range bookings {
		myLogger.Logger.Println(booking)
	}
}

func getByRows(file *excelize.File) {
	rows, err := file.GetRows("ноя2023-ноя2024")
	if err != nil {
		myLogger.Logger.Fatal(err)
	}

	for i, row := range rows {
		if i == 9 {
			for _, col := range row {
				myLogger.Logger.Print(col, "\t")
			}
		}
	}
}
func datePlusOneDay(date string) (string, error) {
	timeDate, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return "", err
	}
	return timeDate.Add(time.Hour * 24).Format(time.DateOnly), nil
}

type excelParserSettings struct {
	file      *excelize.File
	startDate string
	startRow  int
	dateMap   map[string]string
}

func (e excelParserSettings) getAllBookingForApartment(roomNumber string) ([]BookingInfo, error) {
	date := e.startDate
	var value string
	bookingInfoSlice := make([]BookingInfo, 0, 15)

	for {
		colum, ok := e.dateMap[date]
		if !ok {
			err := errors.New("date not found in Map of date")
			myLogger.Logger.Println(fmt.Errorf("%v: %w", date, err))
			break
		}
		v, err := getCellValueByName(e.file, colum+strconv.Itoa(e.startRow))
		if err != nil {
			return nil, err
		}

		switch {
		case value == "":
			value = v

		case value != "" && v == value:
			value = v
		case value != v:
			bookingInfo, err := parseValue(value, roomNumber)
			if err != nil {
				myLogger.Logger.Println(err)
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

func parseValue(cellValue string, roomNumber string) (BookingInfo, error) {
	s := strings.FieldsFunc(cellValue, func(c rune) bool {
		return c == ':' || c == ','
	})

	bookingInfo := BookingInfo{RoomNumber: roomNumber}
	bookingInfo.parseOutBookingInfo(s)
	return bookingInfo, nil
}
func (b *BookingInfo) parseOutBookingInfo(s []string) {
	for k, v := range s {
		v = strings.TrimFunc(v, func(c rune) bool {
			return unicode.IsSpace(c)
		})

		//проверим на наличие следующего индекса в слайсе
		if len(s) < k+1+1 {
			return
		}

		switch v {
		case "Name":
			b.GuestName = s[k+1]
		case "Check in":
			b.CheckIn = s[k+1]
		case "Check out":
			b.CheckOut = s[k+1]
		case "Price":
			b.Price = s[k+1]
		case "Cleaning price":
			b.CleaningPrice = s[k+1]
		case "Electricity and water payment":
			b.ElectricityAndWaterPayment = s[k+1]
		case "Adult":
			b.Adult = s[k+1]
		case "children":
			b.Children = s[k+1]
		case "Phone":
			b.Phone = s[k+1]
		case "Description":
			b.Description = s[k+1]
		}
	}
}

type BookingInfo struct {
	RoomNumber                 string
	CheckIn                    string
	CheckOut                   string
	GuestName                  string
	Phone                      string
	Price                      string
	CleaningPrice              string
	ElectricityAndWaterPayment string
	Adult                      string
	Children                   string
	Description                string
}

func writeInMapDatesAndColumesNames(file *excelize.File, i int) (map[string]string, error) {
	datesMap := make(map[string]string)
	for {
		str, err := excelize.ColumnNumberToName(i)
		if err != nil {
			return nil, err
		}
		yearAndMonth, err := getCellValueByName(file, str+"1")
		if err != nil {
			return nil, err
		}
		if yearAndMonth == "" {
			myLogger.Logger.Println(errors.New("end of months in sheet"))
			break
		}
		value2, err := getCellValueByName(file, str+"2")
		if err != nil {
			return nil, err
		}
		if value2 == "" {
			myLogger.Logger.Println(errors.New("end of date in sheet"))
			break
		}

		yearStr, monthStr, found := strings.Cut(yearAndMonth, ", ")
		if !found {
			return nil, errors.New("can not split the string which contain year and month because there is no ',' in the string")
		}

		year, err := strconv.Atoi(yearStr)
		if err != nil {
			return nil, err
		}

		day, err := strconv.Atoi(value2)
		if err != nil {
			return nil, err
		}
		month, err := monthConverter(monthStr)
		if err != nil {
			return nil, err
		}

		date := time.Date(year, month, day, 12, 0, 0, 0, time.UTC)

		datesMap[date.Format(time.DateOnly)] = str

		i++
	}
	return datesMap, nil
}

func monthConverter(month string) (time.Month, error) {
	switch month {
	case "January", "январь":
		return time.January, nil
	case "February", "февраль":
		return time.February, nil
	case "March", "март":
		return time.March, nil
	case "April", "апрель":
		return time.April, nil
	case "May", "май":
		return time.May, nil
	case "June", "июнь":
		return time.June, nil
	case "July", "июль":
		return time.July, nil
	case "August", "август":
		return time.August, nil
	case "September", "сентябрь":
		return time.September, nil
	case "October", "октябрь":
		return time.October, nil
	case "November", "ноябрь":
		return time.November, nil
	case "December", "декабрь":
		return time.December, nil
	}
	return 0, errors.New("error parsing name of month")
}

func getCellValueByName(file *excelize.File, cellName string) (string, error) {
	str, err := file.GetCellValue("ноя2023-ноя2024", cellName)
	if err != nil {
		return "", err
	}
	return str, nil
}

func getMergeCells(file *excelize.File) {
	mergecells, err := file.GetMergeCells("ноя2023-ноя2024")
	if err != nil {
		myLogger.Logger.Fatal(err)
	}
	for _, mergecell := range mergecells {
		myLogger.Logger.Println(mergecell.GetCellValue())
		myLogger.Logger.Println(mergecell.GetStartAxis())
		myLogger.Logger.Println(mergecell.GetEndAxis())
	}
}
