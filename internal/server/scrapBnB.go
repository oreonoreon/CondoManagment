package server

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handle) ScrapBnBPost(c *gin.Context) {
	type Request struct {
		Filter   string `json:"filter"`
		Location string `json:"location"`
		CheckIn  string `json:"checkIn"`
		CheckOut string `json:"checkOut"`
	}
	request := new(Request)
	err := c.BindJSON(request)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	data, err := h.ServiceBnB.ScrapDataFromAirBnB(c.Request.Context(), request.Filter, request.Location, request.CheckIn, request.CheckOut)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *Handle) ScrapBnBLocationNameUpdate(c *gin.Context) {
	type Request struct {
		RoomID       int64  `json:"roomID,string"`
		LocationName string `json:"locationName"`
	}
	request := new(Request)
	err := c.BindJSON(request)
	if err != nil {
		zap.L().Error("ScrapBnBLocationNameUpdate", zap.Error(err))
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	coordinate, err := h.ServiceBnB.UpdateLocationName(c.Request.Context(), request.RoomID, request.LocationName)
	if err != nil {
		zap.L().Error("ScrapBnBLocationNameUpdate", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, coordinate)
}

func (h *Handle) ScrapBnbRoomUnderstandableTypePatch(c *gin.Context) {
	type Request struct {
		RoomID             int64  `json:"roomID,string"`
		UnderstandableType string `json:"UnderstandableType"`
	}
	request := new(Request)
	err := c.BindJSON(request)
	if err != nil {
		zap.L().Error("ScrapBnbRoomUnderstandableTypePatch", zap.Error(err))
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	room, err := h.ServiceBnB.UpdateRoomUnderstandableType(c.Request.Context(), request.RoomID, request.UnderstandableType)
	if err != nil {
		zap.L().Error("ScrapBnbRoomUnderstandableTypePatch", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, room)
}
