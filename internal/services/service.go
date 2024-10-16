package services

import (
	"awesomeProject/internal/excelCalendarScraper"
	"awesomeProject/internal/repo"
	"context"
	"github.com/google/uuid"
	"log"
	"strconv"
	"time"
)

type Service struct {
	Storage *repo.Repository
}

func NewService(storage *repo.Repository) *Service {
	return &Service{Storage: storage}
}

func (s *Service) CreateReservation(ctx context.Context, booking excelCalendarScraper.BookingInfo) error {
	g, err := dbModelConvertGuest(booking)
	if err != nil {
		return err
	}
	log.Println(*g)

	guest, err := s.Storage.FindGuestByPhoneNumber(ctx, booking.Phone)
	if err != nil {
		return err
	}
	if guest == nil {
		guest, err = s.Storage.CreateGuest(ctx, *g)
		if err != nil {
			return err
		}
	}

	reserv, err := dbModelConvert(booking, guest.GuestID)
	if err != nil {
		return err
	}
	log.Println(*reserv)
	_, err = s.Storage.Create(ctx, *reserv)
	if err != nil {
		return err
	}
	return nil
}

func dbModelConvertGuest(booking excelCalendarScraper.BookingInfo) (*repo.Guest, error) {
	id := uuid.New()
	return &repo.Guest{
		GuestID:     id,
		Name:        booking.GuestName,
		Phone:       booking.Phone,
		Description: "",
	}, nil
}

func dbModelConvert(booking excelCalendarScraper.BookingInfo, uuid uuid.UUID) (*repo.Reservation, error) {
	timeFormat := "02.01.2006"
	checkIn, err := time.Parse(timeFormat, booking.CheckIn)
	if err != nil {
		return nil, err
	}
	checkOut, err := time.Parse(timeFormat, booking.CheckOut)
	if err != nil {
		return nil, err
	}

	var price int
	if booking.Price != "" {
		price, err = strconv.Atoi(booking.Price)
		if err != nil {
			return nil, err
		}
	}

	var cleaningPrice int
	if booking.CleaningPrice != "" {
		cleaningPrice, err = strconv.Atoi(booking.CleaningPrice)
		if err != nil {
			return nil, err
		}
	}

	var adult int
	if booking.Adult != "" {
		adult, err = strconv.Atoi(booking.Adult)
		if err != nil {
			return nil, err
		}
	}

	var children int
	if booking.Children != "" {
		children, err = strconv.Atoi(booking.Children)
		if err != nil {
			return nil, err
		}
	}

	return &repo.Reservation{
		Oid:                        0,
		RoomNumber:                 booking.RoomNumber,
		GuestID:                    uuid,
		CheckIn:                    checkIn,
		CheckOut:                   checkOut,
		Price:                      price,
		CleaningPrice:              cleaningPrice,
		ElectricityAndWaterPayment: booking.ElectricityAndWaterPayment,
		Adult:                      adult,
		Children:                   children,
		Description:                booking.Description,
	}, nil
}

func stringToInt(str ...string) ([]int, error) {
	i := make([]int, 0)
	for _, s := range str {
		var v int
		var err error
		if s != "" {
			v, err = strconv.Atoi(s)
			if err != nil {
				return nil, err
			}

		} else {
			v = 0
		}
		i = append(i, v)
	}
	return i, nil
}
