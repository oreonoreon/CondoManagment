package excelCalendarScraper

import (
	"errors"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"regexp"
	"strconv"
)

func (s SheetConfig) apartmentsMap(file *excelize.File) (map[string]int, error) {
	apartmentsMap := make(map[string]int)
	for {
		str, err := excelize.ColumnNumberToName(s.ApartmentColumNum)
		if err != nil {
			return nil, err
		}

		value, err := getCellValueByName(file, str+strconv.Itoa(s.StartRow), s.Name)
		if err != nil {
			return nil, err
		}

		if value != "" {
			values, err := validate(value)
			if err != nil {
				return nil, err
			}
			apartmentsMap[values[0]] = s.StartRow
			s.StartRow++
		} else {
			break
		}
	}
	zap.L().Debug("apartmentsMap", zap.Any("apartmentsMap", apartmentsMap))
	return apartmentsMap, nil
}

func validate(value string) ([]string, error) {
	r := regexp.MustCompile(`[A-Z]{1,2}[0-9]{3}-?[0-9]{0,3}`)
	sliceStrings := r.FindAllString(value, -1)
	if len(sliceStrings) == 0 {
		zap.L().Debug("validate value", zap.Any("value", value))
		return nil, errors.New("validate failed")
	}
	return sliceStrings, nil
}
