package models

import "awesomeProject/internal/entities"

type SheetConfig struct {
	Path               string `yaml:"path"`
	Name               string `yaml:"name"`
	StartColumNum      int    `yaml:"startColumNum"`
	YearAndMonthRowNum int    `yaml:"yearAndMonthRowNum"`
	DaysRowNum         int    `yaml:"daysRowNum"`
	StartRow           int    `yaml:"startRow"`
	ApartmentColumNum  int    `yaml:"apartmentColumNum"`
}

func DbConvertToModel(settings entities.Settings) SheetConfig {
	return SheetConfig{
		Path:               settings.ExcelPath,
		Name:               settings.SheetName,
		StartColumNum:      settings.SheetStartColumNum,
		YearAndMonthRowNum: settings.SheetYearAndMonthRowNum,
		DaysRowNum:         settings.SheetDaysRowNum,
		StartRow:           settings.SheetStartRow,
		ApartmentColumNum:  settings.SheetApartmentColumNum,
	}
}
