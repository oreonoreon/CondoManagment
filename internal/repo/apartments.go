package repo

import (
	"awesomeProject/internal/entities"
	"context"
	"database/sql"
)

// apartmentInScope — условие видимости апартамента a для ApartmentScope (аргументы — apartmentScopeArgs):
// активный (или есть право на неактивные) и доступ выдан пользователю лично или его роли
// (выдача роли не учитывается при apartments.assigned_only), либо есть право на все.
const apartmentInScope = `($3 OR a.is_active) AND ($2 OR EXISTS (
	SELECT 1 FROM apartment_access aa
	WHERE aa.room_number = a.room_number AND ((aa.role = $1 AND NOT $5) OR aa.user_id = $4)))`

func apartmentScopeArgs(scope entities.ApartmentScope) []any {
	return []any{scope.Role, scope.All, scope.Inactive, scope.UserID, scope.AssignedOnly}
}

const apartmentColumns = `a.id, a.room_number, a.description, a.airbnb_calendar, a.is_active`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanApartment(row rowScanner) (entities.Apartment, error) {
	var apartment entities.Apartment
	var roomNumber, description, airbnbCalendar sql.NullString
	err := row.Scan(&apartment.Oid, &roomNumber, &description, &airbnbCalendar, &apartment.IsActive)
	apartment.RoomNumber = roomNumber.String
	apartment.Description = description.String
	apartment.AirbnbCalendar = airbnbCalendar.String
	return apartment, err
}

// ReadApartmentsFor возвращает апартаменты, видимые пользователю.
func (db *Repository) ReadApartmentsFor(ctx context.Context, scope entities.ApartmentScope) ([]entities.Apartment, error) {
	rows, err := db.PostgreSQL.QueryContext(ctx,
		`SELECT `+apartmentColumns+` FROM apartments a WHERE `+apartmentInScope+` ORDER BY a.room_number`,
		apartmentScopeArgs(scope)...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	apartments := make([]entities.Apartment, 0)
	for rows.Next() {
		apartment, err := scanApartment(rows)
		if err != nil {
			return nil, err
		}
		apartments = append(apartments, apartment)
	}
	return apartments, rows.Err()
}
