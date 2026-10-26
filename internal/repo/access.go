package repo

import (
	"awesomeProject/internal/entities"
	"context"
	"database/sql"
	"errors"
)

// ReadAllowedRooms возвращает номера апартаментов, видимых пользователю.
func (db *Repository) ReadAllowedRooms(ctx context.Context, scope entities.ApartmentScope) ([]string, error) {
	rows, err := db.PostgreSQL.QueryContext(ctx,
		`SELECT a.room_number FROM apartments a WHERE a.room_number IS NOT NULL AND `+apartmentInScope,
		apartmentScopeArgs(scope)...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rooms := make([]string, 0)
	for rows.Next() {
		var room string
		if err := rows.Scan(&room); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, rows.Err()
}

// ReadRoomByReservationID возвращает номер апартамента брони. found=false, если брони нет.
func (db *Repository) ReadRoomByReservationID(ctx context.Context, reservationID int) (room string, found bool, err error) {
	return db.readRoom(ctx, "SELECT room_number FROM reservations WHERE id=$1", reservationID)
}

// ReadRoomByReservationInfoID возвращает номер апартамента брони, к которой относится reservation_info.
func (db *Repository) ReadRoomByReservationInfoID(ctx context.Context, reservationInfoID int) (room string, found bool, err error) {
	return db.readRoom(ctx,
		"SELECT r.room_number FROM reservation_info ri JOIN reservations r ON r.id = ri.reservation_id WHERE ri.id=$1",
		reservationInfoID)
}

func (db *Repository) readRoom(ctx context.Context, query string, id int) (string, bool, error) {
	var room string
	err := db.PostgreSQL.QueryRowContext(ctx, query, id).Scan(&room)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return room, true, nil
}
