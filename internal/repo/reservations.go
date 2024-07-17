package repo

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Reservaton struct {
	Oid        int
	RoomNumber string    `json:"room_number"`
	GuestID    uuid.UUID `json:"guest_id"`
	CheckIn    time.Time `db:"check_in"`
	CheckOut   time.Time `db:"check_out"`
}

type Reservator interface {
	Delete(ctx context.Context)
	Create(ctx context.Context, r Reservaton) (*Reservaton, error)
	Read(ctx context.Context, checkin, checkout string) ([]Reservaton, error)
	ReadWithRoomNumber(ctx context.Context, roomNumber, checkin, checkout string) ([]Reservaton, error)
	Update(ctx context.Context, id int, checkIn, checkOut string) error
}

func (db *DBPostgreSQl) Delete(ctx context.Context) {

}

func (db *DBPostgreSQl) Create(ctx context.Context, r Reservaton) (*Reservaton, error) {
	reservation := new(Reservaton)
	queryContext := db.PostgreSQL.QueryRowContext(ctx,
		"INSERT INTO Reservations (room_number, guest_id, check_in, check_out) VALUES ($1, $2, $3,$4) Returning *",
		r.RoomNumber, r.GuestID, r.CheckIn, r.CheckOut)

	err := queryContext.Scan(&reservation.Oid, &reservation.RoomNumber, &reservation.GuestID, &reservation.CheckIn, &reservation.CheckOut)
	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (db *DBPostgreSQl) Read(ctx context.Context, checkin, checkout string) ([]Reservaton, error) {
	reservations := make([]Reservaton, 0, 0)
	queryContext, err := db.PostgreSQL.QueryContext(ctx,
		"Select * from Reservations where (check_in, check_out) OVERLAPS ($1, $2)",
		checkin, checkout)
	if err != nil {
		return nil, err
	}
	defer queryContext.Close()

	for queryContext.Next() {
		var reservation Reservaton
		err = queryContext.Scan(&reservation.Oid, &reservation.RoomNumber, &reservation.GuestID, &reservation.CheckIn, &reservation.CheckOut)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	return reservations, nil
}

func (db *DBPostgreSQl) ReadWithRoomNumber(ctx context.Context, roomNumber, checkin, checkout string) ([]Reservaton, error) {
	reservations := make([]Reservaton, 0, 0)
	queryContext, err := db.PostgreSQL.QueryContext(ctx,
		"Select * from Reservations where (check_in, check_out) OVERLAPS ($1, $2) AND room_number=$3",
		checkin, checkout, roomNumber)
	if err != nil {
		return nil, err
	}
	defer queryContext.Close()

	for queryContext.Next() {
		var reservation Reservaton
		err = queryContext.Scan(&reservation.Oid, &reservation.RoomNumber, &reservation.GuestID, &reservation.CheckIn, &reservation.CheckOut)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	return reservations, nil
}

func (db *DBPostgreSQl) Update(ctx context.Context, id int, checkIn, checkOut string) error {
	reservation := new(Reservaton)
	queryContext := db.PostgreSQL.QueryRowContext(ctx, "UPDATE Reservations SET check_in=$1, check_out=$2 where id=$3 Returning *",
		checkIn, checkOut, id)

	err := queryContext.Scan(&reservation.Oid, &reservation.RoomNumber, &reservation.GuestID, &reservation.CheckIn, &reservation.CheckOut)
	if err != nil {
		return err
	}
	return nil
}
