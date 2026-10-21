package server

import (
	"awesomeProject/internal/entities"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func bookingRoom(b entities.Booking) string { return b.Reservation.RoomNumber }

// GetBookingsByCheckIn GET /calendar/bookings/check-in/:date
// Возвращает слайс Booking у которых дата check_in совпадает с указанной.
// Формат даты: YYYY-MM-DD
func (h *Handle) GetBookingsByCheckIn(c *gin.Context) {
	date, ok := parseDateParam(c)
	if !ok {
		return
	}
	access, ok := h.roomAccess(c)
	if !ok {
		return
	}
	bookings, err := h.TransactionalService.GetBookingByCheckIn(c.Request.Context(), date)
	if err != nil {
		zap.L().Error("GetBookingsByCheckIn", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	view, ok := bookingsViewFor(c, filterByRoom(access, bookings, bookingRoom))
	if !ok {
		return
	}
	c.IndentedJSON(http.StatusOK, view)
}

// GetBookingsByCheckOut GET /calendar/bookings/check-out/:date
// Возвращает слайс Booking у которых дата check_out совпадает с указанной.
// Формат даты: YYYY-MM-DD
func (h *Handle) GetBookingsByCheckOut(c *gin.Context) {
	date, ok := parseDateParam(c)
	if !ok {
		return
	}
	access, ok := h.roomAccess(c)
	if !ok {
		return
	}
	bookings, err := h.TransactionalService.GetBookingByCheckOut(c.Request.Context(), date)
	if err != nil {
		zap.L().Error("GetBookingsByCheckOut", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	view, ok := bookingsViewFor(c, filterByRoom(access, bookings, bookingRoom))
	if !ok {
		return
	}
	c.IndentedJSON(http.StatusOK, view)
}

// parseDateParam парсит параметр :date формата YYYY-MM-DD из пути.
// При ошибке пишет 400 и возвращает false.
func parseDateParam(c *gin.Context) (time.Time, bool) {
	dateStr := c.Param("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid date format, expected YYYY-MM-DD: "+err.Error())
		return time.Time{}, false
	}
	return date, true
}
