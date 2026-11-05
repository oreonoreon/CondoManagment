package entities

import "github.com/google/uuid"

type Apartment struct {
	Oid            int    `json:"oid"`
	RoomNumber     string `json:"room_number"`
	Description    string `json:"description"`
	AirbnbCalendar string `json:"airbnbCalendar"`
	IsActive       bool   `json:"is_active"`
	// Кондо и группа: nil — апартамент ещё никуда не отнесён
	CondoID   *int   `json:"condo_id"`
	CondoName string `json:"condo_name"`
	GroupID   *int   `json:"group_id"`
	GroupName string `json:"group_name"`
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

// ApartmentPlacement — куда отнести апартаменты: nil снимает кондо или группу.
type ApartmentPlacement struct {
	ApartmentIDs []int `json:"apartment_ids"`
	CondoID      *int  `json:"condo_id"`
	GroupID      *int  `json:"group_id"`
}

// Condo — кондо с группами (для админки).
type Condo struct {
	ID             int              `json:"id"`
	Name           string           `json:"name"`
	ApartmentCount int              `json:"apartment_count"`
	Groups         []ApartmentGroup `json:"groups"`
}

// ApartmentGroup — группа апартаментов внутри кондо.
type ApartmentGroup struct {
	ID             int    `json:"id"`
	CondoID        int    `json:"condo_id"`
	Name           string `json:"name"`
	ApartmentCount int    `json:"apartment_count"`
}

// ApartmentUpdate — изменяемые поля апартамента; nil — не менять.
type ApartmentUpdate struct {
	Description    *string `json:"description"`
	AirbnbCalendar *string `json:"airbnbCalendar"`
	IsActive       *bool   `json:"is_active"`
}
