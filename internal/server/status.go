package server

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// GetStatusTypes GET /calendar/status-types
// Возвращает справочник активных типов статусов для отображения чекбоксов/бейджей в UI.
func (h *Handle) GetStatusTypes(c *gin.Context) {
	statusTypes, err := h.TransactionalService.ListStatusTypes(c.Request.Context())
	if err != nil {
		zap.L().Error("GetStatusTypes", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, statusTypes)
}

// GetReservationStatuses GET /calendar/reservations/:id/statuses
// Возвращает текущие активные статусы указанной брони.
func (h *Handle) GetReservationStatuses(c *gin.Context) {
	reservationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid id: "+err.Error())
		return
	}

	statuses, err := h.TransactionalService.GetReservationStatuses(c.Request.Context(), reservationID)
	if err != nil {
		zap.L().Error("GetReservationStatuses", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, statuses)
}

// ToggleReservationStatus POST /calendar/reservations/:id/statuses/:statusTypeId
// Переключает статус брони: включает, если сейчас не активен, либо снимает (сохраняя историю), если активен.
func (h *Handle) ToggleReservationStatus(c *gin.Context) {
	reservationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid id: "+err.Error())
		return
	}
	statusTypeID, err := strconv.Atoi(c.Param("statusTypeId"))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid statusTypeId: "+err.Error())
		return
	}

	rawUserID, exists := c.Get("userID")
	if !exists {
		c.String(http.StatusUnauthorized, "authentication required")
		return
	}
	userIDStr, ok := rawUserID.(string)
	if !ok {
		c.String(http.StatusInternalServerError, "invalid userID in session")
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		zap.L().Error("ToggleReservationStatus: parse userID", zap.Error(err))
		c.String(http.StatusInternalServerError, "invalid userID in session")
		return
	}

	result, err := h.TransactionalService.ToggleReservationStatus(c.Request.Context(), reservationID, statusTypeID, userID)
	if err != nil {
		zap.L().Error("ToggleReservationStatus", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, result)
}
