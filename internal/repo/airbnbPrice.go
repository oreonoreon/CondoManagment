package repo

import (
	"awesomeProject/internal/entities"
	"context"
	"database/sql"
	"errors"
)

func (db *Repository) GetAirbnbPrice(ctx context.Context, p entities.AirbnbPrice) (*entities.AirbnbPrice, error) {
	price := new(entities.AirbnbPrice)
	query := "Select * from airbnb_price where roomid=$1 AND check_in = $2 AND check_out=$3 AND scraping_date=$4"

	queryContext := db.PostgreSQL.QueryRowContext(
		ctx,
		query,
		p.RoomID,
		p.CheckIn,
		p.CheckOut,
		p.ScrapingDate,
	)

	err := queryContext.Scan(
		&price.RoomID,
		&price.Price,
		&price.CheckIn,
		&price.CheckOut,
		&price.ScrapingDate,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return price, nil
}

func (db *Repository) CreateAirbnbPrice(ctx context.Context, p entities.AirbnbPrice) (*entities.AirbnbPrice, error) {
	price := new(entities.AirbnbPrice)
	query := "INSERT INTO airbnb_price (roomid, price, check_in,check_out,scraping_date) VALUES ($1, $2, $3,$4,$5) RETURNING *"

	queryContext := db.PostgreSQL.QueryRowContext(
		ctx,
		query,
		p.RoomID,
		p.Price,
		p.CheckIn,
		p.CheckOut,
		p.ScrapingDate,
	)

	err := queryContext.Scan(
		&price.RoomID,
		&price.Price,
		&price.CheckIn,
		&price.CheckOut,
		&price.ScrapingDate,
	)
	if err != nil {
		return nil, err
	}

	return price, nil
}
