package services

import (
	"awesomeProject/internal/entities"
	"context"

	"github.com/google/uuid"
)

// RoomAccess описывает, к каким апартаментам у пользователя есть доступ.
//
// Обычный режим: rooms — апартаменты, закрытые для роли (роль не указана в apartments.access);
// всё остальное, включая комнаты, которых нет в apartments, доступно.
// Режим ownedOnly (право apartments.owned_only): rooms — апартаменты пользователя
// из apartment_owners, и доступны только они.
type RoomAccess struct {
	ownedOnly bool
	rooms     map[string]struct{}
}

func NewForbiddenRoomAccess(forbidden ...string) RoomAccess {
	return RoomAccess{rooms: toSet(forbidden)}
}

func NewOwnedRoomAccess(owned ...string) RoomAccess {
	return RoomAccess{ownedOnly: true, rooms: toSet(owned)}
}

func (a RoomAccess) Allowed(room string) bool {
	_, listed := a.rooms[room]
	if a.ownedOnly {
		return listed
	}
	return !listed
}

func toSet(items []string) map[string]struct{} {
	set := make(map[string]struct{}, len(items))
	for _, item := range items {
		set[item] = struct{}{}
	}
	return set
}

type ServiceAccess struct {
	storage StorageAccess
}

type StorageAccess interface {
	ReadForbiddenRooms(ctx context.Context, role string) ([]string, error)
	ReadOwnedRooms(ctx context.Context, userID uuid.UUID) ([]string, error)
	ReadRoomByReservationID(ctx context.Context, reservationID int) (string, bool, error)
	ReadRoomByReservationInfoID(ctx context.Context, reservationInfoID int) (string, bool, error)
}

func NewServiceAccess(storage StorageAccess) *ServiceAccess {
	return &ServiceAccess{storage: storage}
}

func (s *ServiceAccess) RoomAccessFor(ctx context.Context, user entities.CurrentUser) (RoomAccess, error) {
	if user.Permissions.Has(entities.PermApartmentsOwnedOnly) {
		owned, err := s.storage.ReadOwnedRooms(ctx, user.ID)
		if err != nil {
			return RoomAccess{}, err
		}
		return NewOwnedRoomAccess(owned...), nil
	}

	forbidden, err := s.storage.ReadForbiddenRooms(ctx, user.Role)
	if err != nil {
		return RoomAccess{}, err
	}
	return NewForbiddenRoomAccess(forbidden...), nil
}

func (s *ServiceAccess) RoomByReservationID(ctx context.Context, reservationID int) (string, bool, error) {
	return s.storage.ReadRoomByReservationID(ctx, reservationID)
}

func (s *ServiceAccess) RoomByReservationInfoID(ctx context.Context, reservationInfoID int) (string, bool, error) {
	return s.storage.ReadRoomByReservationInfoID(ctx, reservationInfoID)
}
