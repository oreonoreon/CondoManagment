package services

import (
	"context"
)

// RoomAccess — множество апартаментов, закрытых для роли.
// Апартаменты, которых нет в таблице apartments (например, ручные уборки
// по сторонним комнатам), считаются доступными.
type RoomAccess map[string]struct{}

func (a RoomAccess) Allowed(room string) bool {
	_, forbidden := a[room]
	return !forbidden
}

type ServiceAccess struct {
	storage StorageAccess
}

type StorageAccess interface {
	ReadForbiddenRooms(ctx context.Context, role string) ([]string, error)
	ReadRoomByReservationID(ctx context.Context, reservationID int) (string, bool, error)
	ReadRoomByReservationInfoID(ctx context.Context, reservationInfoID int) (string, bool, error)
}

func NewServiceAccess(storage StorageAccess) *ServiceAccess {
	return &ServiceAccess{storage: storage}
}

func (s *ServiceAccess) RoomAccessForRole(ctx context.Context, role string) (RoomAccess, error) {
	rooms, err := s.storage.ReadForbiddenRooms(ctx, role)
	if err != nil {
		return nil, err
	}
	access := make(RoomAccess, len(rooms))
	for _, room := range rooms {
		access[room] = struct{}{}
	}
	return access, nil
}

func (s *ServiceAccess) RoomByReservationID(ctx context.Context, reservationID int) (string, bool, error) {
	return s.storage.ReadRoomByReservationID(ctx, reservationID)
}

func (s *ServiceAccess) RoomByReservationInfoID(ctx context.Context, reservationInfoID int) (string, bool, error) {
	return s.storage.ReadRoomByReservationInfoID(ctx, reservationInfoID)
}
