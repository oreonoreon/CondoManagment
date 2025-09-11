package test

import (
	"awesomeProject/internal/services"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"awesomeProject/internal/entities"
)

// --- Mock для StorageReservation ---
type MockStorageReservation struct {
	mock.Mock
}

func (m *MockStorageReservation) UpdateReservation(ctx context.Context, r entities.Reservation) (*entities.Reservation, error) {
	return nil, nil
}
func (m *MockStorageReservation) Create(ctx context.Context, r entities.Reservation) (*entities.Reservation, error) {
	return nil, nil
}
func (m *MockStorageReservation) ReadALLByRoomNumber(ctx context.Context, roomNumber string) ([]entities.Reservation, error) {
	return nil, nil
}
func (m *MockStorageReservation) ReadWithRoomNumber(ctx context.Context, roomNumber string, checkin, checkout time.Time) ([]entities.Reservation, error) {
	return nil, nil
}
func (m *MockStorageReservation) FindBookingByGuestUUID(ctx context.Context, uuid uuid.UUID) ([]entities.Reservation, error) {
	return nil, nil
}
func (m *MockStorageReservation) Delete(ctx context.Context, id int) (*entities.Reservation, error) {
	return nil, nil
}
func (m *MockStorageReservation) GetReservationByID(ctx context.Context, id int) (*entities.Reservation, error) {
	return nil, nil
}

// --- Mock для StorageGuest ---
type MockStorageGuest struct {
	mock.Mock
}

func (m *MockStorageGuest) UpdateGuest(ctx context.Context, g entities.Guest) (*entities.Guest, error) {
	return nil, nil
}
func (m *MockStorageGuest) CreateGuest(ctx context.Context, g entities.Guest) (*entities.Guest, error) {
	return nil, nil
}
func (m *MockStorageGuest) FindGuestByPhoneNumber(ctx context.Context, phone string) (*entities.Guest, error) {
	return nil, nil
}
func (m *MockStorageGuest) ReadGuest(ctx context.Context, guestID uuid.UUID) (*entities.Guest, error) {
	return nil, nil
}

func TestFindMiddlePriceForPeriod_CasesIndividually(t *testing.T) {
	type testCase struct {
		name         string
		bookings     []entities.Reservation
		expectedDays int
		expectedAvg  int
	}

	start := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC)

	tests := []testCase{
		{
			name: "checkIn < start, checkOut < end \n checkIn > start, checkOut > end\n checkIn > start && checkOut < end\n checkIn < start && checkOut > end",
			bookings: []entities.Reservation{
				// Case A: checkIn < start, checkOut < end
				{
					CheckIn:          time.Date(2025, 5, 30, 0, 0, 0, 0, time.UTC),
					CheckOut:         time.Date(2025, 6, 3, 0, 0, 0, 0, time.UTC),
					PriceForOneNight: 1000,
					Price:            3000,
				},
				// Case B: checkIn > start, checkOut > end
				{
					CheckIn:          time.Date(2025, 6, 8, 0, 0, 0, 0, time.UTC),
					CheckOut:         time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
					PriceForOneNight: 1500,
					Price:            10500,
				},
				// Case C: checkIn > start && checkOut < end
				{
					CheckIn:          time.Date(2025, 6, 4, 0, 0, 0, 0, time.UTC),
					CheckOut:         time.Date(2025, 6, 6, 0, 0, 0, 0, time.UTC),
					PriceForOneNight: 1200,
					Price:            2400,
				},
				// Case D: checkIn < start && checkOut > end
				{
					CheckIn:          time.Date(2025, 5, 30, 0, 0, 0, 0, time.UTC),
					CheckOut:         time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
					PriceForOneNight: 1100,
					Price:            16500,
				},
			},
			expectedDays: 15,
			expectedAvg:  17300 / 15,
		},

		{
			name: "Case A: checkIn < start, checkOut == end",
			bookings: []entities.Reservation{
				{
					CheckIn:          time.Date(2025, 5, 25, 0, 0, 0, 0, time.UTC),
					CheckOut:         end,
					PriceForOneNight: 1000,
					Price:            15000,
				},
			},
			expectedDays: 9,
			expectedAvg:  1000,
		},
		{
			name: "Case B: checkIn == start, checkOut > end",
			bookings: []entities.Reservation{
				{
					CheckIn:          start,
					CheckOut:         time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
					PriceForOneNight: 1500,
					Price:            22500,
				},
			},
			expectedDays: 9,
			expectedAvg:  1500,
		},
		{
			name: "Case C: checkIn > start, checkOut < end",
			bookings: []entities.Reservation{
				{
					CheckIn:          time.Date(2025, 6, 2, 0, 0, 0, 0, time.UTC),
					CheckOut:         time.Date(2025, 6, 5, 0, 0, 0, 0, time.UTC),
					PriceForOneNight: 1200,
					Price:            3600,
				},
			},
			expectedDays: 3,
			expectedAvg:  1200,
		},
		{
			name: "Case D: checkIn < start, checkOut > end",
			bookings: []entities.Reservation{
				{
					CheckIn:          time.Date(2025, 5, 25, 0, 0, 0, 0, time.UTC),
					CheckOut:         time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
					PriceForOneNight: 1100,
					Price:            22000,
				},
			},
			expectedDays: 9,
			expectedAvg:  1100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRes := new(MockStorageReservation)
			mockGuest := new(MockStorageGuest)
			s := services.NewService(mockRes, mockGuest)

			mockRes.On("ReadWithRoomNumber", mock.Anything, "101", start, end).
				Return(tt.bookings, nil)

			avg, err := s.FindMiddlePriceForPeriod(context.Background(), "101", "2025-06-01", "2025-06-10")

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedAvg, avg)
		})
	}
}
