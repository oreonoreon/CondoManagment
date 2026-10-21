package server

import (
	"awesomeProject/internal/entities"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// bookingIDKey — id брони отдаётся всегда: без него календарь не может отрисовать событие.
const bookingIDKey = "id"

// bookingFieldKeys связывает поле из права booking.field.<поле> с JSON-ключами entities.Booking.
// Ключи, которых здесь нет, роли без booking.field.all не отдаются (белый список),
// поэтому новые поля Booking по умолчанию скрыты от ограниченных ролей.
var bookingFieldKeys = map[string][]string{
	entities.BookingFieldRoomNumber:       {"roomNumber"},
	entities.BookingFieldGuestID:          {"guest_id", "guest_uuid"},
	entities.BookingFieldGuestName:        {"name"},
	entities.BookingFieldGuestPhone:       {"phone"},
	entities.BookingFieldGuestDescription: {"guestDescription"},
	entities.BookingFieldCheckIn:          {"check_in"},
	entities.BookingFieldCheckOut:         {"check_out"},
	entities.BookingFieldPrice:            {"price"},
	entities.BookingFieldCleaningPrice:    {"cleaning_price"},
	entities.BookingFieldElectricity:      {"electricity_and_water_payment"},
	entities.BookingFieldAdult:            {"adult"},
	entities.BookingFieldChildren:         {"children"},
	entities.BookingFieldDescription:      {"reservationDescription"},
	entities.BookingFieldDays:             {"days"},
	entities.BookingFieldPriceForNight:    {"price_for_night"},
	entities.BookingFieldReservationInfo:  {"reservation_info"},
	entities.BookingFieldStatuses:         {"statuses"},
}

// allowedBookingKeys возвращает JSON-ключи брони, разрешённые правами.
func allowedBookingKeys(perms entities.Permissions) map[string]struct{} {
	allowed := map[string]struct{}{bookingIDKey: {}}
	for field, keys := range bookingFieldKeys {
		if perms.CanSeeBookingField(field) {
			for _, key := range keys {
				allowed[key] = struct{}{}
			}
		}
	}
	return allowed
}

// viewBookings отдаёт брони как есть для booking.field.all, иначе — только разрешённые поля.
func viewBookings(perms entities.Permissions, bookings []entities.Booking) (any, error) {
	if perms.CanSeeAllBookingFields() {
		return bookings, nil
	}

	raw, err := json.Marshal(bookings)
	if err != nil {
		return nil, err
	}
	var items []map[string]json.RawMessage
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, err
	}

	allowed := allowedBookingKeys(perms)
	for _, item := range items {
		for key := range item {
			if _, ok := allowed[key]; !ok {
				delete(item, key)
			}
		}
	}
	return items, nil
}

// bookingsViewFor возвращает брони в объёме, разрешённом правами пользователя запроса.
// При ошибке сам пишет ответ и возвращает false.
func bookingsViewFor(c *gin.Context, bookings []entities.Booking) (any, bool) {
	user, ok := mustCurrentUser(c)
	if !ok {
		return nil, false
	}
	view, err := viewBookings(user.Permissions, bookings)
	if err != nil {
		zap.L().Error("bookingsViewFor", zap.Error(err))
		c.String(http.StatusInternalServerError, "internal error")
		return nil, false
	}
	return view, true
}

// bookingViewFor — bookingsViewFor для одной брони.
func bookingViewFor(c *gin.Context, booking entities.Booking) (any, bool) {
	user, ok := mustCurrentUser(c)
	if !ok {
		return nil, false
	}
	view, err := viewBooking(user.Permissions, booking)
	if err != nil {
		zap.L().Error("bookingViewFor", zap.Error(err))
		c.String(http.StatusInternalServerError, "internal error")
		return nil, false
	}
	return view, true
}

// viewBooking — viewBookings для одной брони.
func viewBooking(perms entities.Permissions, booking entities.Booking) (any, error) {
	view, err := viewBookings(perms, []entities.Booking{booking})
	if err != nil {
		return nil, err
	}
	if items, ok := view.([]map[string]json.RawMessage); ok {
		return items[0], nil
	}
	return booking, nil
}
