package repo

import "context"

type Apartment struct {
	oid            int
	RoomNumber     string `json:"room_number"`
	Description    string
	AirbnbCalendar string
}

func (db *DBPostgreSQl) CreateApartment(ctx context.Context, apartment Apartment) error {
	_, err := db.PostgreSQL.QueryContext(ctx,
		"INSERT INTO Apartments (room_number, description,airbnb_calendar) VALUES ($1,$2,$3)",
		apartment.RoomNumber, apartment.Description, apartment.AirbnbCalendar)
	if err != nil {
		return err
	}
	return nil
}
func (db *DBPostgreSQl) UpdateApartment() {

}
func (db *DBPostgreSQl) DeleteApartment() {

}

func (db *DBPostgreSQl) ReadApartment(ctx context.Context, roomNumber string) (*Apartment, error) {
	apartment := new(Apartment)
	queryContext := db.PostgreSQL.QueryRowContext(ctx,
		"Select * from Apartments where room_number=$1",
		roomNumber)

	err := queryContext.Scan(&apartment.oid, &apartment.RoomNumber, &apartment.Description, &apartment.AirbnbCalendar)
	if err != nil {
		return nil, err
	}
	return apartment, nil
}

func (db *DBPostgreSQl) ReadApartmentAll(ctx context.Context) ([]Apartment, error) {
	apartments := make([]Apartment, 0)
	queryContext, err := db.PostgreSQL.QueryContext(ctx,
		"Select * from Apartments",
	)
	if err != nil {
		return nil, err
	}
	defer queryContext.Close()

	for queryContext.Next() {
		var apartment Apartment
		err = queryContext.Scan(&apartment.oid, &apartment.RoomNumber, &apartment.Description, &apartment.AirbnbCalendar)
		if err != nil {
			return nil, err
		}
		apartments = append(apartments, apartment)
	}

	return apartments, nil
}
