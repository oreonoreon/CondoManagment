package main

import (
	"awesomeProject/internal/calendarScraper"
	"awesomeProject/internal/excelCalendarScraper"
	"awesomeProject/internal/myLogger"
	"awesomeProject/internal/repo"
	"fmt"
	_ "modernc.org/sqlite"
	"net/http"
)

func main() {
	//create Db model
	//model, err := repo.DataBasePostgreSQl()
	//if err != nil {
	//	panic(err)
	//}
	//defer model.PostgreSQL.Close()

	//err = calendarScraper.DownloadCalendar(model, "A606")
	//if err != nil {
	//	myLogger.Logger.Println(err)
	//}
	excelCalendarScraper.ExcelCalendarScraper()
	//server(model)
}

func server(model repo.DBSql) {
	http.HandleFunc("/", Reservation(model))
	http.HandleFunc("/apartment", Apartments(model))
	http.HandleFunc("/sync", SynchroniseCalendar(model))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func SynchroniseCalendar(dbModel repo.DBSql) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sync" {
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		}

		switch r.Method {
		case "GET":
			http.ServeFile(w, r, "sync.html")
		case "POST":
			if r.FormValue("sync") == "Sync" {
				calendarScraper.ScrapAll(r.Context(), dbModel)
			}
		}
	}
}
func Apartments(dbModel repo.DBSql) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/apartment" {
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		}

		switch r.Method {
		case "GET":
			http.ServeFile(w, r, "apartment.html")
		case "POST":
			apart := repo.Apartment{
				RoomNumber:     r.FormValue("room_number"),
				Description:    r.FormValue("description"),
				AirbnbCalendar: r.FormValue("airbnb_calendar"),
			}
			err := dbModel.CreateApartment(r.Context(), apart)
			if err != nil {
				myLogger.Logger.Printf("ParseForm() err: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.ServeFile(w, r, "apartment.html")
		case "PUT":
			fallthrough
		case "DELETE":
			fallthrough
		default:
			fmt.Fprintf(w, "Sorry, only GET and POST and PUT and DELETE methods are supported.")
		}
	}
}

var Checkin string
var Checkout string

func Reservation(dbModel repo.DBSql) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		}

		switch r.Method {
		case "GET":
			http.ServeFile(w, r, "form.html")
		case "POST":
			// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}

			Checkin = r.FormValue("checkin")
			Checkout = r.FormValue("checkout")

			reservation, err := dbModel.Read(r.Context(), Checkin, Checkout)

			if err != nil {
				fmt.Fprintf(w, "Read: %v", err)
			}

			//_, err = dbModel.Create(r.Context(), reservation.RoomNumber, Checkin, Checkout)
			//
			//if err != nil {
			//	fmt.Fprintf(w, "Create: %v", err)
			//}

			//templ, err := template.ParseFiles("./response.html")
			//if err != nil {
			//	myLogger.Logger.Println(err)
			//}
			//templ.Execute(w, results)

			fmt.Fprintf(w, "%v", reservation)
		default:
			fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
		}
	}
}
