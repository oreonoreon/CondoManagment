package repo

import (
	"awesomeProject/internal/entities"
	"context"
	"database/sql"
	"errors"
)

func (db *Repository) Set(ctx context.Context, settings entities.Settings) (*entities.Settings, error) {
	set := new(entities.Settings)
	query := "INSERT INTO settings (" +
		"excel_path ," +
		"sheet_name ," +
		"sheet_startColumNum ," +
		"sheet_yearAndMonthRowNum ," +
		"sheet_daysRowNum ," +
		"sheet_startRow ," +
		"sheet_apartmentColumNum" +
		") " +
		"VALUES ($1,$2,$3,$4,$5,$6,$7) Returning *"

	row := db.PostgreSQL.QueryRowContext(
		ctx,
		query,
		settings.ExcelPath,
		settings.SheetName,
		settings.SheetStartColumNum,
		settings.SheetYearAndMonthRowNum,
		settings.SheetDaysRowNum,
		settings.SheetStartRow,
		settings.SheetApartmentColumNum,
	)

	err := row.Scan(
		&set.ExcelPath,
		&set.SheetName,
		&set.SheetStartColumNum,
		&set.SheetYearAndMonthRowNum,
		&set.SheetDaysRowNum,
		&set.SheetStartRow,
		&set.SheetApartmentColumNum,
	)
	if err != nil {
		return nil, err
	}
	return set, nil
}

func (db *Repository) Get(ctx context.Context, sheetName string) (*entities.Settings, error) {
	set := new(entities.Settings)
	query := "SELECT DISTINCT * from settings where sheet_name=$1"
	row := db.PostgreSQL.QueryRowContext(
		ctx,
		query,
		sheetName,
	)

	err := row.Scan(
		&set.ExcelPath,
		&set.SheetName,
		&set.SheetStartColumNum,
		&set.SheetYearAndMonthRowNum,
		&set.SheetDaysRowNum,
		&set.SheetStartRow,
		&set.SheetApartmentColumNum,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return set, nil
}

func (db *Repository) SettingsUpdate(ctx context.Context, settings entities.Settings) (*entities.Settings, error) {
	set := new(entities.Settings)
	query := "UPDATE settings SET " +
		"excel_path=$1 ," +
		"sheet_startColumNum=$2 ," +
		"sheet_yearAndMonthRowNum=$3 ," +
		"sheet_daysRowNum=$4 ," +
		"sheet_startRow=$5 ," +
		"sheet_apartmentColumNum=$6 " +
		"WHERE " +
		"sheet_name=$7 " +
		"Returning *"

	row := db.PostgreSQL.QueryRowContext(
		ctx,
		query,
		settings.ExcelPath,
		settings.SheetStartColumNum,
		settings.SheetYearAndMonthRowNum,
		settings.SheetDaysRowNum,
		settings.SheetStartRow,
		settings.SheetApartmentColumNum,
		settings.SheetName,
	)

	err := row.Scan(
		&set.ExcelPath,
		&set.SheetName,
		&set.SheetStartColumNum,
		&set.SheetYearAndMonthRowNum,
		&set.SheetDaysRowNum,
		&set.SheetStartRow,
		&set.SheetApartmentColumNum,
	)
	if err != nil {
		return nil, err
	}
	return set, nil
}
