package entities

type Booking struct {
	Guest
	Reservation
	ReservationInfo ReservationInfo     `json:"reservation_info"`
	Statuses        []ReservationStatus `json:"statuses,omitempty"`
}
