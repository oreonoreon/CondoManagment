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

	var bookingsFromExcel = []entities.Booking{
		{
			Guest: entities.Guest{
				Phone: "+79130317799",
			},
			Reservation: entities.Reservation{
				RoomNumber: "ME206",
				CheckIn:    time.Date(2025, time.July, 10, 0, 0, 0, 0, time.UTC),
				CheckOut:   time.Date(2025, time.July, 31, 0, 0, 0, 0, time.UTC),
				Price:      31500,
			},
		},
		{
			Guest: entities.Guest{
				Phone: "+79069168888",
			},
			Reservation: entities.Reservation{
				RoomNumber: "ME206",
				CheckIn:    time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC),
				CheckOut:   time.Date(2025, time.August, 11, 0, 0, 0, 0, time.UTC),
				Price:      14000,
			},
		},
		{
			Guest: entities.Guest{
				Phone: "+79036635367",
			},
			Reservation: entities.Reservation{
				RoomNumber: "ME206",
				CheckIn:    time.Date(2025, time.October, 14, 0, 0, 0, 0, time.UTC),
				CheckOut:   time.Date(2025, time.November, 14, 0, 0, 0, 0, time.UTC),
				Price:      0,
			},
		},
		{
			Guest: entities.Guest{
				Phone: "+79233274004",
			},
			Reservation: entities.Reservation{
				RoomNumber: "ME206",
				CheckIn:    time.Date(2026, time.February, 17, 0, 0, 0, 0, time.UTC),
				CheckOut:   time.Date(2026, time.March, 5, 0, 0, 0, 0, time.UTC),
				Price:      68800,
			},
		},
	}

	var bookingsFromDB = []entities.Booking{
		{
			Guest: entities.Guest{
				Phone: "+79130317799",
			},
			Reservation: entities.Reservation{
				RoomNumber: "ME206",
				CheckIn:    time.Date(2025, time.July, 10, 0, 0, 0, 0, time.UTC),
				CheckOut:   time.Date(2025, time.July, 31, 0, 0, 0, 0, time.UTC),
				Price:      31500, // суммарная цена (из price)
			},
		},
		{
			Guest: entities.Guest{
				Phone: "+79069168888",
			},
			Reservation: entities.Reservation{
				RoomNumber: "ME206",
				CheckIn:    time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC),
				CheckOut:   time.Date(2025, time.August, 11, 0, 0, 0, 0, time.UTC),
				Price:      0,
			},
		},
		{
			Guest: entities.Guest{
				Phone: "+79036635367",
			},
			Reservation: entities.Reservation{
				RoomNumber: "ME206",
				CheckIn:    time.Date(2025, time.October, 14, 0, 0, 0, 0, time.UTC),
				CheckOut:   time.Date(2025, time.November, 14, 0, 0, 0, 0, time.UTC),
				Price:      0,
			},
		},
		{
			Guest: entities.Guest{
				Phone: "+79233274004",
			},
			Reservation: entities.Reservation{
				RoomNumber: "ME206",
				CheckIn:    time.Date(2026, time.February, 17, 0, 0, 0, 0, time.UTC),
				CheckOut:   time.Date(2026, time.March, 5, 0, 0, 0, 0, time.UTC),
				Price:      68800,
			},
		},
	}

	expectedExcel := []entities.Booking{
		{
			Guest: entities.Guest{
				Phone: "+79069168888",
			},
			Reservation: entities.Reservation{
				RoomNumber: "ME206",
				CheckIn:    time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC),
				CheckOut:   time.Date(2025, time.August, 11, 0, 0, 0, 0, time.UTC),
				Price:      14000,
			},
		},
	}

	expectedDB := []entities.Booking{
		{
			Guest: entities.Guest{
				Phone: "+79130317799",
			},
			Reservation: entities.Reservation{
				RoomNumber: "ME206",
				CheckIn:    time.Date(2025, time.July, 10, 0, 0, 0, 0, time.UTC),
				CheckOut:   time.Date(2025, time.July, 31, 0, 0, 0, 0, time.UTC),
				Price:      31500, // суммарная цена (из price)
			},
		},
		{
			Guest: entities.Guest{
				Phone: "+79036635367",
			},
			Reservation: entities.Reservation{
				RoomNumber: "ME206",
				CheckIn:    time.Date(2025, time.October, 14, 0, 0, 0, 0, time.UTC),
				CheckOut:   time.Date(2025, time.November, 14, 0, 0, 0, 0, time.UTC),
				Price:      0,
			},
		},
		{
			Guest: entities.Guest{
				Phone: "+79233274004",
			},
			Reservation: entities.Reservation{
				RoomNumber: "ME206",
				CheckIn:    time.Date(2026, time.February, 17, 0, 0, 0, 0, time.UTC),
				CheckOut:   time.Date(2026, time.March, 5, 0, 0, 0, 0, time.UTC),
				Price:      68800,
			},
		},
	}

	t.Run("Удаляет совпадения и удаляет лишние из БД", func(t *testing.T) {
		//bookingsFromDB := []entities.Booking{bookingInDB}
		//bookingsFromExcel := []entities.Booking{bookingInExcel}

		// Ожидаем вызов Delete для bookingToDelete
		mockStorageReservation.On("Delete", ctx, 0).
			Return(nil, nil).Once()

		remainingExcel, remainingDB, err := serviceExcel.SynchronizeSliceFromDBandExcel(ctx, bookingsFromDB, bookingsFromExcel)

		assert.NoError(t, err)
		assert.Equal(t, expectedExcel, remainingExcel, "excel slice")
		assert.Equal(t, expectedDB, remainingDB, "DB slice")

	})

}
