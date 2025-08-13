package test

import (
	"awesomeProject/internal/entities"
	"awesomeProject/internal/services"
	"context"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// --- Моки StorageSettings ---

type MockStorageSettings struct {
	mock.Mock
}

func (m *MockStorageSettings) Get(ctx context.Context, sheetName string) (*entities.Settings, error) {
	return nil, nil
}

func (m *MockStorageSettings) Set(ctx context.Context, settings entities.Settings) (*entities.Settings, error) {
	return nil, nil
}
func (m *MockStorageSettings) SettingsUpdate(ctx context.Context, settings entities.Settings) (*entities.Settings, error) {
	return nil, nil
}

func Test_synchronizeSliceFromDBandExcel(t *testing.T) {
	ctx := context.Background()
	mockStorageSettings := new(MockStorageSettings)

	mockStorageReservation := new(MockStorageReservation)

	mockStorageGuest := new(MockStorageGuest)

	serviceExcel := services.NewServiceExcel(
		mockStorageSettings,
		services.NewService(mockStorageReservation, mockStorageGuest),
	)

	bookingTimeIn := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	bookingTimeOut := time.Date(2025, 6, 5, 0, 0, 0, 0, time.UTC)

	bookingInDB := entities.Booking{

		Guest: entities.Guest{
			Phone: "123",
		},
		Reservation: entities.Reservation{
			RoomNumber: "101",
			CheckIn:    bookingTimeIn,
			CheckOut:   bookingTimeOut,
			Price:      1000,
		},
	}

	bookingInExcel := entities.Booking{
		Guest: entities.Guest{
			Phone: "123",
		},
		Reservation: entities.Reservation{
			RoomNumber: "101",
			CheckIn:    bookingTimeIn,
			CheckOut:   bookingTimeOut,
			Price:      1000,
		},
	}
	var expectedExcel []entities.Booking
	expectedDB := []entities.Booking{bookingInDB}

	t.Run("Удаляет совпадения и удаляет лишние из БД", func(t *testing.T) {
		bookingsFromDB := []entities.Booking{bookingInDB}
		bookingsFromExcel := []entities.Booking{bookingInExcel}

		// Ожидаем вызов Delete для bookingToDelete
		mockStorageReservation.On("Delete", ctx, 0).
			Return(nil, nil).Once()

		remainingExcel, remainingDB, err := serviceExcel.SynchronizeSliceFromDBandExcel(ctx, bookingsFromDB, bookingsFromExcel)

		assert.NoError(t, err)
		assert.Equal(t, expectedExcel, remainingExcel, "excel slice")
		assert.Equal(t, expectedDB, remainingDB, "DB slice")

	})

}
