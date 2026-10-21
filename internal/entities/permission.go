package entities

import (
	"slices"

	"github.com/google/uuid"
)

// Права ролей (таблица role_permissions). Описание каждого — в migrations/014_roles_and_permissions.sql.
const (
	PermCalendarRead        = "calendar.read"
	PermCalendarWrite       = "calendar.write"
	PermTodayRead           = "today.read"
	PermCleaningManage      = "cleaning.manage"
	PermReportsRead         = "reports.read"
	PermUsersManage         = "users.manage"
	PermToolsUse            = "tools.use"
	PermApartmentsOwnedOnly = "apartments.owned_only"

	// PermBookingFieldPrefix + имя поля разрешает отдавать это поле брони.
	PermBookingFieldPrefix = "booking.field."
	PermBookingFieldAll    = PermBookingFieldPrefix + "all"
)

// Поля брони, видимость которых настраивается правами booking.field.<поле>.
const (
	BookingFieldRoomNumber       = "room_number"
	BookingFieldGuestID          = "guest_id"
	BookingFieldGuestName        = "guest_name"
	BookingFieldGuestPhone       = "guest_phone"
	BookingFieldGuestDescription = "guest_description"
	BookingFieldCheckIn          = "check_in"
	BookingFieldCheckOut         = "check_out"
	BookingFieldPrice            = "price"
	BookingFieldCleaningPrice    = "cleaning_price"
	BookingFieldElectricity      = "electricity_and_water_payment"
	BookingFieldAdult            = "adult"
	BookingFieldChildren         = "children"
	BookingFieldDescription      = "description"
	BookingFieldDays             = "days"
	BookingFieldPriceForNight    = "price_for_night"
	BookingFieldReservationInfo  = "reservation_info"
	BookingFieldStatuses         = "statuses"
)

// Permissions — набор прав роли.
type Permissions map[string]struct{}

func NewPermissions(names ...string) Permissions {
	p := make(Permissions, len(names))
	for _, name := range names {
		p[name] = struct{}{}
	}
	return p
}

func (p Permissions) Has(name string) bool {
	_, ok := p[name]
	return ok
}

// CanSeeAllBookingFields — роли отдаются все поля брони без обрезки.
func (p Permissions) CanSeeAllBookingFields() bool {
	return p.Has(PermBookingFieldAll)
}

// CanSeeBookingField — роли отдаётся поле брони field (см. BookingField*).
func (p Permissions) CanSeeBookingField(field string) bool {
	return p.CanSeeAllBookingFields() || p.Has(PermBookingFieldPrefix+field)
}

// List возвращает права отсортированным списком (для отдачи на фронт).
func (p Permissions) List() []string {
	list := make([]string, 0, len(p))
	for name := range p {
		list = append(list, name)
	}
	slices.Sort(list)
	return list
}

// CurrentUser — пользователь текущего запроса с актуальными ролью и правами из БД.
type CurrentUser struct {
	ID          uuid.UUID
	Name        string
	Role        string
	Permissions Permissions
}
