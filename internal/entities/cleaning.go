package entities

import "time"

type Cleaning struct {
	ID            int       `db:"id"             json:"id"`
	ReservationID *int      `db:"reservation_id" json:"reservation_id,omitempty"`
	CleaningTime  time.Time `db:"cleaning_time"  json:"cleaning_time"`
	Room          string    `db:"room"           json:"room"`
	CleaningPrice int       `db:"cleaning_price" json:"cleaning_price"`
	LaundryPrice  int       `db:"laundry_price"  json:"laundry_price"`
	AgentName     string    `db:"agent_name"     json:"agent_name"`
	Description   string    `db:"description"    json:"description"`
	Paid          bool      `db:"paid"           json:"paid"`
}
