package report

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

// TotalPriceReportSheetName название листа с отчётом суммарных цен по апартаментам
const TotalPriceReportSheetName = "TotalPriceReport"

// TotalPriceForPeriodReport строит xlsx файл, где по вертикали (столбец A, начиная со строки 2)
// расположены названия апартаментов, а по горизонтали (строка 1, начиная со столбца B) - месяцы
// указанного периода.
//
// pricesByApartmentAndMonth: имя апартамента -> ключ месяца ("01.2006") -> суммарная цена за месяц
// apartmentNames: порядок строк (названия апартаментов)
// months: порядок и набор столбцов (первое число каждого месяца периода)
func TotalPriceForPeriodReport(pricesByApartmentAndMonth map[string]map[string]int, apartmentNames []string, months []time.Time) ([]byte, error) {
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			zap.L().Error("TotalPriceForPeriodReport/file.Close()", zap.Error(err))
		}
	}()

	sheetName := TotalPriceReportSheetName
	sheetIndex, err := file.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать лист %s: %w", sheetName, err)
	}
	file.SetActiveSheet(sheetIndex)
	if err := file.DeleteSheet("Sheet1"); err != nil {
		return nil, fmt.Errorf("не удалось удалить лист по умолчанию: %w", err)
	}

	headerStyle, err := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return nil, fmt.Errorf("не удалось создать стиль заголовка: %w", err)
	}

	if err := SetCellStyleAndValue(file, sheetName, "A1", headerStyle, "Апартамент"); err != nil {
		return nil, err
	}

	monthColumn := make(map[string]string, len(months))
	for colIdx, month := range months {
		colName, err := excelize.ColumnNumberToName(colIdx + 2) // столбец B и далее
		if err != nil {
			return nil, fmt.Errorf("не удалось получить имя столбца для месяца %s: %w", month.Format("01.2006"), err)
		}
		monthKey := month.Format("01.2006")
		monthColumn[monthKey] = colName

		cell := colName + "1"
		if err := SetCellStyleAndValue(file, sheetName, cell, headerStyle, monthKey); err != nil {
			return nil, err
		}
	}

	for rowIdx, apartmentName := range apartmentNames {
		row := rowIdx + 2

		cellA := fmt.Sprintf("A%d", row)
		if err := file.SetCellValue(sheetName, cellA, apartmentName); err != nil {
			return nil, fmt.Errorf("не удалось установить значение в %s: %w", cellA, err)
		}

		monthPrices := pricesByApartmentAndMonth[apartmentName]
		for _, month := range months {
			monthKey := month.Format("01.2006")
			cell := fmt.Sprintf("%s%d", monthColumn[monthKey], row)

			price := 0
			if monthPrices != nil {
				price = monthPrices[monthKey]
			}

			if err := file.SetCellValue(sheetName, cell, price); err != nil {
				return nil, fmt.Errorf("не удалось установить значение в %s: %w", cell, err)
			}
		}
	}

	if err := file.SetColWidth(sheetName, "A", "A", 20); err != nil {
		return nil, fmt.Errorf("не удалось установить ширину столбца: %w", err)
	}

	buffer, err := file.WriteToBuffer()
	if err != nil {
		zap.L().Error("TotalPriceForPeriodReport/file.WriteToBuffer()", zap.Error(err))
		return nil, err
	}

	return buffer.Bytes(), nil
}
