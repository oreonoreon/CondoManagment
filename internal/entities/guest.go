package entities

import "github.com/google/uuid"

type Guest struct {
	GuestID     uuid.UUID
	Name        string
	Phone       string
	Description string
}
