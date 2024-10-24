package repo

import (
	"awesomeProject/internal/entities"
	"awesomeProject/internal/myLogger"
	"context"
	"database/sql"
	"github.com/google/uuid"
)

//type Guester interface {
//	CreateGuest(ctx context.Context, g entities.Guest) (*entities.Guest, error)
//	ReadGuest(ctx context.Context, guestID uuid.UUID) (*entities.Guest, error)
//	FindGuestByPhoneNumber(ctx context.Context, phone string) (*entities.Guest, error)
//}

func (db *Repository) CreateGuest(ctx context.Context, g entities.Guest) (*entities.Guest, error) {
	guest := new(entities.Guest)
	queryContext := db.PostgreSQL.QueryRowContext(ctx,
		"INSERT INTO Guests (guest_id, name, phone,description) VALUES ($1, $2, $3, $4) RETURNING *",
		g.GuestID, g.Name, Nullable(g.Phone), g.Description)

	var sPhone sql.NullString
	err := queryContext.Scan(&guest.GuestID, &guest.Name, &sPhone, &guest.Description)
	if err != nil {
		return nil, err
	}

	if sPhone.Valid {
		guest.Phone = sPhone.String
	} else {
		guest.Phone = ""
	}

	return guest, nil
}

func (db *Repository) ReadGuest(ctx context.Context, guestID uuid.UUID) (*entities.Guest, error) {
	guest := new(entities.Guest)
	queryContext := db.PostgreSQL.QueryRowContext(ctx,
		"Select * from Guests where guest_id=$1",
		guestID)

	var sPhone sql.NullString
	err := queryContext.Scan(&guest.GuestID, &guest.Name, &sPhone, &guest.Description)
	if err != nil {
		return nil, err
	}

	if sPhone.Valid {
		guest.Phone = sPhone.String
	} else {
		guest.Phone = ""
	}

	return guest, nil
}

func Nullable(field string) interface{} {
	if len(field) == 0 {
		return sql.NullString{}
	}
	return field
}

func (db *Repository) FindGuestByPhoneNumber(ctx context.Context, phone string) (*entities.Guest, error) {
	guest := new(entities.Guest)
	queryContext := db.PostgreSQL.QueryRowContext(ctx,
		"Select * from Guests where phone=$1",
		phone)

	err := queryContext.Scan(&guest.GuestID, &guest.Name, &guest.Phone, &guest.Description)
	if err == db.PostgreSQL.ErrNoRows() {
		myLogger.Logger.Println(err)
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return guest, nil
}
