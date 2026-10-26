package entities

import (
	"slices"

	"github.com/google/uuid"
)

// Права ролей (таблица role_permissions). Описание каждого — в PermissionCatalog.
const (
	PermCalendarRead       = "calendar.read"
	PermCalendarWrite      = "calendar.write"
	PermTodayRead          = "today.read"
	PermCleaningManage     = "cleaning.manage"
	PermReportsRead        = "reports.read"
	PermUsersManage        = "users.manage"
	PermToolsUse           = "tools.use"
	PermApartmentsAll      = "apartments.all"
	PermApartmentsInactive = "apartments.inactive"
	// PermApartmentsAssignedOnly — роль видит только апартаменты, выданные пользователю лично;
	// выдать такой роли апартамент целиком нельзя (например, владельцы — каждому свои).
	PermApartmentsAssignedOnly = "apartments.assigned_only"
	PermApartmentsManage       = "apartments.manage"
	PermRolesManage        = "roles.manage"

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

// ApartmentScope — какие апартаменты видны пользователю.
type ApartmentScope struct {
	Role     string
	UserID   uuid.UUID
	All          bool // apartments.all — все апартаменты без выдачи доступа
	Inactive     bool // apartments.inactive — включая неактивные
	AssignedOnly bool // apartments.assigned_only — только выданные лично, выдача роли не действует
}

func (u CurrentUser) ApartmentScope() ApartmentScope {
	return ApartmentScope{
		Role:         u.Role,
		UserID:       u.ID,
		All:          u.Permissions.Has(PermApartmentsAll),
		Inactive:     u.Permissions.Has(PermApartmentsInactive),
		AssignedOnly: u.Permissions.Has(PermApartmentsAssignedOnly),
	}
}

// PermissionInfo — право и его описание для админки.
type PermissionInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Group       string `json:"group"`
}

// PermissionCatalog — все права, которые понимает бэкенд. Админка выдаёт ролям только их.
func PermissionCatalog() []PermissionInfo {
	catalog := []PermissionInfo{
		{PermCalendarRead, "Просмотр календаря бронирований", "Календарь"},
		{PermCalendarWrite, "Создание, изменение и удаление броней, статусов, депозитов", "Календарь"},
		{PermTodayRead, "Вкладка «Сегодня»: заезды и выезды по дате", "Календарь"},
		{PermCleaningManage, "Уборки (по любым комнатам)", "Календарь"},
		{PermReportsRead, "Отчёты", "Календарь"},
		{PermToolsUse, "Excel-синхронизация и парсер Airbnb", "Инструменты"},
		{PermApartmentsAll, "Все апартаменты без выдачи доступа", "Апартаменты"},
		{PermApartmentsInactive, "Видеть неактивные апартаменты", "Апартаменты"},
		{PermApartmentsAssignedOnly, "Только апартаменты, выданные пользователю лично (роли целиком не выдаются)", "Апартаменты"},
		{PermUsersManage, "Админка: пользователи", "Админка"},
		{PermApartmentsManage, "Админка: апартаменты и доступ к ним", "Админка"},
		{PermRolesManage, "Админка: роли и права", "Админка"},
		{PermBookingFieldAll, "Все поля брони", "Поля брони"},
	}
	for _, f := range bookingFieldInfo {
		catalog = append(catalog, PermissionInfo{PermBookingFieldPrefix + f.field, f.description, "Поля брони"})
	}
	return catalog
}

var bookingFieldInfo = []struct{ field, description string }{
	{BookingFieldRoomNumber, "Номер апартамента"},
	{BookingFieldGuestName, "Имя гостя"},
	{BookingFieldGuestPhone, "Телефон гостя"},
	{BookingFieldGuestDescription, "Описание гостя"},
	{BookingFieldGuestID, "ID гостя"},
	{BookingFieldCheckIn, "Заезд"},
	{BookingFieldCheckOut, "Выезд"},
	{BookingFieldDays, "Количество ночей"},
	{BookingFieldAdult, "Взрослые"},
	{BookingFieldChildren, "Дети"},
	{BookingFieldPrice, "Цена"},
	{BookingFieldPriceForNight, "Цена за ночь"},
	{BookingFieldCleaningPrice, "Цена уборки"},
	{BookingFieldElectricity, "Электричество и вода"},
	{BookingFieldDescription, "Описание брони"},
	{BookingFieldReservationInfo, "Депозит, предоплата, фактические заезд и выезд"},
	{BookingFieldStatuses, "Статусы брони"},
}

// IsKnownPermission — право есть в PermissionCatalog.
func IsKnownPermission(name string) bool {
	for _, p := range PermissionCatalog() {
		if p.Name == name {
			return true
		}
	}
	return false
}
