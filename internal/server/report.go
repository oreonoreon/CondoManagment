package server

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handle) Report(c *gin.Context) {
	type ReportRequest struct {
		RoomNumber string `json:"room_number"`
		Start      string `json:"start"`
		End        string `json:"end"`
	}

	request := ReportRequest{}

	err := c.BindJSON(&request)
	if err != nil {
		zap.L().Error("CreateBookingPost", zap.Error(err))
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	fileData, err := h.Service.CreateReport(c.Request.Context(), request.RoomNumber, request.Start, request.End)
	if err != nil {
		zap.L().Error("CreateReport", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Disposition", "attachment; filename=\""+request.RoomNumber+".xlsx\"")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileData)
}

func (h *Handle) MiddlePriceForPeriodReport(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "middlepriceReport.html", nil)
		return
	}

	start := c.PostForm("start")
	end := c.PostForm("end")
	apartments, err := h.ServiceApartment.GetAllApartment(c.Request.Context())
	if err != nil {
		zap.L().Error("FindMiddlePriceForPeriod", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	priceMap, err := h.Service.FindMiddlePriceForPeriodReport(c.Request.Context(), apartments, start, end)
	if err != nil {
		zap.L().Error("FindMiddlePriceForPeriod", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, priceMap)
}

func (h *Handle) TotalPriceForPeriodReport(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "totalpriceReport.html", nil)
		return
	}

	start := c.PostForm("start")
	end := c.PostForm("end")
	apartments, err := h.ServiceApartment.GetAllApartment(c.Request.Context())
	if err != nil {
		zap.L().Error("TotalPriceForPeriodReport", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	priceMap, err := h.Service.FindTotalPriceForPeriodReport(c.Request.Context(), apartments, start, end)
	if err != nil {
		zap.L().Error("TotalPriceForPeriodReport", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "totalpriceReport.html", gin.H{
		"Info": priceMap,
	})
}
