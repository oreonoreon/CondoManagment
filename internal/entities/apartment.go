package entities

type Apartment struct {
	Oid            int
	RoomNumber     string `json:"room_number"`
	Description    string
	AirbnbCalendar string
}
