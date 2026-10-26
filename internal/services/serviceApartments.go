package services

import (
	"awesomeProject/internal/entities"
	"context"
)

type ServiceApartment struct {
	storageApartment StorageApartment
}

type StorageApartment interface {
	ReadApartmentsFor(ctx context.Context, scope entities.ApartmentScope) ([]entities.Apartment, error)
}

func NewServiceApartment(storageApartment StorageApartment) *ServiceApartment {
	return &ServiceApartment{
		storageApartment,
	}
}

// GetApartmentsFor возвращает апартаменты, видимые пользователю: выданные его роли или лично ему
// (все — при праве apartments.all), неактивные — только при праве apartments.inactive.
func (s *ServiceApartment) GetApartmentsFor(ctx context.Context, user entities.CurrentUser) ([]entities.Apartment, error) {
	return s.storageApartment.ReadApartmentsFor(ctx, user.ApartmentScope())
}
