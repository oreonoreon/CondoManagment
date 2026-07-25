package services

import (
	"awesomeProject/internal/entities"
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ServiceStatus struct {
	storage StorageStatusRepo
}

type StorageStatusRepo interface {
	ListStatusTypes(ctx context.Context) ([]entities.StatusType, error)
	GetActiveReservationStatuses(ctx context.Context, reservationID int) ([]entities.ReservationStatus, error)
	GetActiveReservationStatusByType(ctx context.Context, reservationID, statusTypeID int) (*entities.ReservationStatus, error)
	CreateReservationStatus(ctx context.Context, rs entities.ReservationStatus) (*entities.ReservationStatus, error)
	DeactivateReservationStatus(ctx context.Context, id int) (*entities.ReservationStatus, error)
}

func NewServiceStatus(storage StorageStatusRepo) *ServiceStatus {
	return &ServiceStatus{storage: storage}
}

// List возвращает справочник активных типов статусов (для чекбоксов/бейджей в UI).
func (s *ServiceStatus) List(ctx context.Context) ([]entities.StatusType, error) {
	result, err := s.storage.ListStatusTypes(ctx)
	if err != nil {
		zap.L().Error("ServiceStatus.List", zap.Error(err))
		return nil, err
	}
	return result, nil
}

// GetByReservation возвращает текущие активные статусы конкретной брони.
func (s *ServiceStatus) GetByReservation(ctx context.Context, reservationID int) ([]entities.ReservationStatus, error) {
	result, err := s.storage.GetActiveReservationStatuses(ctx, reservationID)
	if err != nil {
		zap.L().Error("ServiceStatus.GetByReservation", zap.Error(err))
		return nil, err
	}
	return result, nil
}

// Toggle переключает статус брони: если статус этого типа сейчас активен - снимает его (сохраняя историю),
// если не активен/отсутствует - включает новой записью с указанием, кто установил (setBy).
func (s *ServiceStatus) Toggle(ctx context.Context, reservationID, statusTypeID int, setBy uuid.UUID) (*entities.ReservationStatus, error) {
	existing, err := s.storage.GetActiveReservationStatusByType(ctx, reservationID, statusTypeID)
	if err != nil {
		zap.L().Error("ServiceStatus.Toggle: GetActiveReservationStatusByType", zap.Error(err))
		return nil, err
	}

	if existing != nil {
		result, err := s.storage.DeactivateReservationStatus(ctx, existing.ID)
		if err != nil {
			zap.L().Error("ServiceStatus.Toggle: DeactivateReservationStatus", zap.Error(err))
			return nil, err
		}
		return result, nil
	}

	newStatus := entities.ReservationStatus{
		ReservationID: reservationID,
		StatusTypeID:  statusTypeID,
		IsActive:      true,
		SetBy:         &setBy,
	}
	result, err := s.storage.CreateReservationStatus(ctx, newStatus)
	if err != nil {
		zap.L().Error("ServiceStatus.Toggle: CreateReservationStatus", zap.Error(err))
		return nil, err
	}
	return result, nil
}
