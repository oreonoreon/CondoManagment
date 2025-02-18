package excelCalendarScraper

import (
	"errors"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

func (s SheetConfig) writeInMapDatesAndColumesNames(file *excelize.File) (map[string]string, error) {
	datesMap := make(map[string]string)
	for {
		str, err := excelize.ColumnNumberToName(s.StartColumNum)
		if err != nil {
			return nil, err
		}

		YearAndMonthRowNum := strconv.Itoa(s.YearAndMonthRowNum)

		yearAndMonth, err := getCellValueByName(file, str+YearAndMonthRowNum, s.Name)
		if err != nil {
			return nil, err
		}
		if yearAndMonth == "" {
			zap.L().Debug("writeInMapDatesAndColumesNames", zap.Error(errors.New("end of months in sheet")))
			break
		}

		DaysRowNum := strconv.Itoa(s.DaysRowNum)

		value2, err := getCellValueByName(file, str+DaysRowNum, s.Name)
		if err != nil {
			return nil, err
		}
		if value2 == "" {
			zap.L().Debug("writeInMapDatesAndColumesNames", zap.Error(errors.New("end of date in sheet")))
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

		s.StartColumNum++
	}
	return datesMap, nil
}

func monthConverter(month string) (time.Month, error) {
	switch month {
	case "January", "january", "январь":
		return time.January, nil
	case "February", "february", "февраль":
		return time.February, nil
	case "March", "march", "март":
		return time.March, nil
	case "April", "april", "апрель":
		return time.April, nil
	case "May", "may", "май":
		return time.May, nil
	case "June", "june", "июнь":
		return time.June, nil
	case "July", "july", "июль":
		return time.July, nil
	case "August", "august", "август":
		return time.August, nil
	case "September", "september", "сентябрь":
		return time.September, nil
	case "October", "october", "октябрь":
		return time.October, nil
	case "November", "november", "ноябрь":
		return time.November, nil
	case "December", "december", "декабрь":
		return time.December, nil
	}
	return 0, errors.New("error parsing name of month")
}
