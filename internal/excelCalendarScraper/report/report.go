package report

import (
	"awesomeProject/internal/entities"
	"fmt"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"strconv"
)

func ReportForOwner(info []entities.Booking) (string, error) {
	file, err := excelize.OpenFile("./etc/Book_1.xlsx")
	if err != nil {
		return "", err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			zap.L().Error("ReportForOwner/file.Close()", zap.Error(err))
		}
	}()

	var sheetName = "Report"
	targetCells := []string{"A2", "B2", "C2", "D2", "E2", "F2", "G2", "H2", "I2", "J2", "K2", "L2"}

	cellsStyleID := make([]int, 0, 12)
	for _, sourceCell := range targetCells {
		// Получаем стиль
		styleID, err := file.GetCellStyle(sheetName, sourceCell)
		if err != nil {
			return "", fmt.Errorf("не удалось получить стиль из %s: %w", sourceCell, err)
		}
		cellsStyleID = append(cellsStyleID, styleID)
	}

	// Получаем высоту строки 2
	height, err := file.GetRowHeight(sheetName, 2)
	if err != nil {
		return "", fmt.Errorf("не удалось получить высоту строки: %w", err)
	}

	for i, bookingInfo := range info {
		row := strconv.Itoa(i + 2)

		err := file.SetRowHeight(sheetName, i+2, height)
		if err != nil {
			return "", fmt.Errorf("не удалось установить высоту строки %s: %w", row, err)
		}

		//A
		err = file.SetCellStyle(sheetName, "A"+row, "A"+row, cellsStyleID[0])
		if err != nil {
			return "", fmt.Errorf("не удалось применить стиль к %s: %w", "A"+row, err)
		}
		err = file.SetCellValue(sheetName, "A"+row, bookingInfo.RoomNumber)
		if err != nil {
			return "", err
		}

		//B
		err = file.SetCellStyle(sheetName, "B"+row, "B"+row, cellsStyleID[1])
		if err != nil {
			return "", fmt.Errorf("не удалось применить стиль к %s: %w", "A"+row, err)
		}
		err = file.SetCellValue(sheetName, "B"+row, bookingInfo.CheckIn)
		if err != nil {
			return "", err
		}

		//C
		err = file.SetCellStyle(sheetName, "C"+row, "C"+row, cellsStyleID[2])
		if err != nil {
			return "", fmt.Errorf("не удалось применить стиль к %s: %w", "A"+row, err)
		}
		err = file.SetCellValue(sheetName, "C"+row, bookingInfo.CheckOut)
		if err != nil {
			return "", err
		}

		//D
		err = file.SetCellStyle(sheetName, "D"+row, "D"+row, cellsStyleID[3])
		if err != nil {
			return "", fmt.Errorf("не удалось применить стиль к %s: %w", "A"+row, err)
		}
		err = file.SetCellValue(sheetName, "D"+row, bookingInfo.Days)
		if err != nil {
			return "", err
		}

		//E
		err = file.SetCellStyle(sheetName, "E"+row, "E"+row, cellsStyleID[4])
		if err != nil {
			return "", fmt.Errorf("не удалось применить стиль к %s: %w", "A"+row, err)
		}
		err = file.SetCellValue(sheetName, "E"+row, bookingInfo.Name)
		if err != nil {
			return "", err
		}

		//F
		err = file.SetCellStyle(sheetName, "F"+row, "F"+row, cellsStyleID[5])
		if err != nil {
			return "", fmt.Errorf("не удалось применить стиль к %s: %w", "A"+row, err)
		}
		err = file.SetCellValue(sheetName, "F"+row, bookingInfo.Phone)
		if err != nil {
			return "", err
		}

		//G
		err = file.SetCellStyle(sheetName, "G"+row, "G"+row, cellsStyleID[6])
		if err != nil {
			return "", fmt.Errorf("не удалось применить стиль к %s: %w", "A"+row, err)
		}
		err = file.SetCellValue(sheetName, "G"+row, bookingInfo.Adult)
		if err != nil {
			return "", err
		}

		//H
		err = file.SetCellStyle(sheetName, "H"+row, "H"+row, cellsStyleID[7])
		if err != nil {
			return "", fmt.Errorf("не удалось применить стиль к %s: %w", "A"+row, err)
		}
		err = file.SetCellValue(sheetName, "H"+row, bookingInfo.Children)
		if err != nil {
			return "", err
		}

		//I
		err = file.SetCellStyle(sheetName, "I"+row, "I"+row, cellsStyleID[8])
		if err != nil {
			return "", fmt.Errorf("не удалось применить стиль к %s: %w", "A"+row, err)
		}
		err = file.SetCellValue(sheetName, "I"+row, bookingInfo.Price)
		if err != nil {
			return "", err
		}

		//J
		err = file.SetCellStyle(sheetName, "J"+row, "J"+row, cellsStyleID[9])
		if err != nil {
			return "", fmt.Errorf("не удалось применить стиль к %s: %w", "A"+row, err)
		}
		err = file.SetCellValue(sheetName, "J"+row, bookingInfo.Price*70/100)
		if err != nil {
			return "", err
		}

		//K
		err = file.SetCellStyle(sheetName, "K"+row, "K"+row, cellsStyleID[10])
		if err != nil {
			return "", fmt.Errorf("не удалось применить стиль к %s: %w", "A"+row, err)
		}
		err = file.SetCellValue(sheetName, "K"+row, bookingInfo.Price*30/100)
		if err != nil {
			return "", err
		}

		//L
		err = file.SetCellStyle(sheetName, "L"+row, "L"+row, cellsStyleID[11])
		if err != nil {
			return "", fmt.Errorf("не удалось применить стиль к %s: %w", "A"+row, err)
		}
		err = file.SetCellValue(sheetName, "L"+row, bookingInfo.PriceForOneNight)
		if err != nil {
			return "", err
		}

	}

	err = file.SaveAs("./etc/" + info[0].RoomNumber + ".xlsx")
	if err != nil {
		zap.L().Error("ReportForOwner/file.SaveAs()", zap.Error(err))
	}

	return file.Path, nil
}

func SetCellStyleAndValue(file *excelize.File, sheetName string, cell string, styleID int, value interface{}) error {
	err := file.SetCellStyle(sheetName, cell, cell, styleID)
	if err != nil {
		return fmt.Errorf("не удалось применить стиль к %s: %w", cell, err)
	}

	err = file.SetCellValue(sheetName, cell, value)
	if err != nil {
		return fmt.Errorf("не удалось установить значения в %s: %w", cell, err)
	}
	return nil
}

//func MontlyPriceReport(priceMap map[string]int) (string, error) {
//	file, err := excelize.OpenFile("./etc/Prices.xlsx")
//	if err != nil {
//		return "", err
//	}
//
//	defer func() {
//		err := file.Close()
//		if err != nil {
//			zap.L().Error("ReportForOwner/file.Close()", zap.Error(err))
//		}
//	}()
//
//	var sheetName = "Price"
//	var apartCellStart = "A2"
//	value, err := file.GetCellValue(sheetName, apartCellStart)
//	if err != nil {
//		return "", err
//	}
//	price := priceMap[value]
//
//}
