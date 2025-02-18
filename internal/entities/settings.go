package entities

type Settings struct {
	ExcelPath               string `db:"excel_path"`
	SheetName               string `db:"sheet_name"`
	SheetStartColumNum      int    `db:"sheet_startColumNum"`
	SheetYearAndMonthRowNum int    `db:"sheet_yearAndMonthRowNum"`
	SheetDaysRowNum         int    `db:"sheet_daysRowNum"`
	SheetStartRow           int    `db:"sheet_startRow"`
	SheetApartmentColumNum  int    `db:"sheet_apartmentColumNum"`
}
