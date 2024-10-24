package services

import (
	"awesomeProject/internal/entities"
	"awesomeProject/internal/models"
	"context"
	"log"
)

//type Service struct {
//	storage *repo.Repository
//}
//
//func NewService(storage *repo.Repository) *Service {
//	return &Service{storage: storage}
//}

type Service struct {
	storage Storage
}

type Storage interface {
	Create(ctx context.Context, r entities.Reservation) (*entities.Reservation, error)
	CreateGuest(ctx context.Context, g entities.Guest) (*entities.Guest, error)
	FindGuestByPhoneNumber(ctx context.Context, phone string) (*entities.Guest, error)
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) CreateReservation(ctx context.Context, booking models.BookingInfo) error {
	g, err := booking.DbModelConvertGuest()
	if err != nil {
		return err
	}
	log.Println(*g)

	guest, err := s.storage.FindGuestByPhoneNumber(ctx, booking.Phone)
	if err != nil {
		return err
	}
	if guest == nil {
		guest, err = s.storage.CreateGuest(ctx, *g)
		if err != nil {
			return err
		}
	}

	reserv, err := booking.DbModelConvert(guest.GuestID)
	if err != nil {
		return err
	}
	log.Println(*reserv)
	_, err = s.storage.Create(ctx, *reserv)
	if err != nil {
		return err
	}
	return nil
}

// Создание репорта для собственика
func CreateReport() {

}

// Оповещение собственика о предстоящем бронирование
func FutureBooking() {

}

// свободные квартиры на эти даты
func FreeApartmentForDates() {

}

// цены расчитаные по переуду, к примеру низкий сезон переходящий в высокий
func AproksimatePrice() {

}

//func dbModelConvertGuest(booking excelCalendarScraper.BookingInfo) (*entities.Guest, error) {
//	id := uuid.New()
//	return &entities.Guest{
//		GuestID:     id,
//		Name:        booking.GuestName,
//		Phone:       booking.Phone,
//		Description: "",
//	}, nil
//}
//
//func dbModelConvert(booking excelCalendarScraper.BookingInfo, uuid uuid.UUID) (*entities.Reservation, error) {
//	timeFormat := "02.01.2006"
//	checkIn, err := time.Parse(timeFormat, booking.CheckIn)
//	if err != nil {
//		return nil, err
//	}
//	checkOut, err := time.Parse(timeFormat, booking.CheckOut)
//	if err != nil {
//		return nil, err
//	}
//
//	var price int
//	if booking.Price != "" {
//		price, err = strconv.Atoi(booking.Price)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	var cleaningPrice int
//	if booking.CleaningPrice != "" {
//		cleaningPrice, err = strconv.Atoi(booking.CleaningPrice)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	var adult int
//	if booking.Adult != "" {
//		adult, err = strconv.Atoi(booking.Adult)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	var children int
//	if booking.Children != "" {
//		children, err = strconv.Atoi(booking.Children)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return &entities.Reservation{
//		Oid:                        0,
//		RoomNumber:                 booking.RoomNumber,
//		GuestID:                    uuid,
//		CheckIn:                    checkIn,
//		CheckOut:                   checkOut,
//		Price:                      price,
//		CleaningPrice:              cleaningPrice,
//		ElectricityAndWaterPayment: booking.ElectricityAndWaterPayment,
//		Adult:                      adult,
//		Children:                   children,
//		Description:                booking.Description,
//	}, nil
//}
//
//func stringToInt(str ...string) ([]int, error) {
//	i := make([]int, 0)
//	for _, s := range str {
//		var v int
//		var err error
//		if s != "" {
//			v, err = strconv.Atoi(s)
//			if err != nil {
//				return nil, err
//			}
//
//		} else {
//			v = 0
//		}
//		i = append(i, v)
//	}
//	return i, nil
//}
