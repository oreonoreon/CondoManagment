package entities

type Apartment struct {
	Oid            int    `json:"oid"`
	RoomNumber     string `json:"room_number"`
	Description    string `json:"description"`
	AirbnbCalendar string `json:"airbnbCalendar"`
}
