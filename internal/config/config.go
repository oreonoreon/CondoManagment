package config

import (
	"awesomeProject/internal/entities"
	"awesomeProject/internal/excelCalendarScraper/models"
	"go.uber.org/zap"
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

type ConfigEnv struct {
	FrontURL     string
	IsProduction bool
	DatabaseDSN  string
	//SessionSecret string
	//ServerPort    string
}

func LoadEnv() *ConfigEnv {
	conf := new(ConfigEnv)

	ginMode, ok := os.LookupEnv("GIN_MODE")
	if !ok {
		conf.IsProduction = false
		zap.L().Info("GIN_MODE not set, default to debug")
	} else {
		if ginMode == "release" {
			conf.IsProduction = true
		} else {
			conf.IsProduction = false
		}
	}

	allowOrigin, ok := os.LookupEnv("FRONT_URL")
	if !ok {
		allowOrigin = "http://localhost:5173" // явно указываем для разработки
	}
	conf.FrontURL = allowOrigin

	conf.DatabaseDSN = buildDataSourceName()

	return conf
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func buildDataSourceName() string {
	defaultDatabaseDSN := "postgres://oreonoreon:12345@localhost:5432/postgres?sslmode=disable" // "postgres://oreonoreon:12345@postgres:5432/postgres?sslmode=disable"

	databaseDSN, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		DB_HOST := os.Getenv("DB_HOST")
		DB_USER := os.Getenv("DB_USER")
		DB_PASS := os.Getenv("DB_PASS")
		DB_NAME := os.Getenv("DB_NAME")

		if DB_HOST != "" && DB_USER != "" && DB_PASS != "" && DB_NAME != "" {
			return "postgres://" + DB_USER + ":" + DB_PASS + "@" + DB_HOST + ":5432/" + DB_NAME + "?sslmode=disable"
		}

		return defaultDatabaseDSN
	} else {
		return databaseDSN
	}
}
