package repo

import (
	"context"
	"github.com/google/uuid"
)

type DBSql interface {
	Reservator
	CreateApartment(ctx context.Context, apartment Apartment) error
	ReadApartment(ctx context.Context, roomNumber string) (*Apartment, error)
	ReadApartmentAll(ctx context.Context) ([]Apartment, error)
	CreateGuest(ctx context.Context, g Guest) (*Guest, error)
	ReadGuest(ctx context.Context, guestID uuid.UUID) (*Guest, error)
}
