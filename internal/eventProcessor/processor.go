package eventProcessor

import (
	"awesomeProject/internal/excelCalendarScraper"
	"awesomeProject/internal/services"
	"golang.org/x/net/context"
	"log"
)

func ProcesFullfilingDBWithDataFromExcel(ctx context.Context, service *services.Service, e excelCalendarScraper.ExcelFilePath) error {
	bookings, err := e.ExcelCalendarScraper()
	if err != nil {
		return err
	}

	//for test
	for _, booking := range bookings {
		log.Println(booking)
	}

	for _, booking := range bookings {
		err = service.CreateReservation(ctx, booking)
	}

	return nil
}
