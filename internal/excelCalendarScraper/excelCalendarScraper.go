package excelCalendarScraper

import (
	"awesomeProject/internal/myLogger"
	"awesomeProject/internal/repo"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
	"time"
)

func ExcelCalendarScraper(model repo.DBSql) {
	file, err := excelize.OpenFile("CalendarCheckInCheckOut.xlsx")
	if err != nil {
		myLogger.Logger.Println(err)
		return
	}
	defer file.Close()

	//getByRows(file)
	datesMap, err := writeInMapDatesAndColumesNames(file, 7)
	if err != nil {
		myLogger.Logger.Println(err)
		return
	}

	reservations, guests, err := getAllBookingForApartment("A602", file, "2023-11-01", 9, datesMap)
	if err != nil {
		myLogger.Logger.Println(err)
		return
	}
	myLogger.Logger.Println(reservations)
	myLogger.Logger.Println(guests)

	//getCellValueByName(file, "HX9")
	//getMergeCells(file)
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

func getAllBookingForApartment(roomNumber string, file *excelize.File, startDate string, startRow int, dateMap map[string]string) ([]repo.Reservaton, []repo.Guest, error) {
	date := startDate
	var value string
	reservation := repo.Reservaton{RoomNumber: roomNumber}
	reservations := make([]repo.Reservaton, 0, 15)
	var guest repo.Guest
	guests := make([]repo.Guest, 0, 15)

	for {
		colum, ok := dateMap[date]
		if !ok {
			err := errors.New("date not found in Map of date")
			myLogger.Logger.Println(fmt.Errorf("%v: %w", date, err))
			break
		}
		v, err := getCellValueByName(file, colum+strconv.Itoa(startRow))
		if err != nil {
			return nil, nil, err
		}

		switch {
		case value == "":
			//todo исправить что бы повторно не конвертить в time.Time
			timeDate, err := time.Parse(time.DateOnly, date)
			if err != nil {
				return nil, nil, err
			}
			value = v
			reservation.CheckIn = timeDate
			reservation.CheckOut = timeDate
		case value != "" && v == value:
			//todo исправить что бы повторно не конвертить в time.Time
			timeDate, err := time.Parse(time.DateOnly, date)
			if err != nil {
				return nil, nil, err
			}
			value = v
			reservation.CheckOut = timeDate
		case value != v:
			timeDate, err := time.Parse(time.DateOnly, date)
			if err != nil {
				return nil, nil, err
			}
			guest = repo.Guest{GuestID: uuid.New(), Name: parseName(value), Phone: parsePhone(value), Description: parseDescription(value)}
			reservation.GuestID = guest.GuestID

			guests = append(guests, guest)
			guest = repo.Guest{}

			reservations = append(reservations, reservation)
			reservation = repo.Reservaton{RoomNumber: roomNumber}
			reservation.CheckIn = timeDate
			reservation.CheckOut = timeDate
			value = v
		}

		// add one day to date
		date, err = datePlusOneDay(date)
		if err != nil {
			return nil, nil, err
		}
	}

	return reservations, guests, nil
}
func parsePhone(cellValue string) string {
	return ""
}

func parseName(cellValue string) string {
	before, _, boo := strings.Cut(cellValue, " ")
	if boo {
		return before
	}
	return ""
}

func parseDescription(cellValue string) string {
	return cellValue
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
