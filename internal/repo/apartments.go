package repo

import (
	"awesomeProject/internal/entities"
	"context"
)

//type Apartmenter interface {
//	CreateApartment(ctx context.Context, apartment entities.Apartment) error
//	ReadApartment(ctx context.Context, roomNumber string) (*entities.Apartment, error)
//	ReadApartmentAll(ctx context.Context) ([]entities.Apartment, error)
//}

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

func (db *Repository) ReadApartmentAll(ctx context.Context) ([]entities.Apartment, error) {
	apartments := make([]entities.Apartment, 0)
	queryContext, err := db.PostgreSQL.QueryContext(ctx,
		"Select * from Apartments",
	)
	if err != nil {
		return nil, err
	}
	defer queryContext.Close()

	for queryContext.Next() {
		var apartment entities.Apartment
		err = queryContext.Scan(&apartment.Oid, &apartment.RoomNumber, &apartment.Description, &apartment.AirbnbCalendar)
		if err != nil {
			return nil, err
		}
		apartments = append(apartments, apartment)
	}

	return apartments, nil
}
