package repo

import (
	"awesomeProject/internal/entities"
	"context"
	"database/sql"
	"errors"
	"time"
)

const reservationInfoColumns = `id, reservation_id, deposit, deposit_currency, prepayment, payment_on_checkin, actual_check_in, actual_check_out`

func scanReservationInfo(row *sql.Row, result *entities.ReservationInfo) error {
	return row.Scan(
		&result.ID,
		&result.ReservationID,
		&result.Deposit,
		&result.DepositCurrency,
		&result.Prepayment,
		&result.PaymentOnCheckin,
		&result.ActualCheckIn,
		&result.ActualCheckOut,
	)
}

func (db *Repository) CreateReservationInfo(ctx context.Context, ri entities.ReservationInfo) (*entities.ReservationInfo, error) {
	runner := getRunner(ctx, db.PostgreSQL)
	result := new(entities.ReservationInfo)

	query := `INSERT INTO reservation_info (reservation_id, deposit, deposit_currency, prepayment, payment_on_checkin, actual_check_in, actual_check_out)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING ` + reservationInfoColumns

	row := runner.QueryRowContext(ctx, query,
		ri.ReservationID, ri.Deposit, ri.DepositCurrency,
		ri.Prepayment, ri.PaymentOnCheckin,
		ri.ActualCheckIn, ri.ActualCheckOut,
	)
	if err := scanReservationInfo(row, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (db *Repository) GetReservationInfoByReservationID(ctx context.Context, reservationID int) (*entities.ReservationInfo, error) {
	runner := getRunner(ctx, db.PostgreSQL)
	result := new(entities.ReservationInfo)

	query := `SELECT ` + reservationInfoColumns + ` FROM reservation_info WHERE reservation_id = $1`

	row := runner.QueryRowContext(ctx, query, reservationID)
	err := scanReservationInfo(row, result)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *Repository) GetReservationInfoByID(ctx context.Context, id int) (*entities.ReservationInfo, error) {
	runner := getRunner(ctx, db.PostgreSQL)
	result := new(entities.ReservationInfo)

	query := `SELECT ` + reservationInfoColumns + ` FROM reservation_info WHERE id = $1`

	row := runner.QueryRowContext(ctx, query, id)
	err := scanReservationInfo(row, result)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *Repository) UpdateReservationInfo(ctx context.Context, ri entities.ReservationInfo) (*entities.ReservationInfo, error) {
	runner := getRunner(ctx, db.PostgreSQL)
	result := new(entities.ReservationInfo)

	query := `UPDATE reservation_info
		SET deposit=$2, deposit_currency=$3, prepayment=$4, payment_on_checkin=$5, actual_check_in=$6, actual_check_out=$7
		WHERE id=$1
		RETURNING ` + reservationInfoColumns

	row := runner.QueryRowContext(ctx, query,
		ri.ID, ri.Deposit, ri.DepositCurrency,
		ri.Prepayment, ri.PaymentOnCheckin,
		ri.ActualCheckIn, ri.ActualCheckOut,
	)
	if err := scanReservationInfo(row, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (db *Repository) UpdateReservationInfoByReservationID(ctx context.Context, ri entities.ReservationInfo) (*entities.ReservationInfo, error) {
	runner := getRunner(ctx, db.PostgreSQL)
	result := new(entities.ReservationInfo)

	query := `UPDATE reservation_info
		SET deposit=$2, deposit_currency=$3, prepayment=$4, payment_on_checkin=$5, actual_check_in=$6, actual_check_out=$7
		WHERE reservation_id=$1
		RETURNING ` + reservationInfoColumns

	row := runner.QueryRowContext(ctx, query,
		ri.ReservationID, ri.Deposit, ri.DepositCurrency,
		ri.Prepayment, ri.PaymentOnCheckin,
		ri.ActualCheckIn, ri.ActualCheckOut,
	)
	if err := scanReservationInfo(row, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (db *Repository) DeleteReservationInfo(ctx context.Context, id int) (*entities.ReservationInfo, error) {
	runner := getRunner(ctx, db.PostgreSQL)
	result := new(entities.ReservationInfo)

	query := `DELETE FROM reservation_info WHERE id=$1 RETURNING ` + reservationInfoColumns

	row := runner.QueryRowContext(ctx, query, id)
	err := scanReservationInfo(row, result)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *Repository) GetReservationInfosByActualCheckIn(ctx context.Context, date time.Time) ([]entities.ReservationInfo, error) {
	return db.getReservationInfosByDate(ctx, "actual_check_in", date)
}

func (db *Repository) GetReservationInfosByActualCheckOut(ctx context.Context, date time.Time) ([]entities.ReservationInfo, error) {
	return db.getReservationInfosByDate(ctx, "actual_check_out", date)
}

// getReservationInfosByDate выбирает записи reservation_info по дате (без учёта времени) указанного поля.
func (db *Repository) getReservationInfosByDate(ctx context.Context, column string, date time.Time) ([]entities.ReservationInfo, error) {
	runner := getRunner(ctx, db.PostgreSQL)

	query := `SELECT ` + reservationInfoColumns + ` FROM reservation_info
		WHERE DATE(` + column + `) = $1`

	rows, err := runner.QueryContext(ctx, query, date.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]entities.ReservationInfo, 0)
	for rows.Next() {
		var ri entities.ReservationInfo
		if err := rows.Scan(
			&ri.ID,
			&ri.ReservationID,
			&ri.Deposit,
			&ri.DepositCurrency,
			&ri.Prepayment,
			&ri.PaymentOnCheckin,
			&ri.ActualCheckIn,
			&ri.ActualCheckOut,
		); err != nil {
			return nil, err
		}
		result = append(result, ri)
	}
	return result, rows.Err()
}
