package server

import (
	"awesomeProject/internal/entities"
	"encoding/json"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fullBooking — бронь со всеми заполненными полями (чтобы omitempty ничего не прятал).
func fullBooking() entities.Booking {
	id := uuid.New()
	return entities.Booking{
		Guest: entities.Guest{GuestID: id, Name: "Ivan", Phone: "+7", Description: "vip"},
		Reservation: entities.Reservation{
			Oid: 7, RoomNumber: "A605", GuestID: id,
			CheckIn: time.Date(2026, 5, 8, 14, 0, 0, 0, time.UTC), CheckOut: time.Date(2026, 5, 15, 11, 0, 0, 0, time.UTC),
			Price: 7000, CleaningPrice: 1500, ElectricityAndWaterPayment: "счётчики",
			Adult: 2, Children: 1, Description: "late", Days: 7, PriceForOneNight: 1000,
		},
		ReservationInfo: entities.ReservationInfo{ID: 1, ReservationID: 7, Deposit: 100},
		Statuses:        []entities.ReservationStatus{{ID: 1}},
	}
}

func keysOf(t *testing.T, view any) []string {
	raw, err := json.Marshal(view)
	require.NoError(t, err)
	var item map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(raw, &item))
	keys := make([]string, 0, len(item))
	for k := range item {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

func permsForFields(fields ...string) entities.Permissions {
	perms := entities.NewPermissions()
	for _, f := range fields {
		perms[entities.PermBookingFieldPrefix+f] = struct{}{}
	}
	return perms
}

func TestViewBookingAgent(t *testing.T) {
	perms := permsForFields("room_number", "guest_name", "check_in", "check_out", "adult", "children", "days")
	view, err := viewBooking(perms, fullBooking())
	require.NoError(t, err)
	assert.Equal(t, []string{"adult", "check_in", "check_out", "children", "days", "id", "name", "roomNumber"}, keysOf(t, view))
}

func TestViewBookingOwner(t *testing.T) {
	perms := permsForFields("guest_name", "guest_phone", "room_number", "check_in", "check_out", "price", "adult", "children")
	view, err := viewBooking(perms, fullBooking())
	require.NoError(t, err)
	assert.Equal(t, []string{"adult", "check_in", "check_out", "children", "id", "name", "phone", "price", "roomNumber"}, keysOf(t, view))
}

func TestViewBookingAllFieldsUnchanged(t *testing.T) {
	b := fullBooking()
	view, err := viewBooking(entities.NewPermissions(entities.PermBookingFieldAll), b)
	require.NoError(t, err)
	assert.Equal(t, b, view)
}

// Каждый JSON-ключ Booking должен быть привязан к какому-то booking.field.<поле>,
// иначе после переименования тега поле молча пропадёт у ограниченных ролей.
func TestBookingFieldKeysCoverBooking(t *testing.T) {
	all := entities.NewPermissions()
	for field := range bookingFieldKeys {
		all[entities.PermBookingFieldPrefix+field] = struct{}{}
	}
	b := fullBooking()
	full, err := json.Marshal(b)
	require.NoError(t, err)
	limited, err := viewBooking(all, b)
	require.NoError(t, err)
	limitedRaw, err := json.Marshal(limited)
	require.NoError(t, err)
	assert.JSONEq(t, string(full), string(limitedRaw))
}
