package server

import (
	"awesomeProject/internal/entities"
	"awesomeProject/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	currentUserContextKey = "currentUser"
	roomAccessContextKey  = "roomAccess"
	errForbiddenRoom      = "access to this apartment is denied"
)

// currentUser возвращает пользователя запроса, загруженного SessionAuthMiddleware.
func currentUser(c *gin.Context) (entities.CurrentUser, bool) {
	value, ok := c.Get(currentUserContextKey)
	if !ok {
		return entities.CurrentUser{}, false
	}
	user, ok := value.(entities.CurrentUser)
	return user, ok
}

// mustCurrentUser — как currentUser, но при отсутствии пользователя пишет 401.
func mustCurrentUser(c *gin.Context) (entities.CurrentUser, bool) {
	user, ok := currentUser(c)
	if !ok {
		c.String(http.StatusUnauthorized, "authentication required")
	}
	return user, ok
}

// RequirePermission пропускает запрос, только если у роли пользователя есть все указанные права.
func RequirePermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := currentUser(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}
		for _, permission := range permissions {
			if !user.Permissions.Has(permission) {
				zap.L().Warn("RequirePermission: access denied",
					zap.String("path", c.Request.URL.Path),
					zap.String("role", user.Role),
					zap.String("required", permission),
				)
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
				return
			}
		}
		c.Next()
	}
}

// Уборки (/calendar/cleaning) доступны ролям с правом cleaning.manage, но намеренно
// не ограничиваются по апартаментам: такие роли видят и редактируют уборки по любым комнатам,
// включая отсутствующие в apartments.

// roomAccess возвращает доступные пользователю апартаменты (кешируется на время запроса).
// При ошибке сам пишет ответ и возвращает false.
func (h *Handle) roomAccess(c *gin.Context) (services.RoomAccess, bool) {
	if cached, ok := c.Get(roomAccessContextKey); ok {
		return cached.(services.RoomAccess), true
	}

	user, ok := mustCurrentUser(c)
	if !ok {
		return services.RoomAccess{}, false
	}

	access, err := h.access.RoomAccessFor(c.Request.Context(), user)
	if err != nil {
		zap.L().Error("roomAccess", zap.Error(err))
		c.String(http.StatusInternalServerError, "internal error")
		return services.RoomAccess{}, false
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
