package repo

import (
	"context"
	"github.com/google/uuid"
)

//type DBSql interface {
//	Update(ctx context.Context, id int, checkIn, checkOut string) error
//	Read(ctx context.Context, checkin, checkout string) ([]Reservaton, error)
//	ReadWithRoomNumber(ctx context.Context, roomNumber, checkin, checkout string) ([]Reservaton, error)
//	Create(ctx context.Context, r Reservaton) (*Reservaton, error)
//	CreateApartment(ctx context.Context, apartment Apartment) error
//	ReadApartment(ctx context.Context, roomNumber string) (*Apartment, error)
//	ReadApartmentAll(ctx context.Context) ([]Apartment, error)
//	CreateGuest(ctx context.Context, g Guest) (*Guest, error)
//	ReadGuest(ctx context.Context, guestID uuid.UUID) (*Guest, error)
//}

type DBSql interface {
	Reservator
	CreateApartment(ctx context.Context, apartment Apartment) error
	ReadApartment(ctx context.Context, roomNumber string) (*Apartment, error)
	ReadApartmentAll(ctx context.Context) ([]Apartment, error)
	CreateGuest(ctx context.Context, g Guest) (*Guest, error)
	ReadGuest(ctx context.Context, guestID uuid.UUID) (*Guest, error)
}
