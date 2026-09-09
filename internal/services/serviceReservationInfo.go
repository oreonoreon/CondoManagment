package services

import (
	"awesomeProject/internal/entities"
	"context"
	"go.uber.org/zap"
	"time"
)

type ServiceReservationInfo struct {
	storage StorageReservationInfoRepo
}

type StorageReservationInfoRepo interface {
	CreateReservationInfo(ctx context.Context, ri entities.ReservationInfo) (*entities.ReservationInfo, error)
	GetReservationInfoByID(ctx context.Context, id int) (*entities.ReservationInfo, error)
	GetReservationInfoByReservationID(ctx context.Context, reservationID int) (*entities.ReservationInfo, error)
	UpdateReservationInfoByID(ctx context.Context, ri entities.ReservationInfo) (*entities.ReservationInfo, error)
	UpdateReservationInfoByReservationID(ctx context.Context, ri entities.ReservationInfo) (*entities.ReservationInfo, error)
	DeleteReservationInfo(ctx context.Context, id int) (*entities.ReservationInfo, error)
	GetReservationInfosByActualCheckIn(ctx context.Context, date time.Time) ([]entities.ReservationInfo, error)
	GetReservationInfosByActualCheckOut(ctx context.Context, date time.Time) ([]entities.ReservationInfo, error)
}

func NewServiceReservationInfo(storage StorageReservationInfoRepo) *ServiceReservationInfo {
	return &ServiceReservationInfo{storage: storage}
}

func (s *ServiceReservationInfo) CreateReservationInfo(ctx context.Context, ri entities.ReservationInfo) (*entities.ReservationInfo, error) {
	result, err := s.storage.CreateReservationInfo(ctx, ri)
	if err != nil {
		zap.L().Error("CreateReservationInfo", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (s *ServiceReservationInfo) GetReservationInfoByID(ctx context.Context, id int) (*entities.ReservationInfo, error) {
	result, err := s.storage.GetReservationInfoByID(ctx, id)
	if err != nil {
		zap.L().Error("GetReservationInfoByID", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (s *ServiceReservationInfo) GetReservationInfoByReservationID(ctx context.Context, reservationID int) (*entities.ReservationInfo, error) {
	result, err := s.storage.GetReservationInfoByReservationID(ctx, reservationID)
	if err != nil {
		zap.L().Error("GetReservationInfoByReservationID", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (s *ServiceReservationInfo) UpdateReservationInfoByID(ctx context.Context, ri entities.ReservationInfo) (*entities.ReservationInfo, error) {
	result, err := s.storage.UpdateReservationInfoByID(ctx, ri)
	if err != nil {
		zap.L().Error("UpdateReservationInfo", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (s *ServiceReservationInfo) UpdateReservationInfoByReservationID(ctx context.Context, ri entities.ReservationInfo) (*entities.ReservationInfo, error) {
	result, err := s.storage.UpdateReservationInfoByReservationID(ctx, ri)
	if err != nil {
		zap.L().Error("UpdateReservationInfo", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (s *ServiceReservationInfo) DeleteReservationInfo(ctx context.Context, id int) (*entities.ReservationInfo, error) {
	result, err := s.storage.DeleteReservationInfo(ctx, id)
	if err != nil {
		zap.L().Error("DeleteReservationInfo", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (s *ServiceReservationInfo) GetReservationInfosByActualCheckIn(ctx context.Context, date time.Time) ([]entities.ReservationInfo, error) {
	result, err := s.storage.GetReservationInfosByActualCheckIn(ctx, date)
	if err != nil {
		zap.L().Error("GetReservationInfosByActualCheckIn", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (s *ServiceReservationInfo) GetReservationInfosByActualCheckOut(ctx context.Context, date time.Time) ([]entities.ReservationInfo, error) {
	result, err := s.storage.GetReservationInfosByActualCheckOut(ctx, date)
	if err != nil {
		zap.L().Error("GetReservationInfosByActualCheckOut", zap.Error(err))
		return nil, err
	}
	return result, nil
}
