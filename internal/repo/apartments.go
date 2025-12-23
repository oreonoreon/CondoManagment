package repo

import (
	"awesomeProject/internal/entities"
	"context"
	"database/sql"
	"errors"
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

type Company struct {
	Name       string
	Apartments []string
	Users      []User
}

type Apartment struct {
	RoomNumber     string
	Description    string
	AirbnbCalendar string
}

type User struct {
	name string
	role Role
}

type Role struct {
	Name string
}
