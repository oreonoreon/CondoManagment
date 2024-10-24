package main

import (
	"awesomeProject/internal/eventProcessor"
	"awesomeProject/internal/excelCalendarScraper"
	"awesomeProject/internal/myLogger"
	"awesomeProject/internal/repo"
	"awesomeProject/internal/services"
	"context"
	"github.com/gopsql/standard"
	"gopkg.in/yaml.v3"
	"os"
)

func main() {
	var excelFilePath excelCalendarScraper.ExcelFilePath
	conf, err := os.ReadFile("./etc/config.yaml")
	if err != nil {
		myLogger.Logger.Println(err)
	}

	err = yaml.Unmarshal(conf, &excelFilePath)
	if err != nil {
		myLogger.Logger.Println(err)
	}

	myLogger.Logger.Println(excelFilePath)

	//create Db connection
	db, err := repo.ConnectionPostgreSQl()
	if err != nil {
		panic(err)
	}

	defer func(db *standard.DB) {
		err := db.Close()
		if err != nil {
			myLogger.Logger.Println(err)
		}
	}(db)

	postgre := repo.NewRepository(db)

	s := services.NewService(postgre)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = eventProcessor.ProcesFullfilingDBWithDataFromExcel(ctx, s, excelFilePath)
	if err != nil {
		myLogger.Logger.Println(err)
	}

	//server.Server(postgre)
}
