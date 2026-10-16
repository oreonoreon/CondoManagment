package server

import (
	"awesomeProject/internal/entities"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

func cleaningRoom(cl entities.Cleaning) string { return cl.Room }

// requireCleaningPayload проверяет доступ к комнате и брони, указанным в теле запроса.
func (h *Handle) requireCleaningPayload(c *gin.Context, req entities.Cleaning) bool {
	if !h.requireRooms(c, req.Room) {
		return false
	}
	if req.ReservationID != nil && !h.requireReservation(c, *req.ReservationID) {
		return false
	}
	return true
}

// requireCleaning проверяет доступ к существующей записи уборки. Несуществующую пропускает.
func (h *Handle) requireCleaning(c *gin.Context, id int) bool {
	cleaning, err := h.TransactionalService.GetCleaningByID(c.Request.Context(), id)
	if err != nil {
		zap.L().Error("requireCleaning", zap.Error(err))
		c.String(http.StatusInternalServerError, "internal error")
		return false
	}
	if cleaning == nil {
		return true
	}
	return h.requireRooms(c, cleaning.Room)
}

// GetAllCleaning GET /calendar/cleaning
func (h *Handle) GetAllCleaning(c *gin.Context) {
	access, ok := h.roomAccess(c)
	if !ok {
		return
	}
	cleanings, err := h.TransactionalService.GetAllCleaning(c.Request.Context())
	if err != nil {
		zap.L().Error("GetAllCleaning", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, filterByRoom(access, cleanings, cleaningRoom))
}

// GetCleaningByID GET /calendar/cleaning/:id
func (h *Handle) GetCleaningByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid id: "+err.Error())
		return
	}
	cleaning, err := h.TransactionalService.GetCleaningByID(c.Request.Context(), id)
	if err != nil {
		zap.L().Error("GetCleaningByID", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if cleaning == nil {
		c.String(http.StatusNotFound, "cleaning record not found")
		return
	}
	if !h.requireRooms(c, cleaning.Room) {
		return
	}
	c.IndentedJSON(http.StatusOK, cleaning)
}

// CreateCleaning POST /calendar/cleaning
func (h *Handle) CreateCleaning(c *gin.Context) {
	var req entities.Cleaning
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("CreateCleaning: bind json", zap.Error(err))
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if !h.requireCleaningPayload(c, req) {
		return
	}
	cleaning, err := h.TransactionalService.CreateCleaningManual(c.Request.Context(), req)
	if err != nil {
		zap.L().Error("CreateCleaning", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusCreated, cleaning)
}

// UpdateCleaning PATCH /calendar/cleaning/:id
func (h *Handle) UpdateCleaning(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid id: "+err.Error())
		return
	}
	var req entities.Cleaning
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("UpdateCleaning: bind json", zap.Error(err))
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	req.ID = id
	if !h.requireCleaning(c, id) || !h.requireCleaningPayload(c, req) {
		return
	}
	cleaning, err := h.TransactionalService.UpdateCleaning(c.Request.Context(), req)
	if err != nil {
		zap.L().Error("UpdateCleaning", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, cleaning)
}

// DeleteCleaning DELETE /calendar/cleaning/:id
func (h *Handle) DeleteCleaning(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid id: "+err.Error())
		return
	}
	if !h.requireCleaning(c, id) {
		return
	}
	cleaning, err := h.TransactionalService.DeleteCleaning(c.Request.Context(), id)
	if err != nil {
		zap.L().Error("DeleteCleaning", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if cleaning == nil {
		c.String(http.StatusNotFound, "cleaning record not found")
		return
	}
	c.IndentedJSON(http.StatusOK, cleaning)
}

// GetCleaningByDate GET /calendar/cleaning/date/:date
// Формат даты: 2006-01-02
func (h *Handle) GetCleaningByDate(c *gin.Context) {
	dateStr := c.Param("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid date format, expected YYYY-MM-DD: "+err.Error())
		return
	}
	access, ok := h.roomAccess(c)
	if !ok {
		return
	}
	cleanings, err := h.TransactionalService.GetCleaningByDate(c.Request.Context(), date)
	if err != nil {
		zap.L().Error("GetCleaningByDate", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, filterByRoom(access, cleanings, cleaningRoom))
}
