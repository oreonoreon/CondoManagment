package server

import (
	"awesomeProject/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	roleAdmin = "admin"

	roomAccessContextKey = "roomAccess"
	errForbiddenRoom     = "access to this apartment is denied"
)

// RequireRole пропускает запрос только для пользователей с указанной ролью.
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleStr, err := getRoleFromContext(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if roleStr != role {
			zap.L().Warn("RequireRole: access denied",
				zap.String("path", c.Request.URL.Path),
				zap.String("role", roleStr),
				zap.String("required", role),
			)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}
		c.Next()
	}
}

// Уборки (/calendar/cleaning) намеренно не ограничиваются по apartments.access:
// любая роль видит и редактирует уборки по любым комнатам, включая отсутствующие в apartments.

// roomAccess возвращает закрытые для роли пользователя апартаменты (кешируется на время запроса).
// При ошибке сам пишет ответ и возвращает false.
func (h *Handle) roomAccess(c *gin.Context) (services.RoomAccess, bool) {
	if cached, ok := c.Get(roomAccessContextKey); ok {
		return cached.(services.RoomAccess), true
	}

	roleStr, err := getRoleFromContext(c)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		return nil, false
	}

	access, err := h.access.RoomAccessForRole(c.Request.Context(), roleStr)
	if err != nil {
		zap.L().Error("roomAccess", zap.Error(err))
		c.String(http.StatusInternalServerError, "internal error")
		return nil, false
	}

	c.Set(roomAccessContextKey, access)
	return access, true
}

// requireRooms проверяет доступ ко всем перечисленным апартаментам. При отказе пишет 403.
func (h *Handle) requireRooms(c *gin.Context, rooms ...string) bool {
	access, ok := h.roomAccess(c)
	if !ok {
		return false
	}
	for _, room := range rooms {
		if !access.Allowed(room) {
			zap.L().Warn("requireRooms: access denied", zap.String("path", c.Request.URL.Path), zap.String("room", room))
			c.String(http.StatusForbidden, errForbiddenRoom)
			return false
		}
	}
	return true
}

// requireReservation проверяет доступ к апартаменту брони. Несуществующую бронь пропускает —
// дальше хендлер сам вернёт "не найдено".
func (h *Handle) requireReservation(c *gin.Context, reservationID int) bool {
	room, found, err := h.access.RoomByReservationID(c.Request.Context(), reservationID)
	if err != nil {
		zap.L().Error("requireReservation", zap.Error(err))
		c.String(http.StatusInternalServerError, "internal error")
		return false
	}
	if !found {
		return true
	}
	return h.requireRooms(c, room)
}

// requireReservationInfo проверяет доступ к апартаменту брони, к которой относится reservation_info.
func (h *Handle) requireReservationInfo(c *gin.Context, reservationInfoID int) bool {
	room, found, err := h.access.RoomByReservationInfoID(c.Request.Context(), reservationInfoID)
	if err != nil {
		zap.L().Error("requireReservationInfo", zap.Error(err))
		c.String(http.StatusInternalServerError, "internal error")
		return false
	}
	if !found {
		return true
	}
	return h.requireRooms(c, room)
}

// filterByRoom оставляет только элементы с доступными апартаментами.
func filterByRoom[T any](access services.RoomAccess, items []T, room func(T) string) []T {
	result := make([]T, 0, len(items))
	for _, item := range items {
		if access.Allowed(room(item)) {
			result = append(result, item)
		}
	}
	return result
}
