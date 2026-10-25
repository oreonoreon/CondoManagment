package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

// ReadForbiddenRooms возвращает номера апартаментов, к которым у роли нет доступа
// (роль отсутствует в apartments.access).
func (db *Repository) ReadForbiddenRooms(ctx context.Context, role string) ([]string, error) {
	rows, err := db.PostgreSQL.QueryContext(ctx,
		"SELECT room_number FROM apartments WHERE NOT ($1 = ANY(COALESCE(access, '{}')))",
		role,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rooms := make([]string, 0)
	for rows.Next() {
		var room sql.NullString
		if err := rows.Scan(&room); err != nil {
			return nil, err
		}
		if room.Valid {
			rooms = append(rooms, room.String)
		}
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

// ReadOwnedRooms возвращает номера апартаментов, владельцем которых является пользователь.
func (db *Repository) ReadOwnedRooms(ctx context.Context, userID uuid.UUID) ([]string, error) {
	rows, err := db.PostgreSQL.QueryContext(ctx,
		"SELECT room_number FROM apartment_owners WHERE user_id = $1",
		userID,
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
