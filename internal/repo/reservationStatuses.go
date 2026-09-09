package repo

import (
	"awesomeProject/internal/entities"
	"context"
	"database/sql"
	"errors"
)

// GetActiveReservationStatuses возвращает список активных статусов для брони с подтянутыми названиями/цветом.
func (db *Repository) GetActiveReservationStatuses(ctx context.Context, reservationID int) ([]entities.ReservationStatus, error) {
	runner := getRunner(ctx, db.PostgreSQL)

	query := `SELECT rs.id, rs.reservation_id, rs.status_type_id, rs.is_active, rs.set_by, rs.set_at, rs.removed_at,
			st.code, st.name, st.color
		FROM reservation_statuses rs
		JOIN status_types st ON st.id = rs.status_type_id
		WHERE rs.reservation_id = $1 AND rs.is_active = TRUE
		ORDER BY st.sort_order`

	rows, err := runner.QueryContext(ctx, query, reservationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]entities.ReservationStatus, 0)
	for rows.Next() {
		var rs entities.ReservationStatus
		var setBy sql.NullString
		var removedAt sql.NullTime
		if err := rows.Scan(
			&rs.ID, &rs.ReservationID, &rs.StatusTypeID, &rs.IsActive, &setBy, &rs.SetAt, &removedAt,
			&rs.StatusCode, &rs.StatusName, &rs.StatusColor,
		); err != nil {
			return nil, err
		}
		if setBy.Valid {
			if id, err := parseUUID(setBy.String); err == nil {
				rs.SetBy = &id
			}
		}
		if removedAt.Valid {
			rs.RemovedAt = &removedAt.Time
		}
		result = append(result, rs)
	}
	return result, rows.Err()
}

// GetActiveReservationStatusByType возвращает текущую активную запись статуса указанного типа для брони, если есть.
func (db *Repository) GetActiveReservationStatusByType(ctx context.Context, reservationID, statusTypeID int) (*entities.ReservationStatus, error) {
	runner := getRunner(ctx, db.PostgreSQL)
	result := new(entities.ReservationStatus)

	query := `SELECT id, reservation_id, status_type_id, is_active, set_by, set_at, removed_at
		FROM reservation_statuses
		WHERE reservation_id = $1 AND status_type_id = $2 AND is_active = TRUE`

	row := runner.QueryRowContext(ctx, query, reservationID, statusTypeID)
	var setBy sql.NullString
	var removedAt sql.NullTime
	err := row.Scan(&result.ID, &result.ReservationID, &result.StatusTypeID, &result.IsActive, &setBy, &result.SetAt, &removedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if setBy.Valid {
		if id, err := parseUUID(setBy.String); err == nil {
			result.SetBy = &id
		}
	}
	if removedAt.Valid {
		result.RemovedAt = &removedAt.Time
	}
	return result, nil
}

// CreateReservationStatus включает статус для брони (новая запись истории).
func (db *Repository) CreateReservationStatus(ctx context.Context, rs entities.ReservationStatus) (*entities.ReservationStatus, error) {
	runner := getRunner(ctx, db.PostgreSQL)
	result := new(entities.ReservationStatus)

	query := `INSERT INTO reservation_statuses (reservation_id, status_type_id, is_active, set_by, set_at)
		VALUES ($1, $2, TRUE, $3, now())
		RETURNING id, reservation_id, status_type_id, is_active, set_by, set_at, removed_at`

	row := runner.QueryRowContext(ctx, query, rs.ReservationID, rs.StatusTypeID, rs.SetBy)
	var setBy sql.NullString
	var removedAt sql.NullTime
	if err := row.Scan(&result.ID, &result.ReservationID, &result.StatusTypeID, &result.IsActive, &setBy, &result.SetAt, &removedAt); err != nil {
		return nil, err
	}
	if setBy.Valid {
		if id, err := parseUUID(setBy.String); err == nil {
			result.SetBy = &id
		}
	}
	if removedAt.Valid {
		result.RemovedAt = &removedAt.Time
	}
	return result, nil
}

// DeactivateReservationStatus снимает активный статус (is_active=false), сохраняя историю.
func (db *Repository) DeactivateReservationStatus(ctx context.Context, id int) (*entities.ReservationStatus, error) {
	runner := getRunner(ctx, db.PostgreSQL)
	result := new(entities.ReservationStatus)

	query := `UPDATE reservation_statuses
		SET is_active = FALSE, removed_at = now()
		WHERE id = $1
		RETURNING id, reservation_id, status_type_id, is_active, set_by, set_at, removed_at`

	row := runner.QueryRowContext(ctx, query, id)
	var setBy sql.NullString
	var removedAt sql.NullTime
	if err := row.Scan(&result.ID, &result.ReservationID, &result.StatusTypeID, &result.IsActive, &setBy, &result.SetAt, &removedAt); err != nil {
		return nil, err
	}
	if setBy.Valid {
		if id, err := parseUUID(setBy.String); err == nil {
			result.SetBy = &id
		}
	}
	if removedAt.Valid {
		result.RemovedAt = &removedAt.Time
	}
	return result, nil
}
