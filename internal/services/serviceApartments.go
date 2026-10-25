package services

import (
	"awesomeProject/internal/entities"
	"context"

	"github.com/google/uuid"
)

type ServiceApartment struct {
	storageApartment StorageApartment
}

type StorageApartment interface {
	ReadApartmentAll(ctx context.Context, role string) ([]entities.Apartment, error)
	ReadApartmentsOwnedBy(ctx context.Context, userID uuid.UUID) ([]entities.Apartment, error)
}

func NewServiceApartment(storageApartment StorageApartment) *ServiceApartment {
	return &ServiceApartment{
		storageApartment,
	}
}

func (s *ServiceApartment) GetAllApartment(ctx context.Context, role string) ([]entities.Apartment, error) {
	apartment, err := s.storageApartment.ReadApartmentAll(ctx, role)
	if err != nil {
		return nil, err
	}
	return apartment, nil
}

// GetApartmentsFor возвращает апартаменты, доступные пользователю: свои (apartments.owned_only)
// или по роли через apartments.access.
func (s *ServiceApartment) GetApartmentsFor(ctx context.Context, user entities.CurrentUser) ([]entities.Apartment, error) {
	if user.Permissions.Has(entities.PermApartmentsOwnedOnly) {
		return s.storageApartment.ReadApartmentsOwnedBy(ctx, user.ID)
	}
	return s.GetAllApartment(ctx, user.Role)
}
