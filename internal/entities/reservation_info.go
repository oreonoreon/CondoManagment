package entities

type ReservationInfo struct {
	ID               int    `db:"id"                 json:"id"`
	ReservationID    int    `db:"reservation_id"     json:"reservation_id"`
	Deposit          int    `db:"deposit"            json:"deposit"`
	DepositCurrency  string `db:"deposit_currency"   json:"deposit_currency"`
	Prepayment       int    `db:"prepayment"         json:"prepayment"`
	PaymentOnCheckin int    `db:"payment_on_checkin" json:"payment_on_checkin"`
	Notes            string `db:"notes"              json:"notes"`
}
