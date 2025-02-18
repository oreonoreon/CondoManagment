package services

import (
	"awesomeProject/internal/entities"
	"context"
)

type ServiceSettings struct {
	storageSettings StorageSettings
}

type StorageSettings interface {
	Get(ctx context.Context, sheetName string) (*entities.Settings, error)
	Set(ctx context.Context, settings entities.Settings) (*entities.Settings, error)
	SettingsUpdate(ctx context.Context, settings entities.Settings) (*entities.Settings, error)
}

func NewServiceSettings(storageSettings StorageSettings) *ServiceSettings {
	return &ServiceSettings{storageSettings: storageSettings}
}

type SettingsConverter interface {
	SettingsConvert() entities.Settings
}

func (s *ServiceSettings) Set(ctx context.Context, converter SettingsConverter) (*entities.Settings, error) {
	settings := converter.SettingsConvert()

	set, err := s.storageSettings.Get(ctx, settings.SheetName)
	if err != nil {
		return nil, err
	}

	if *set != settings {
		set, err = s.storageSettings.SettingsUpdate(ctx, settings)
		if err != nil {
			return nil, err
		}
	}

	return set, err
}
