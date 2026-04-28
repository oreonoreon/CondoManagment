package repo

import (
	"awesomeProject/internal/entities"
	"context"
	"database/sql"
	"errors"
)

const reservationInfoColumns = `id, reservation_id, deposit, deposit_currency, prepayment, payment_on_checkin, notes`

func scanReservationInfo(row *sql.Row, result *entities.ReservationInfo) error {
	return row.Scan(
		&result.ID,
		&result.ReservationID,
		&result.Deposit,
		&result.DepositCurrency,
		&result.Prepayment,
		&result.PaymentOnCheckin,
		&result.Notes,
	)
}

func (db *Repository) CreateReservationInfo(ctx context.Context, ri entities.ReservationInfo) (*entities.ReservationInfo, error) {
	runner := getRunner(ctx, db.PostgreSQL)
	result := new(entities.ReservationInfo)

	query := `INSERT INTO reservation_info (reservation_id, deposit, deposit_currency, prepayment, payment_on_checkin, notes)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING ` + reservationInfoColumns

	row := runner.QueryRowContext(ctx, query,
		ri.ReservationID, ri.Deposit, ri.DepositCurrency,
		ri.Prepayment, ri.PaymentOnCheckin, ri.Notes,
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
		SET deposit=$2, deposit_currency=$3, prepayment=$4, payment_on_checkin=$5, notes=$6
		WHERE id=$1
		RETURNING ` + reservationInfoColumns

	row := runner.QueryRowContext(ctx, query,
		ri.ID, ri.Deposit, ri.DepositCurrency,
		ri.Prepayment, ri.PaymentOnCheckin, ri.Notes,
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
