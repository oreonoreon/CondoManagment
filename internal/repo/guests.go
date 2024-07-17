package repo

import (
	"context"
	"github.com/google/uuid"
)

type Guest struct {
	GuestID     uuid.UUID
	Name        string
	Phone       string
	Description string
}

func (db *DBPostgreSQl) CreateGuest(ctx context.Context, g Guest) (*Guest, error) {
	guest := new(Guest)
	queryContext := db.PostgreSQL.QueryRowContext(ctx,
		"INSERT INTO Guests (guest_id, name, phone,description) VALUES ($1, $2, $3, $4) RETURNING *",
		g.GuestID, g.Name, g.Phone, g.Description)

	err := queryContext.Scan(&guest.GuestID, &guest.Name, &guest.Phone, &guest.Description)
	if err != nil {
		return nil, err
	}
	return guest, nil
}

func (db *DBPostgreSQl) ReadGuest(ctx context.Context, guestID uuid.UUID) (*Guest, error) {
	guest := new(Guest)
	queryContext := db.PostgreSQL.QueryRowContext(ctx,
		"Select * from Guests where guest_id=$1",
		guestID)

	err := queryContext.Scan(&guest.GuestID, &guest.Name, &guest.Phone, &guest.Description)
	if err != nil {
		return nil, err
	}
	return guest, nil
}
