package config

import (
	"awesomeProject/internal/entities"
	"awesomeProject/internal/excelCalendarScraper/models"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	LogPath            string `yaml:"logPath"`
	DebugLogLevel      bool   `yaml:"debugLogLevel"`
	models.SheetConfig `yaml:"sheetConfig"`
}

func InitConfig() Config {
	var config Config
	var sheetConfig models.SheetConfig

	config.SheetConfig = sheetConfig

	conf, err := os.ReadFile("./etc/config.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(conf, &config)
	if err != nil {
		panic(err)
	}
	return config
}

func (c *Config) SettingsConvert() entities.Settings {
	return entities.Settings{
		ExcelPath:               c.Path,
		SheetName:               c.Name,
		SheetStartColumNum:      c.StartColumNum,
		SheetYearAndMonthRowNum: c.YearAndMonthRowNum,
		SheetDaysRowNum:         c.DaysRowNum,
		SheetStartRow:           c.StartRow,
		SheetApartmentColumNum:  c.ApartmentColumNum,
	}
}
