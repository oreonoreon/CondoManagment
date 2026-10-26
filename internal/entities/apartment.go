package entities

import "github.com/google/uuid"

type Apartment struct {
	Oid            int    `json:"oid"`
	RoomNumber     string `json:"room_number"`
	Description    string `json:"description"`
	AirbnbCalendar string `json:"airbnbCalendar"`
	IsActive       bool   `json:"is_active"`
}

// ApartmentAccessUser — пользователь, которому лично выдан доступ к апартаменту.
type ApartmentAccessUser struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Role string    `json:"role"`
}

// ApartmentWithAccess — апартамент и выданный к нему доступ (для админки).
type ApartmentWithAccess struct {
	Apartment
	Roles []string              `json:"roles"`
	Users []ApartmentAccessUser `json:"users"`
}

// ApartmentAccess — кому выдан доступ к апартаменту.
type ApartmentAccess struct {
	Roles   []string    `json:"roles"`
	UserIDs []uuid.UUID `json:"user_ids"`
}

// ApartmentUpdate — изменяемые поля апартамента; nil — не менять.
type ApartmentUpdate struct {
	Description    *string `json:"description"`
	AirbnbCalendar *string `json:"airbnbCalendar"`
	IsActive       *bool   `json:"is_active"`
}
