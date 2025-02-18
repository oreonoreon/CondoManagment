package server

import (
	"awesomeProject/internal/services"
	"fmt"
	"go.uber.org/zap"
	"html/template"
	"net/http"
)

type Handle struct {
	*services.Service
	*services.ServiceSettings
	*services.ServiceExcel
}

func NewHandle(serviceReservation *services.Service, serviceSettings *services.ServiceSettings, serviceExcel *services.ServiceExcel) Handle {
	return Handle{
		serviceReservation,
		serviceSettings,
		serviceExcel,
	}
}

func (h *Handle) SynchroniseBookings(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/sync" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "html/sync.html")
	case "POST":
		sheetName := r.FormValue("sheet_name")
		roomNumber := r.FormValue("room_number")
		start := r.FormValue("start")
		end := r.FormValue("end")

		templ, err := template.ParseFiles("html/ExcelRes.html")
		if err != nil {
			zap.L().Error("SynchroniseBookings", zap.Error(err))
		}

		if roomNumber == "All" || roomNumber == "all" {
			bookings, err := h.ServiceExcel.GetAllBookingForPeriod(r.Context(), sheetName, start, end)
			if err != nil {
				fmt.Fprint(w, err)
				return
			}
			templ.Execute(w, bookings)
		} else {
			bookings, err := h.ServiceExcel.GetBookingForPeriod(r.Context(), sheetName, roomNumber, start, end)
			if err != nil {
				fmt.Fprint(w, err)
				return
			}
			templ.Execute(w, bookings)
		}

	}
}
