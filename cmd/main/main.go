package main

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/repo"
	"awesomeProject/internal/server"
	"awesomeProject/internal/services"
	"context"
	"github.com/gopsql/standard"
	"go.uber.org/zap"
)

func main() {
	conf := config.InitConfig()

	zapConfig := zap.NewProductionConfig()
	if conf.DebugLogLevel {
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	zapConfig.OutputPaths = []string{conf.LogPath}
	logger, err := zapConfig.Build()

	defer zap.L().Sync()

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	zap.L().Info("APPLICATION START")

	//create Db connection
	db, err := repo.ConnectionPostgreSQl()
	if err != nil {
		panic(err)
	}

	defer func(db *standard.DB) {
		err := db.Close()
		if err != nil {
			zap.L().Error("db.Close()", zap.Error(err))
		}
	}(db)

	//db
	postgre := repo.NewRepository(db)

	//services
	serviceReservation := services.NewService(postgre, postgre)
	serviceSettings := services.NewServiceSettings(postgre)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//writing in db the settings
	_, err = serviceSettings.Set(ctx, &conf)
	if err != nil {
		panic(err)
	}

	//services
	serviceExcel := services.NewServiceExcel(postgre)

	//handlers
	handler := server.NewHandle(serviceReservation, serviceSettings, serviceExcel)

	//server
	ser := server.NewServer(handler)
	ser.StartServer()

}

//func excel(ctx context.Context, s *services.Service, conf config.Config) {
//	//config := excelCalendarScraper.NewSheetConfig("ноя2024-ноя2025", 5, "1", "2", 4, 3)
//	searchPeriod := excelCalendarScraper.NewSearchPeriod("2024-11-01", "2025-11-01")
//
//	//config := excelCalendarScraper.NewSheetConfig("ноя2023-ноя2024", 7, "1", "2", 4, 3)
//	//searchPeriod := excelCalendarScraper.NewSearchPeriod("2023-11-01", "2024-11-01")
//
//	err := eventProcessor.P(ctx, s, conf.ExcelFilePath, conf, searchPeriod, "F107")
//	if err != nil {
//		zap.L().Error("P", zap.Error(err))
//	}
//}
