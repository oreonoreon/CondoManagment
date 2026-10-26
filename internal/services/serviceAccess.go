package services

import (
	"awesomeProject/internal/entities"
	"context"
)

// RoomAccess — апартаменты, видимые пользователю (см. entities.ApartmentScope).
type RoomAccess struct {
	allowed map[string]struct{}
}

func NewRoomAccess(allowed ...string) RoomAccess {
	set := make(map[string]struct{}, len(allowed))
	for _, room := range allowed {
		set[room] = struct{}{}
	}
	return RoomAccess{allowed: set}
}

func (a RoomAccess) Allowed(room string) bool {
	_, ok := a.allowed[room]
	return ok
}

type ServiceAccess struct {
	storage StorageAccess
}

type StorageAccess interface {
	ReadAllowedRooms(ctx context.Context, scope entities.ApartmentScope) ([]string, error)
	ReadRoomByReservationID(ctx context.Context, reservationID int) (string, bool, error)
	ReadRoomByReservationInfoID(ctx context.Context, reservationInfoID int) (string, bool, error)
}

func NewServiceAccess(storage StorageAccess) *ServiceAccess {
	return &ServiceAccess{storage: storage}
}

func (s *ServiceAccess) RoomAccessFor(ctx context.Context, user entities.CurrentUser) (RoomAccess, error) {
	rooms, err := s.storage.ReadAllowedRooms(ctx, user.ApartmentScope())
	if err != nil {
		return RoomAccess{}, err
	}
	return NewRoomAccess(rooms...), nil
}

func (s *ServiceAccess) RoomByReservationID(ctx context.Context, reservationID int) (string, bool, error) {
	return s.storage.ReadRoomByReservationID(ctx, reservationID)
}

func (s *ServiceAccess) RoomByReservationInfoID(ctx context.Context, reservationInfoID int) (string, bool, error) {
	return s.storage.ReadRoomByReservationInfoID(ctx, reservationInfoID)
}
