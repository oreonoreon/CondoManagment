package services

import (
	"awesomeProject/internal/entities"
	"context"
	"go.uber.org/zap"
	"time"
)

type ServiceCleaning struct {
	storage StorageCleaningRepo
}

type StorageCleaningRepo interface {
	CreateCleaning(ctx context.Context, c entities.Cleaning) (*entities.Cleaning, error)
	UpdateCleaning(ctx context.Context, c entities.Cleaning) (*entities.Cleaning, error)
	DeleteCleaning(ctx context.Context, id int) (*entities.Cleaning, error)
	GetCleaningByID(ctx context.Context, id int) (*entities.Cleaning, error)
	GetAllCleaning(ctx context.Context) ([]entities.Cleaning, error)
	GetCleaningByDate(ctx context.Context, date time.Time) ([]entities.Cleaning, error)
}

func NewServiceCleaning(storage StorageCleaningRepo) *ServiceCleaning {
	return &ServiceCleaning{storage: storage}
}

func (s *ServiceCleaning) CreateCleaningManual(ctx context.Context, c entities.Cleaning) (*entities.Cleaning, error) {
	result, err := s.storage.CreateCleaning(ctx, c)
	if err != nil {
		zap.L().Error("CreateCleaningManual", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (s *ServiceCleaning) UpdateCleaning(ctx context.Context, c entities.Cleaning) (*entities.Cleaning, error) {
	result, err := s.storage.UpdateCleaning(ctx, c)
	if err != nil {
		zap.L().Error("UpdateCleaning", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (s *ServiceCleaning) DeleteCleaning(ctx context.Context, id int) (*entities.Cleaning, error) {
	result, err := s.storage.DeleteCleaning(ctx, id)
	if err != nil {
		zap.L().Error("DeleteCleaning", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (s *ServiceCleaning) GetCleaningByID(ctx context.Context, id int) (*entities.Cleaning, error) {
	result, err := s.storage.GetCleaningByID(ctx, id)
	if err != nil {
		zap.L().Error("GetCleaningByID", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (s *ServiceCleaning) GetAllCleaning(ctx context.Context) ([]entities.Cleaning, error) {
	result, err := s.storage.GetAllCleaning(ctx)
	if err != nil {
		zap.L().Error("GetAllCleaning", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (s *ServiceCleaning) GetCleaningByDate(ctx context.Context, date time.Time) ([]entities.Cleaning, error) {
	result, err := s.storage.GetCleaningByDate(ctx, date)
	if err != nil {
		zap.L().Error("GetCleaningByDate", zap.Error(err))
		return nil, err
	}
	return result, nil
}
