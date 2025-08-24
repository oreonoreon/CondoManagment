package main

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/repo"
	"awesomeProject/internal/server"
	"awesomeProject/internal/services"
	"context"
	"github.com/gopsql/standard"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func main() {

	conf := config.InitConfig()

	zapConfig := zap.NewProductionConfig()
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
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
	serviceExcel := services.NewServiceExcel(postgre, serviceReservation)
	serviceApartment := services.NewServiceApartment(postgre)
	serviceBnB := services.NewServiceBnB(postgre, postgre, postgre, postgre, postgre)
	servicesUsers := services.NewServiceUsers(postgre)

	//handlers
	handler := server.NewHandle(serviceReservation, serviceSettings, serviceExcel, serviceApartment, serviceBnB, servicesUsers)

	//server
	server.Gin(handler)
}
