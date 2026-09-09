package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

	fileData, err := h.TransactionalService.CreateReport(c.Request.Context(), request.RoomNumber, request.Start, request.End)
	if err != nil {
		zap.L().Error("CreateReport", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	//c.Header("Content-Disposition", "attachment; filename=\""+request.RoomNumber+".xlsx\"")
	//c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileData)

	fileName := request.RoomNumber + ".xlsx"

	// Устанавливаем заголовки ПЕРЕД отправкой данных
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Length", strconv.Itoa(len(fileData)))
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	// Используем Writer напрямую для бинарных данных
	c.Writer.WriteHeader(http.StatusOK)
	_, err = c.Writer.Write(fileData)
	if err != nil {
		zap.L().Error("Write file data", zap.Error(err))
	}
}

//func (h *Handle) MiddlePriceForPeriodReport(c *gin.Context) {
//	roleStr, err := getRoleFromContext(c)
//	if err != nil {
//		c.String(http.StatusUnauthorized, err.Error())
//		return
//	}
//
//	if c.Request.Method == http.MethodGet {
//		c.HTML(http.StatusOK, "middlepriceReport.html", nil)
//		return
//	}
//
//	start := c.PostForm("start")
//	end := c.PostForm("end")
//	apartments, err := h.ServiceApartment.GetAllApartment(c.Request.Context(), roleStr)
//	if err != nil {
//		zap.L().Error("FindMiddlePriceForPeriod", zap.Error(err))
//		c.String(http.StatusInternalServerError, err.Error())
//		return
//	}
//
//	priceMap, err := h.TransactionalService.FindMiddlePriceForPeriodReport(c.Request.Context(), apartments, start, end)
//	if err != nil {
//		zap.L().Error("FindMiddlePriceForPeriod", zap.Error(err))
//		c.String(http.StatusInternalServerError, err.Error())
//		return
//	}
//
//	c.JSON(http.StatusOK, priceMap)
//}

func (h *Handle) TotalPriceForPeriodReport(c *gin.Context) {
	roleStr, err := getRoleFromContext(c)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		return
	}

	type perioud struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	}

	p := new(perioud)
	err = c.BindJSON(p)
	if err != nil {
		zap.L().Error("UpdateBooking", zap.Error(err))
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	apartments, err := h.ServiceApartment.GetAllApartment(c.Request.Context(), roleStr)
	if err != nil {
		zap.L().Error("TotalPriceForPeriodReport", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	priceMap, err := h.TransactionalService.FindTotalPriceForPeriodReport(c.Request.Context(), apartments, p.Start, p.End)
	if err != nil {
		zap.L().Error("TotalPriceForPeriodReport", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, priceMap)
}

// totalPriceReportStartMonth/Year и totalPriceReportEndMonth/Year задают период (помесячно,
// включительно), за который строится xlsx отчёт суммарных цен по апартаментам.
const (
	totalPriceReportStartMonth = 11
	totalPriceReportStartYear  = 2023
	totalPriceReportEndMonth   = 12
	totalPriceReportEndYear    = 2028
)

// TotalPriceForPeriodReportXlsx отдаёт xlsx файл с суммарными ценами по каждому апартаменту
// помесячно за период с totalPriceReportStartMonth.totalPriceReportStartYear по
// totalPriceReportEndMonth.totalPriceReportEndYear.
func (h *Handle) TotalPriceForPeriodReportXlsx(c *gin.Context) {
	roleStr, err := getRoleFromContext(c)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		return
	}

	apartments, err := h.ServiceApartment.GetAllApartment(c.Request.Context(), roleStr)
	if err != nil {
		zap.L().Error("TotalPriceForPeriodReportXlsx", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	fileData, err := h.TransactionalService.TotalPriceForPeriodReportXlsx(
		c.Request.Context(),
		apartments,
		totalPriceReportStartMonth,
		totalPriceReportStartYear,
		totalPriceReportEndMonth,
		totalPriceReportEndYear,
	)
	if err != nil {
		zap.L().Error("TotalPriceForPeriodReportXlsx", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	fileName := "TotalPriceReport.xlsx"

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Length", strconv.Itoa(len(fileData)))
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	c.Writer.WriteHeader(http.StatusOK)
	if _, err := c.Writer.Write(fileData); err != nil {
		zap.L().Error("TotalPriceForPeriodReportXlsx/Write file data", zap.Error(err))
	}
}
