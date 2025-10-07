package repo

import (
	"awesomeProject/internal/entities"
	"awesomeProject/internal/erro"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

func (db *Repository) GetReservationByID(ctx context.Context, id int) (*entities.Reservation, error) {
	reservation := new(entities.Reservation)
	query := "Select * FROM Reservations WHERE id=$1"
	dbRow := db.PostgreSQL.QueryRowContext(ctx, query, id)
	err := dbRow.Scan(
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
		&reservation.Days,
		&reservation.PriceForOneNight,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return reservation, nil
}
func (db *Repository) Delete(ctx context.Context, id int) (*entities.Reservation, error) {
	reservation := new(entities.Reservation)
	query := "DELETE FROM Reservations WHERE id=$1 Returning *"
	dbRow := db.PostgreSQL.QueryRowContext(ctx, query, id)
	err := dbRow.Scan(
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
		&reservation.Days,
		&reservation.PriceForOneNight,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return reservation, nil
}

func (db *Repository) Create(ctx context.Context, r entities.Reservation) (*entities.Reservation, error) {
	runner := getRunner(ctx, db.PostgreSQL) // todo такое использование контекста надо переделать или ввести повсеместно
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
		"days," +
		"price_for_night" +
		") " +
		"VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) Returning *"

	queryContext := runner.QueryRowContext(
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
		&reservation.Days,
		&reservation.PriceForOneNight,
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

func (db *Repository) ReadALLByRoomNumber(ctx context.Context, roomNumber string) ([]entities.Reservation, error) {
	reservations := make([]entities.Reservation, 0, 0)
	queryContext, err := db.PostgreSQL.QueryContext(ctx,
		"Select * from Reservations where room_number=$1 ORDER BY check_in",
		roomNumber)
	if err != nil {
		return nil, err
	}
	defer queryContext.Close()

	for queryContext.Next() {
		var reservation entities.Reservation
		err = queryContext.Scan(
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
			&reservation.Days,
			&reservation.PriceForOneNight,
		)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	return reservations, nil
}

func (db *Repository) ReadWithRoomNumber(ctx context.Context, roomNumber string, checkin, checkout time.Time) ([]entities.Reservation, error) {
	reservations := make([]entities.Reservation, 0, 0)
	queryContext, err := db.PostgreSQL.QueryContext(ctx,
		"Select * from Reservations where (check_in, check_out) OVERLAPS ($1, $2) AND room_number=$3 ORDER BY check_in",
		checkin, checkout, roomNumber)
	if err != nil {
		return nil, err
	}
	defer queryContext.Close()

	for queryContext.Next() {
		var reservation entities.Reservation
		err = queryContext.Scan(
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
			&reservation.Days,
			&reservation.PriceForOneNight,
		)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	return reservations, nil
}

func (db *Repository) UpdateReservation(ctx context.Context, r entities.Reservation) (*entities.Reservation, error) {
	tx := From(ctx)
	if tx == nil {
		return nil, errors.New("context doesn't contain transaction") // todo решить делать так все запросы или через getRunner
	}

	reservation := new(entities.Reservation)
	query := "UPDATE Reservations SET room_number=$2, guest_id=$3, check_in=$4, check_out=$5, price=$6, cleaning_price=$7,electricity_and_water_payment=$8,adult=$9,children=$10,description=$11, days=$12,price_for_night=$13 where id=$1 Returning *"
	queryContext := tx.QueryRowContext(ctx, query,
		r.Oid, r.RoomNumber, r.GuestID, r.CheckIn, r.CheckOut, r.Price, r.CleaningPrice, r.ElectricityAndWaterPayment, r.Adult, r.Children, r.Description, r.Days, r.PriceForOneNight)

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
		&reservation.Days,
		&reservation.PriceForOneNight,
	)
	if err != nil {
		return nil, translatePQ(err)
	}
	return reservation, nil
}

func (db *Repository) FindBookingByGuestUUID(ctx context.Context, uuid uuid.UUID) ([]entities.Reservation, error) {
	return nil, nil
}

// Ошибки Postgres → доменные
func translatePQ(err error) error {
	var pqe *pq.Error
	if errors.As(err, &pqe) {
		switch pqe.Code {
		case "23P01": // exclusion_violation (пересечение диапазонов)
			return erro.ErrMatchWithOtherBooking
		}
	}
	return err
}
