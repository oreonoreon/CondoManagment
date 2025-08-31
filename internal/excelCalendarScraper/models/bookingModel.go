package models

import (
	"awesomeProject/internal/entities"
	"awesomeProject/internal/erro"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type BookingInfo struct {
	RoomNumber                 string
	CheckIn                    string
	CheckOut                   string
	GuestName                  string
	Phone                      string
	Price                      string
	CleaningPrice              string
	ElectricityAndWaterPayment string
	Adult                      string
	Children                   string
	Description                string
}

func (b *BookingInfo) ParseOutBookingInfo(s []string) {
	for k, v := range s {

		//проверим на наличие следующего индекса в слайсе
		if len(s) < k+1+1 {
			return
		}

		switch v {
		case "Name":
			b.GuestName = s[k+1]
		case "Check in":
			b.CheckIn = s[k+1]
		case "Check out":
			b.CheckOut = s[k+1]
		case "Price":
			b.Price = s[k+1]
		case "Cleaning price":
			b.CleaningPrice = s[k+1]
		case "Electricity and water payment":
			b.ElectricityAndWaterPayment = s[k+1]
		case "Adult":
			b.Adult = s[k+1]
		case "children":
			b.Children = s[k+1]
		case "Phone":
			b.Phone = s[k+1]
		case "Description":
			b.Description = s[k+1]
		}
	}
}

func (b *BookingInfo) Validate() error {
	err := b.HaveName()
	if err != nil {
		return err
	}
	err = b.HaveName()
	if err != nil {
		return err
	}
	return nil
}

func (b *BookingInfo) HaveName() error {
	if b.GuestName == "" {
		return erro.ErrBookingExcelModelHaveNotName
	}
	return nil
}

func (b *BookingInfo) HavePhone() error {
	if b.Phone == "+" || b.Phone == "@" {
		return erro.ErrBookingExcelModelHaveNotPhone
	}
	return nil
}

func (b BookingInfo) DbModelConvertGuest() (*entities.Guest, error) {
	id := uuid.New()
	return &entities.Guest{
		GuestID:     id,
		Name:        b.GuestName,
		Phone:       b.Phone,
		Description: "",
	}, nil
}

func TimeConvert(date string) (time.Time, error) {
	timeFormat := "02.01.2006"
	t, err := time.Parse(timeFormat, date)
	if err != nil {
		zap.L().Debug("TimeConvert", zap.Error(err))
		//return time.Time{}, err
	}
	timeFormat = "2006-01-02"
	t, err = time.Parse(timeFormat, date)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (b BookingInfo) DbModelConvert(uuid uuid.UUID) (*entities.Reservation, error) {
	timeFormat := "02.01.2006"
	checkIn, err := time.Parse(timeFormat, b.CheckIn)
	if err != nil {
		return nil, err
	}
	checkOut, err := time.Parse(timeFormat, b.CheckOut)
	if err != nil {
		return nil, err
	}

	var price int
	if b.Price != "" {
		price, err = strconv.Atoi(b.Price)
		if err != nil {
			return nil, err
		}
	}

	var cleaningPrice int
	if b.CleaningPrice != "" {
		cleaningPrice, err = strconv.Atoi(b.CleaningPrice)
		if err != nil {
			return nil, err
		}
	}

	var adult int
	if b.Adult != "" {
		adult, err = strconv.Atoi(b.Adult)
		if err != nil {
			return nil, err
		}
	}

	var children int
	if b.Children != "" {
		children, err = strconv.Atoi(b.Children)
		if err != nil {
			return nil, err
		}
	}

	days := int(checkOut.Sub(checkIn).Hours() / 24)
	priceForOneNight := int(price / days)

	return &entities.Reservation{
		Oid:                        0,
		RoomNumber:                 b.RoomNumber,
		GuestID:                    uuid,
		CheckIn:                    checkIn,
		CheckOut:                   checkOut,
		Price:                      price,
		CleaningPrice:              cleaningPrice,
		ElectricityAndWaterPayment: b.ElectricityAndWaterPayment,
		Adult:                      adult,
		Children:                   children,
		Description:                b.Description,
		Days:                       days,
		PriceForOneNight:           priceForOneNight,
	}, nil
}
