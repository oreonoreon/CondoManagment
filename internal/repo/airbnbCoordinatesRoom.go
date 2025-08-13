package repo

import (
	"awesomeProject/internal/entities"
	"context"
	"database/sql"
	"errors"
)

func (db *Repository) CreateCoordinatesAirbnbRoom(ctx context.Context, c entities.CoordinatesAirbnbRoom) (*entities.CoordinatesAirbnbRoom, error) {
	coordinates := new(entities.CoordinatesAirbnbRoom)
	query := "INSERT INTO coordinates_airbnb_room (roomid, latitude, longitude, locationname) VALUES ($1, $2, $3, $4) RETURNING *"

	queryContext := db.PostgreSQL.QueryRowContext(
		ctx,
		query,
		c.RoomID,
		c.Latitude,
		c.Longitude,
		c.LocationName,
	)

	err := queryContext.Scan(
		&coordinates.RoomID,
		&coordinates.Latitude,
		&coordinates.Longitude,
		&coordinates.LocationName,
	)
	if err != nil {
		return nil, err
	}

	return coordinates, nil
}

func (db *Repository) UpdateCoordinatesAirbnbRoom(ctx context.Context, roomID int64, locationName string) (*entities.CoordinatesAirbnbRoom, error) {
	coordinates := new(entities.CoordinatesAirbnbRoom)
	query := "UPDATE coordinates_airbnb_room SET locationname=$1 where roomid=$2 RETURNING *"

	queryContext := db.PostgreSQL.QueryRowContext(
		ctx,
		query,
		locationName,
		roomID,
	)

	err := queryContext.Scan(
		&coordinates.RoomID,
		&coordinates.Latitude,
		&coordinates.Longitude,
		&coordinates.LocationName,
	)
	//if errors.Is(err, sql.ErrNoRows) {
	//	return nil, nil
	//}
	if err != nil {
		return nil, err
	}

	return coordinates, nil
}

func (db *Repository) GetCoordinatesAirbnbRoom(ctx context.Context, roomID int64) (*entities.CoordinatesAirbnbRoom, error) {
	coordinates := new(entities.CoordinatesAirbnbRoom)
	query := "Select * from coordinates_airbnb_room where roomid=$1"

	queryContext := db.PostgreSQL.QueryRowContext(
		ctx,
		query,
		roomID,
	)

	err := queryContext.Scan(
		&coordinates.RoomID,
		&coordinates.Latitude,
		&coordinates.Longitude,
		&coordinates.LocationName,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return coordinates, nil
}
