package entities

import (
	"github.com/google/uuid"
	"time"
)

type StatusType struct {
	ID        int    `db:"id" json:"id"`
	Code      string `db:"code" json:"code"`
	Name      string `db:"name" json:"name"`
	Color     string `db:"color" json:"color"`
	SortOrder int    `db:"sort_order" json:"sort_order"`
	IsActive  bool   `db:"is_active" json:"is_active"`
}

// ReservationStatus - фактическое присвоение статуса брони (с историей).
// IsActive=true - статус сейчас включен, false - был снят (запись остаётся для истории).
type ReservationStatus struct {
	ID            int        `db:"id" json:"id"`
	ReservationID int        `db:"reservation_id" json:"reservation_id"`
	StatusTypeID  int        `db:"status_type_id" json:"status_type_id"`
	IsActive      bool       `db:"is_active" json:"is_active"`
	SetBy         *uuid.UUID `db:"set_by" json:"set_by,omitempty"`
	SetAt         time.Time  `db:"set_at" json:"set_at"`
	RemovedAt     *time.Time `db:"removed_at" json:"removed_at,omitempty"`

	// Поля из status_types, подтягиваются JOIN'ом для удобства фронта
	StatusCode  string `db:"-" json:"status_code,omitempty"`
	StatusName  string `db:"-" json:"status_name,omitempty"`
	StatusColor string `db:"-" json:"status_color,omitempty"`
}
