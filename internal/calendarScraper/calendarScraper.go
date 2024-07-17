package calendarScraper

import (
	"awesomeProject/internal/myLogger"
	"awesomeProject/internal/repo"
	"context"
	"fmt"
	ics "github.com/arran4/golang-ical"
	"io"
	"net/http"
	"os"
	"time"
)

func ScrapAll(ctx context.Context, dbModel repo.DBSql) {
	apartments, err := dbModel.ReadApartmentAll(ctx)
	if err != nil {
		myLogger.Logger.Printf("ScrapAll err: %v", err)
		return
	}
	for _, apartment := range apartments {
		Scrap(ctx, dbModel, apartment)
	}
}
func Scrap(ctx context.Context, dbModel repo.DBSql, apartment repo.Apartment) {
	res, err := requestCalendar(apartment)
	if err != nil {
		myLogger.Logger.Printf("Scrap err: %v", err)
		return
	}
	defer res.Body.Close()

	events, err := parseCalendar(res.Body)
	if err != nil {
		myLogger.Logger.Printf("Scrap err: %v", err)
		return
	}

	for _, event := range events {
		start, _ := event.GetStartAt()
		end, _ := event.GetEndAt()

		//s := start.Format("2006-01-02")
		//e := end.Format("2006-01-02")

		guest := repo.Guest{
			Description: event.GetProperty(ics.ComponentProperty(ics.PropertySummary)).Value,
		}

		reservation := repo.Reservaton{
			RoomNumber: apartment.RoomNumber,
			CheckIn:    start,
			CheckOut:   end,
		}
		writeEventInDB(ctx, dbModel, reservation, guest)
		//fmt.Println(event.GetStartAt())
		//fmt.Println(event.GetEndAt())
		//fmt.Println(event.GetProperty(ics.ComponentProperty(ics.PropertySummary)).Value)
	}

}

func writeEventInDB(ctx context.Context, dbModel repo.DBSql, r repo.Reservaton, guest repo.Guest) {
	reservation, err := dbModel.ReadWithRoomNumber(ctx, r.RoomNumber, r.CheckIn.Format("2006-01-02"), r.CheckOut.Format("2006-01-02"))
	if err != nil {
		myLogger.Logger.Printf("writeEventInDB err: %v", err)
		return
	}
	if len(reservation) == 0 {
		g, err := dbModel.CreateGuest(ctx, guest)
		if err != nil {
			myLogger.Logger.Printf("writeEventInDB: %v\nReservation: %v\nGuest: %v", err, r, guest)
			return
		}

		r.GuestID = g.GuestID

		_, err = dbModel.Create(ctx, r)
		if err != nil {
			myLogger.Logger.Printf("writeEventInDB: %v", err)
			return
		}
	} else if len(reservation) == 1 {
		//todo придумать что нибудь что бы можно было определить что за гость за бронил даты
		//g, err := dbModel.ReadGuest(ctx, reservation[0].GuestID)
		//if err != nil {
		//	myLogger.Logger.Printf("writeEventInDB: %v", err)
		//	return
		//}

		//апдейтим дату чекаута так при закрытие дат в календаре аэрбнб в плотную, бронирование не отличимо
		err = dbModel.Update(ctx, reservation[0].Oid, reservation[0].CheckIn.Format("2006-01-02"), r.CheckOut.Format("2006-01-02"))
		if err != nil {
			myLogger.Logger.Printf("writeEventInDB: %v", err)
			return
		}

	} else {
		for _, v := range reservation {
			myLogger.Logger.Println("OverBooking apartment: ", v)
		}
		return
	}
}

func requestCalendar(apartment repo.Apartment) (*http.Response, error) {
	res, err := http.Get(apartment.AirbnbCalendar)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func parseCalendar(f io.ReadCloser) ([]*ics.VEvent, error) {
	cal, err := ics.ParseCalendar(f)
	if err != nil {
		return nil, err
	}

	return cal.Events(), nil
}

func DownloadCalendar(dbModel repo.DBSql, roomNumber string) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer cancel()
	apartment, err := dbModel.ReadApartment(ctx, roomNumber)
	if err != nil {
		return fmt.Errorf("downloadCalendar: %w", err)
	}

	res, err := requestCalendar(*apartment)
	if err != nil {
		return fmt.Errorf("downloadCalendar: %w", err)
	}
	defer res.Body.Close()

	out, err := os.Create(fmt.Sprintf("./%v.ics", roomNumber))
	if err != nil {
		return fmt.Errorf("downloadCalendar: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return fmt.Errorf("downloadCalendar: %w", err)
	}
	return nil
}
