package entities

import (
	"github.com/google/uuid"
	"time"
)

type Reservation struct {
	Oid                        int
	RoomNumber                 string    `db:"room_number"`
	GuestID                    uuid.UUID `db:"guest_id"`
	CheckIn                    time.Time `db:"check_in"`
	CheckOut                   time.Time `db:"check_out"`
	Price                      int       `db:"price"`
	CleaningPrice              int       `db:"cleaning_price"`
	ElectricityAndWaterPayment string    `db:"electricity_and_water_payment"`
	Adult                      int       `db:"adult"`
	Children                   int       `db:"children"`
	Description                string    `db:"description"`
	Days                       int       `db:"days"`
	PriceForOneNight           int       `db:"price_for_night"`
}
