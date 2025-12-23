package services

import (
	"awesomeProject/internal/entities"
	"context"
)

type ServiceApartment struct {
	storageApartment StorageApartment
}

type StorageApartment interface {
	ReadApartmentAll(ctx context.Context, role string) ([]entities.Apartment, error)
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
