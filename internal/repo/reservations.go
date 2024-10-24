package repo

import (
	"awesomeProject/internal/entities"
	"context"
)

//type Reservator interface {
//	Delete(ctx context.Context)
//	Create(ctx context.Context, r entities.Reservation) (*entities.Reservation, error)
//	Read(ctx context.Context, checkin, checkout string) ([]entities.Reservation, error)
//	ReadWithRoomNumber(ctx context.Context, roomNumber, checkin, checkout string) ([]entities.Reservation, error)
//	Update(ctx context.Context, id int, checkIn, checkOut string) error
//}

func (db *Repository) Delete(ctx context.Context) {

}

func (db *Repository) Create(ctx context.Context, r entities.Reservation) (*entities.Reservation, error) {
	reservation := new(entities.Reservation)

	query := "INSERT INTO Reservations (" +
		"room_number," +
		" guest_id," +
		" check_in," +
		" check_out," +
		"price," +
		"cleaning_price," +
		"electricity_and_water_payment," +
		"adult," +
		"children," +
		"description," +
		"days,price_for_night" +
		") " +
		"VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) Returning *"

	queryContext := db.PostgreSQL.QueryRowContext(
		ctx,
		query,
		r.RoomNumber,
		r.GuestID,
		r.CheckIn,
		r.CheckOut,
		r.Price,
		r.CleaningPrice,
		r.ElectricityAndWaterPayment,
		r.Adult,
		r.Children,
		r.Description,
		r.Days,
		r.PriceForOneNight,
	)

	err := queryContext.Scan(
		&reservation.Oid,
		&reservation.RoomNumber,
		&reservation.GuestID,
		&reservation.CheckIn,
		&reservation.CheckOut,
		&reservation.Price,
		&reservation.CleaningPrice,
		&reservation.ElectricityAndWaterPayment,
		&reservation.Adult,
		&reservation.Children,
		&reservation.Description,
	)
	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (db *Repository) Read(ctx context.Context, checkin, checkout string) ([]entities.Reservation, error) {
	reservations := make([]entities.Reservation, 0, 0)
	queryContext, err := db.PostgreSQL.QueryContext(ctx,
		"Select * from Reservations where (check_in, check_out) OVERLAPS ($1, $2)",
		checkin, checkout)
	if err != nil {
		return nil, err
	}
	defer queryContext.Close()

	for queryContext.Next() {
		var reservation entities.Reservation
		err = queryContext.Scan(&reservation.Oid, &reservation.RoomNumber, &reservation.GuestID, &reservation.CheckIn, &reservation.CheckOut)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	return reservations, nil
}

func (db *Repository) ReadWithRoomNumber(ctx context.Context, roomNumber, checkin, checkout string) ([]entities.Reservation, error) {
	reservations := make([]entities.Reservation, 0, 0)
	queryContext, err := db.PostgreSQL.QueryContext(ctx,
		"Select * from Reservations where (check_in, check_out) OVERLAPS ($1, $2) AND room_number=$3",
		checkin, checkout, roomNumber)
	if err != nil {
		return nil, err
	}
	defer queryContext.Close()

	for queryContext.Next() {
		var reservation entities.Reservation
		err = queryContext.Scan(&reservation.Oid, &reservation.RoomNumber, &reservation.GuestID, &reservation.CheckIn, &reservation.CheckOut)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	return reservations, nil
}

func (db *Repository) Update(ctx context.Context, id int, checkIn, checkOut string) error {
	reservation := new(entities.Reservation)
	queryContext := db.PostgreSQL.QueryRowContext(ctx, "UPDATE Reservations SET check_in=$1, check_out=$2 where id=$3 Returning *",
		checkIn, checkOut, id)

	err := queryContext.Scan(&reservation.Oid, &reservation.RoomNumber, &reservation.GuestID, &reservation.CheckIn, &reservation.CheckOut)
	if err != nil {
		return err
	}
	return nil
}
