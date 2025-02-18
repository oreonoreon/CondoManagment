package eventProcessor

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/erro"
	"awesomeProject/internal/excelCalendarScraper"
	"awesomeProject/internal/services"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func P(ctx context.Context, service *services.Service, e excelCalendarScraper.ExcelFilePath, config config.Config, searchPeriod excelCalendarScraper.SearchPeriod, roomNumber string) error {
	file, err := excelCalendarScraper.ExcelFileOpen(e.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	sheetMap, err := excelCalendarScraper.NewSheetMap(config)
	if err != nil {
		return err
	}

	sheet, err := excelCalendarScraper.NewSheet(config)
	if err != nil {
		return err
	}

	if roomNumber == "All" {
		for i := range sheetMap.GetApartMap() {
			bookings, err := sheet.GetBookingForPeriod(i, searchPeriod)
			if err != nil {
				fmt.Println(err)
			}

			//for test
			//for _, booking := range bookings {
			//	fmt.Println(booking)
			//	//zap.L().Debug("writeinDb", zap.Any("booking", booking))
			//}

			for _, booking := range bookings {
				err = service.CreateReservation(ctx, booking)
				if err != nil {
					zap.L().Error("writeinDb", zap.Error(err))
				}
				//уведомления пользователя об ошибке создания записи данных
				if !errors.As(erro.ErrFullyMatchOtherBooking, &err) {
					fmt.Println(booking, err)
				}

			}
		}

	} else {
		bookings, err := sheet.GetBookingForPeriod(roomNumber, searchPeriod)
		if err != nil {
			return err
		}

		//for test
		//for _, booking := range bookings {
		//	fmt.Println(booking)
		//	//zap.L().Debug("writeinDb", zap.Any("booking", booking))
		//}

		for _, booking := range bookings {
			err = service.CreateReservation(ctx, booking)
			if err != nil {
				zap.L().Error("writeinDb", zap.Error(err))
			}
			//уведомления пользователя об ошибке создания записи данных
			fmt.Println(err)
		}
	}

	return nil

}
