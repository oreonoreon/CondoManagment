package repo

import (
	"awesomeProject/internal/entities"
	"context"
	"database/sql"
	"errors"
	"time"
)

const cleaningColumns = `id, reservation_id, cleaning_time, room, cleaning_price, laundry_price, agent_name, description, paid`

func scanCleaning(row *sql.Row, result *entities.Cleaning) error {
	return row.Scan(
		&result.ID,
		&result.ReservationID,
		&result.CleaningTime,
		&result.Room,
		&result.CleaningPrice,
		&result.LaundryPrice,
		&result.AgentName,
		&result.Description,
		&result.Paid,
	)
}

func (db *Repository) CreateCleaning(ctx context.Context, c entities.Cleaning) (*entities.Cleaning, error) {
	runner := getRunner(ctx, db.PostgreSQL)
	result := new(entities.Cleaning)

	query := `INSERT INTO cleaning (reservation_id, cleaning_time, room, cleaning_price, laundry_price, agent_name, description, paid)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING ` + cleaningColumns

	row := runner.QueryRowContext(ctx, query,
		c.ReservationID, c.CleaningTime, c.Room,
		c.CleaningPrice, c.LaundryPrice, c.AgentName, c.Description, c.Paid,
	)
	if err := scanCleaning(row, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (db *Repository) UpdateCleaning(ctx context.Context, c entities.Cleaning) (*entities.Cleaning, error) {
	runner := getRunner(ctx, db.PostgreSQL)
	result := new(entities.Cleaning)

	query := `UPDATE cleaning
		SET reservation_id=$2, cleaning_time=$3, room=$4, cleaning_price=$5, laundry_price=$6, agent_name=$7, description=$8, paid=$9
		WHERE id=$1
		RETURNING ` + cleaningColumns

	row := runner.QueryRowContext(ctx, query,
		c.ID, c.ReservationID, c.CleaningTime, c.Room,
		c.CleaningPrice, c.LaundryPrice, c.AgentName, c.Description, c.Paid,
	)
	if err := scanCleaning(row, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (db *Repository) DeleteCleaning(ctx context.Context, id int) (*entities.Cleaning, error) {
	runner := getRunner(ctx, db.PostgreSQL)
	result := new(entities.Cleaning)

	query := `DELETE FROM cleaning WHERE id=$1 RETURNING ` + cleaningColumns

	row := runner.QueryRowContext(ctx, query, id)
	err := scanCleaning(row, result)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *Repository) GetCleaningByID(ctx context.Context, id int) (*entities.Cleaning, error) {
	runner := getRunner(ctx, db.PostgreSQL)
	result := new(entities.Cleaning)

	query := `SELECT ` + cleaningColumns + ` FROM cleaning WHERE id=$1`

	row := runner.QueryRowContext(ctx, query, id)
	err := scanCleaning(row, result)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *Repository) GetAllCleaning(ctx context.Context) ([]entities.Cleaning, error) {
	runner := getRunner(ctx, db.PostgreSQL)

	query := `SELECT ` + cleaningColumns + ` FROM cleaning ORDER BY cleaning_time`

	rows, err := runner.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]entities.Cleaning, 0)
	for rows.Next() {
		var c entities.Cleaning
		if err := rows.Scan(
			&c.ID, &c.ReservationID, &c.CleaningTime, &c.Room,
			&c.CleaningPrice, &c.LaundryPrice, &c.AgentName, &c.Description, &c.Paid,
		); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, rows.Err()
}

func (db *Repository) GetCleaningByDate(ctx context.Context, date time.Time) ([]entities.Cleaning, error) {
	runner := getRunner(ctx, db.PostgreSQL)

	query := `SELECT ` + cleaningColumns + ` FROM cleaning WHERE DATE(cleaning_time) = $1 ORDER BY cleaning_time`

	rows, err := runner.QueryContext(ctx, query, date.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]entities.Cleaning, 0)
	for rows.Next() {
		var c entities.Cleaning
		if err := rows.Scan(
			&c.ID, &c.ReservationID, &c.CleaningTime, &c.Room,
			&c.CleaningPrice, &c.LaundryPrice, &c.AgentName, &c.Description, &c.Paid,
		); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, rows.Err()
}

func (db *Repository) GetCleaningByReservationID(ctx context.Context, reservationID int) (*entities.Cleaning, error) {
	runner := getRunner(ctx, db.PostgreSQL)
	result := new(entities.Cleaning)

	query := `SELECT ` + cleaningColumns + ` FROM cleaning WHERE reservation_id=$1`

	row := runner.QueryRowContext(ctx, query, reservationID)
	err := scanCleaning(row, result)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}
