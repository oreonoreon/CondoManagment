package config

import (
	"awesomeProject/internal/entities"
	"awesomeProject/internal/excelCalendarScraper/models"
	"crypto/rand"
	"encoding/hex"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
)

// minSessionSecretLen — минимальная длина ключа подписи сессий в production.
const minSessionSecretLen = 32

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
	// SessionSecret — ключ подписи cookie сессий (env SESSION_SECRET).
	SessionSecret string
	// SessionSecretGenerated — SESSION_SECRET не задан и ключ сгенерирован случайно (только вне production).
	SessionSecretGenerated bool
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

	conf.SessionSecret, conf.SessionSecretGenerated = loadSessionSecret(conf.IsProduction)

	return conf
}

// loadSessionSecret читает SESSION_SECRET. В production ключ обязателен;
// в разработке без него генерируется случайный (сессии сбросятся после перезапуска).
func loadSessionSecret(isProduction bool) (secret string, generated bool) {
	secret = os.Getenv("SESSION_SECRET")
	if isProduction {
		if len(secret) < minSessionSecretLen {
			panic("SESSION_SECRET must be set and be at least 32 characters long in production")
		}
		return secret, false
	}
	if secret != "" {
		return secret, false
	}

	b := make([]byte, minSessionSecretLen)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b), true
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
