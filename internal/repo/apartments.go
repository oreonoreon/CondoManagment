package repo

import (
	"awesomeProject/internal/entities"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

func (db *Repository) CreateApartment(ctx context.Context, apartment entities.Apartment) error {
	_, err := db.PostgreSQL.QueryContext(ctx,
		"INSERT INTO Apartments (room_number, description,airbnb_calendar) VALUES ($1,$2,$3)",
		apartment.RoomNumber, apartment.Description, apartment.AirbnbCalendar)
	if err != nil {
		return err
	}
	return nil
}
func (db *Repository) UpdateApartment() {

}
func (db *Repository) DeleteApartment() {

}

func (db *Repository) ReadApartment(ctx context.Context, roomNumber string) (*entities.Apartment, error) {
	apartment := new(entities.Apartment)
	queryContext := db.PostgreSQL.QueryRowContext(ctx,
		"Select * from Apartments where room_number=$1",
		roomNumber)

	err := queryContext.Scan(&apartment.Oid, &apartment.RoomNumber, &apartment.Description, &apartment.AirbnbCalendar)
	if err != nil {
		return nil, err
	}
	return apartment, nil
}

func (db *Repository) ReadApartmentAll(ctx context.Context, role string) ([]entities.Apartment, error) {
	apartments := make([]entities.Apartment, 0)
	queryContext, err := db.PostgreSQL.QueryContext(ctx,
		"Select id,room_number,description,airbnb_calendar from Apartments where $1 = ANY(access) ORDER BY room_number",
		role,
	)
	if err != nil {
		return nil, err
	}
	defer queryContext.Close()

	for queryContext.Next() {
		var apartment entities.Apartment
		var airbnbCalendarNull sql.NullString
		err = queryContext.Scan(&apartment.Oid, &apartment.RoomNumber, &apartment.Description, &airbnbCalendarNull)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}

		if airbnbCalendarNull.Valid {
			apartment.AirbnbCalendar = airbnbCalendarNull.String
		} else {
			apartment.AirbnbCalendar = ""
		}

		apartments = append(apartments, apartment)
	}

	return apartments, nil
}

// ReadApartmentsOwnedBy возвращает апартаменты, владельцем которых является пользователь.
func (db *Repository) ReadApartmentsOwnedBy(ctx context.Context, userID uuid.UUID) ([]entities.Apartment, error) {
	rows, err := db.PostgreSQL.QueryContext(ctx,
		`SELECT a.id, a.room_number, a.description, a.airbnb_calendar
		FROM apartments a
		JOIN apartment_owners ao ON ao.room_number = a.room_number
		WHERE ao.user_id = $1
		ORDER BY a.room_number`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	apartments := make([]entities.Apartment, 0)
	for rows.Next() {
		var apartment entities.Apartment
		var description, airbnbCalendar sql.NullString
		if err := rows.Scan(&apartment.Oid, &apartment.RoomNumber, &description, &airbnbCalendar); err != nil {
			return nil, err
		}
		apartment.Description = description.String
		apartment.AirbnbCalendar = airbnbCalendar.String
		apartments = append(apartments, apartment)
	}
	return apartments, rows.Err()
}
