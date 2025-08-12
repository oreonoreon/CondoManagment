package entities

import "time"

type AirbnbData struct {
	AirbnbRoom
	AirbnbPrice
	RatingAirbnbRoom
	CoordinatesAirbnbRoom
	AdditionalInfoBnbPrice
}

func NewAirbnbData(room AirbnbRoom, price AirbnbPrice, rating RatingAirbnbRoom, coordinates CoordinatesAirbnbRoom, additionalInfo AdditionalInfoBnbPrice) AirbnbData {
	return AirbnbData{
		room,
		price,
		rating,
		coordinates,
		additionalInfo,
	}
}

type AirbnbRoom struct {
	RoomID             int64    `json:"roomID,string"`
	Badges             []string `json:"badges"`
	Name               string   `json:"name"`
	Title              string   `json:"title"`
	TypeOfEstate       string   `json:"type"`
	Kind               string   `json:"kind"`
	Category           string   `json:"category"`
	Images             []string `json:"images"`
	UnderstandableType string
	HostID             string
}

type Host struct {
	ID        string
	Name      string
	Trustable bool
}

type RatingAirbnbRoom struct {
	RoomID      int64   `json:"-"`
	Value       float32 `json:"value"`
	ReviewCount int     `json:"reviewCount"`
}

type AirbnbPrice struct {
	RoomID       int64   `json:"-"`
	Price        float32 `json:"price"`
	CheckIn      time.Time
	CheckOut     time.Time
	ScrapingDate time.Time
}

type CoordinatesAirbnbRoom struct {
	RoomID       int64 `json:"-"`
	Latitude     float64
	Longitude    float64
	LocationName string `json:"locationName"`
}

type AdditionalInfoBnbPrice struct {
	Qualifier      string
	CurrencySymbol string
}
