package repo

import (
	"awesomeProject/internal/entities"
	"context"
)

func (db *Repository) CreateRatingAirbnbRoom(ctx context.Context, r entities.RatingAirbnbRoom) (*entities.RatingAirbnbRoom, error) {
	rating := new(entities.RatingAirbnbRoom)
	query := "INSERT INTO rating_airbnb_room (roomid, value, reviewcount) VALUES ($1, $2, $3) RETURNING *"

	queryContext := db.PostgreSQL.QueryRowContext(
		ctx,
		query,
		r.RoomID,
		r.Value,
		r.ReviewCount,
	)

	err := queryContext.Scan(
		&rating.RoomID,
		&rating.Value,
		&rating.ReviewCount,
	)
	if err != nil {
		return nil, err
	}

	return rating, nil
}
