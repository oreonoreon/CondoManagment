package services

import (
	"awesomeProject/internal/entities"
	"context"
	"go.uber.org/zap"
)

type ServiceReservationInfo struct {
	storage StorageReservationInfoRepo
}

type StorageReservationInfoRepo interface {
	CreateReservationInfo(ctx context.Context, ri entities.ReservationInfo) (*entities.ReservationInfo, error)
	GetReservationInfoByID(ctx context.Context, id int) (*entities.ReservationInfo, error)
	GetReservationInfoByReservationID(ctx context.Context, reservationID int) (*entities.ReservationInfo, error)
	UpdateReservationInfo(ctx context.Context, ri entities.ReservationInfo) (*entities.ReservationInfo, error)
	DeleteReservationInfo(ctx context.Context, id int) (*entities.ReservationInfo, error)
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

func (s *ServiceReservationInfo) UpdateReservationInfo(ctx context.Context, ri entities.ReservationInfo) (*entities.ReservationInfo, error) {
	result, err := s.storage.UpdateReservationInfo(ctx, ri)
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
