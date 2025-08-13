package entities

import (
	"github.com/google/uuid"
	"time"
)

type Reservation struct {
	Oid                        int       `db:"id" json:"id"`
	RoomNumber                 string    `db:"room_number" json:"roomNumber"`
	GuestID                    uuid.UUID `db:"guest_id" json:"guest_uuid"`
	CheckIn                    time.Time `db:"check_in" json:"check_in"`
	CheckOut                   time.Time `db:"check_out" json:"check_out"`
	Price                      int       `db:"price" json:"price"`
	CleaningPrice              int       `db:"cleaning_price" json:"cleaning_price"`
	ElectricityAndWaterPayment string    `db:"electricity_and_water_payment" json:"electricity_and_water_payment"`
	Adult                      int       `db:"adult" json:"adult"`
	Children                   int       `db:"children" json:"children"`
	Description                string    `db:"reservationDescription" json:"reservationDescription"`
	Days                       int       `db:"days" json:"days"`
	PriceForOneNight           int       `db:"price_for_night" json:"price_for_night"`
}
