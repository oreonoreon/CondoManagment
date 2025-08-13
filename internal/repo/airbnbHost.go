package repo

import (
	"awesomeProject/internal/entities"
	"context"
	"database/sql"
	"errors"
)

func (db *Repository) GetAirbnbHost(ctx context.Context, id string) (*entities.Host, error) {
	host := new(entities.Host)
	query := "Select * from airbnb_host where id=$1"

	queryContext := db.PostgreSQL.QueryRowContext(
		ctx,
		query,
		id,
	)

	err := queryContext.Scan(
		&host.ID,
		&host.Name,
		&host.Trustable,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return host, nil
}

func (db *Repository) CreateAirbnbHost(ctx context.Context, h entities.Host) (*entities.Host, error) {
	host := new(entities.Host)
	query := "INSERT INTO airbnb_host (id, name, trustable) VALUES ($1, $2, $3) RETURNING *"

	queryContext := db.PostgreSQL.QueryRowContext(
		ctx,
		query,
		h.ID,
		h.Name,
		h.Trustable,
	)

	err := queryContext.Scan(&host.ID, &host.Name, &host.Trustable)
	if err != nil {
		return nil, err
	}

	return host, nil
}
