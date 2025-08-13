package repo

import (
	"awesomeProject/internal/entities"
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

func (db *Repository) CreateAirbnbRoom(ctx context.Context, r entities.AirbnbRoom) (*entities.AirbnbRoom, error) {
	room := new(entities.AirbnbRoom)
	query := "INSERT INTO airbnb_room (roomid, badges, name, title, type, kind, category, images, understandabletype, hostid) VALUES ($1, $2, $3,$4,$5,$6,$7,$8,$9,$10) RETURNING *"
	queryContext := db.PostgreSQL.QueryRowContext(
		ctx,
		query,
		r.RoomID,
		pq.Array(r.Badges),
		r.Name,
		r.Title,
		r.TypeOfEstate,
		r.Kind,
		r.Category,
		pq.Array(r.Images),
		r.UnderstandableType,
		r.HostID,
	)

	err := queryContext.Scan(
		&room.RoomID,
		pq.Array(&room.Badges),
		&room.Name,
		&room.Title,
		&room.TypeOfEstate,
		&room.Kind,
		&room.Category,
		pq.Array(&room.Images),
		&room.UnderstandableType,
		&room.HostID,
	)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (db *Repository) GetAirbnbRoom(ctx context.Context, roomID int64) (*entities.AirbnbRoom, error) {
	room := new(entities.AirbnbRoom)
	query := "Select * from airbnb_room where roomid=$1"
	queryContext := db.PostgreSQL.QueryRowContext(
		ctx,
		query,
		roomID,
	)

	err := queryContext.Scan(
		&room.RoomID,
		pq.Array(&room.Badges),
		&room.Name,
		&room.Title,
		&room.TypeOfEstate,
		&room.Kind,
		&room.Category,
		pq.Array(&room.Images),
		&room.UnderstandableType,
		&room.HostID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (db *Repository) UpdateAirbnbRoomUnderstandableType(ctx context.Context, roomID int64, understandableType string) (*entities.AirbnbRoom, error) {
	room := new(entities.AirbnbRoom)
	query := "UPDATE airbnb_room SET understandabletype=$1 where roomid=$2 RETURNING *"

	queryContext := db.PostgreSQL.QueryRowContext(
		ctx,
		query,
		understandableType,
		roomID,
	)

	err := queryContext.Scan(
		&room.RoomID,
		pq.Array(&room.Badges),
		&room.Name,
		&room.Title,
		&room.TypeOfEstate,
		&room.Kind,
		&room.Category,
		pq.Array(&room.Images),
		&room.UnderstandableType,
		&room.HostID,
	)
	//if errors.Is(err, sql.ErrNoRows) {
	//	return nil, nil
	//}
	if err != nil {
		return nil, err
	}

	return room, nil
}
