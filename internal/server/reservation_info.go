package server

import (
	"awesomeProject/internal/entities"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// UpdateReservationInfo PATCH /calendar/reservation-info/:id
// Обновляет запись reservation_info по её id.
func (h *Handle) UpdateReservationInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid id: "+err.Error())
		return
	}

	var req entities.ReservationInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("UpdateReservationInfo: bind json", zap.Error(err))
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	req.ID = id

	result, err := h.TransactionalService.UpdateReservationInfoByID(c.Request.Context(), req)
	if err != nil {
		zap.L().Error("UpdateReservationInfo", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, result)
}
