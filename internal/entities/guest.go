package entities

import "github.com/google/uuid"

type Guest struct {
	GuestID     uuid.UUID `json:"guest_id"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	Description string    `json:"guestDescription"`
}
