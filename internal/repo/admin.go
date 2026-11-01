package repo

import (
	"awesomeProject/internal/entities"
	"awesomeProject/internal/erro"
	"context"
	"database/sql"
	"errors"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// translateAdminPQ переводит ошибки ограничений Postgres в доменные ошибки админки.
func translateAdminPQ(err error) error {
	var pqe *pq.Error
	if errors.As(err, &pqe) {
		switch pqe.Code {
		case "23505": // unique_violation
			return erro.ErrAlreadyExists
		case "23503": // foreign_key_violation
			return erro.ErrInvalidReference
		}
	}
	return err
}

// withTx выполняет fn в отдельной транзакции.
func (db *Repository) withTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := db.PostgreSQL.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := fn(tx); err != nil {
		return translateAdminPQ(err)
	}
	return translateAdminPQ(tx.Commit())
}

// setClause собирает "col1=$1, col2=$2" из непустых полей обновления.
type setClause struct {
	parts []string
	args  []any
}

func (s *setClause) add(column string, value any) {
	s.args = append(s.args, value)
	s.parts = append(s.parts, column+"=$"+strconv.Itoa(len(s.args)))
}

// ─── Пользователи ───

func (db *Repository) ListUsers(ctx context.Context) ([]entities.UserInfo, error) {
	rows, err := db.PostgreSQL.QueryContext(ctx, "SELECT id, name, phone, role FROM users ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]entities.UserInfo, 0)
	for rows.Next() {
		var u entities.UserInfo
		if err := rows.Scan(&u.ID, &u.Name, &u.Phone, &u.Role); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// UpdateUser меняет пользователя. Возвращает erro.ErrNotFound, если его нет.
func (db *Repository) UpdateUser(ctx context.Context, id uuid.UUID, upd entities.UserUpdate) (*entities.UserInfo, error) {
	var set setClause
	if upd.Name != nil {
		set.add("name", *upd.Name)
	}
	if upd.Phone != nil {
		set.add("phone", *upd.Phone)
	}
	if upd.Role != nil {
		set.add("role", *upd.Role)
	}
	if upd.PasswordHash != nil {
		set.add("passwordhash", *upd.PasswordHash)
	}

	query := "SELECT id, name, phone, role FROM users WHERE id=$1"
	args := []any{id}
	if len(set.parts) > 0 {
		args = append(set.args, id)
		query = "UPDATE users SET " + strings.Join(set.parts, ", ") +
			" WHERE id=$" + strconv.Itoa(len(args)) + " RETURNING id, name, phone, role"
	}

	var u entities.UserInfo
	err := db.PostgreSQL.QueryRowContext(ctx, query, args...).Scan(&u.ID, &u.Name, &u.Phone, &u.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, erro.ErrNotFound
	}
	if err != nil {
		return nil, translateAdminPQ(err)
	}
	return &u, nil
}

// DeleteUser удаляет пользователя. Возвращает erro.ErrNotFound, если его нет.
func (db *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	res, err := db.PostgreSQL.ExecContext(ctx, "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return translateAdminPQ(err)
	}
	return requireAffected(res)
}

func requireAffected(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return erro.ErrNotFound
	}
	return nil
}

// ─── Апартаменты ───

// ListApartmentsWithAccess возвращает все апартаменты (включая неактивные) с выданным доступом.
func (db *Repository) ListApartmentsWithAccess(ctx context.Context) ([]entities.ApartmentWithAccess, error) {
	return db.readApartmentsWithAccess(ctx, nil)
}

// ReadApartmentWithAccess возвращает апартамент с выданным доступом или erro.ErrNotFound.
func (db *Repository) ReadApartmentWithAccess(ctx context.Context, id int) (*entities.ApartmentWithAccess, error) {
	list, err := db.readApartmentsWithAccess(ctx, &id)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, erro.ErrNotFound
	}
	return &list[0], nil
}

func (db *Repository) readApartmentsWithAccess(ctx context.Context, id *int) ([]entities.ApartmentWithAccess, error) {
	rows, err := db.PostgreSQL.QueryContext(ctx,
		`SELECT `+apartmentColumns+` FROM apartments a WHERE $1::int IS NULL OR a.id = $1 ORDER BY a.room_number`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]entities.ApartmentWithAccess, 0)
	index := make(map[string]int)
	for rows.Next() {
		apartment, err := scanApartment(rows)
		if err != nil {
			return nil, err
		}
		index[apartment.RoomNumber] = len(result)
		result = append(result, entities.ApartmentWithAccess{
			Apartment: apartment,
			Roles:     []string{},
			Users:     []entities.ApartmentAccessUser{},
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	grants, err := db.PostgreSQL.QueryContext(ctx,
		`SELECT aa.room_number, aa.role, u.id, u.name, u.role
		FROM apartment_access aa
		LEFT JOIN users u ON u.id = aa.user_id
		ORDER BY aa.role, u.name`)
	if err != nil {
		return nil, err
	}
	defer grants.Close()

	for grants.Next() {
		var room string
		var role, userName, userRole sql.NullString
		var userID uuid.NullUUID
		if err := grants.Scan(&room, &role, &userID, &userName, &userRole); err != nil {
			return nil, err
		}
		i, ok := index[room]
		if !ok {
			continue
		}
		if role.Valid {
			result[i].Roles = append(result[i].Roles, role.String)
		} else if userID.Valid {
			result[i].Users = append(result[i].Users, entities.ApartmentAccessUser{ID: userID.UUID, Name: userName.String, Role: userRole.String})
		}
	}
	return result, grants.Err()
}

// CreateApartment создаёт апартамент вместе с выданным доступом и возвращает его id.
func (db *Repository) CreateApartment(ctx context.Context, apartment entities.Apartment, access entities.ApartmentAccess) (int, error) {
	var id int
	err := db.withTx(ctx, func(tx *sql.Tx) error {
		err := tx.QueryRowContext(ctx,
			`INSERT INTO apartments (room_number, description, airbnb_calendar, is_active) VALUES ($1, $2, $3, $4) RETURNING id`,
			apartment.RoomNumber, apartment.Description, apartment.AirbnbCalendar, apartment.IsActive,
		).Scan(&id)
		if err != nil {
			return err
		}
		return insertApartmentAccess(ctx, tx, apartment.RoomNumber, access)
	})
	return id, err
}

// UpdateApartment меняет апартамент. Возвращает erro.ErrNotFound, если его нет.
func (db *Repository) UpdateApartment(ctx context.Context, id int, upd entities.ApartmentUpdate) error {
	var set setClause
	if upd.Description != nil {
		set.add("description", *upd.Description)
	}
	if upd.AirbnbCalendar != nil {
		set.add("airbnb_calendar", *upd.AirbnbCalendar)
	}
	if upd.IsActive != nil {
		set.add("is_active", *upd.IsActive)
	}
	if len(set.parts) == 0 {
		var exists bool
		err := db.PostgreSQL.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM apartments WHERE id=$1)", id).Scan(&exists)
		if err == nil && !exists {
			err = erro.ErrNotFound
		}
		return err
	}

	args := append(set.args, id)
	res, err := db.PostgreSQL.ExecContext(ctx,
		"UPDATE apartments SET "+strings.Join(set.parts, ", ")+" WHERE id=$"+strconv.Itoa(len(args)), args...)
	if err != nil {
		return translateAdminPQ(err)
	}
	return requireAffected(res)
}

// SetApartmentAccess заменяет весь выданный к апартаменту доступ. erro.ErrNotFound, если апартамента нет.
func (db *Repository) SetApartmentAccess(ctx context.Context, id int, access entities.ApartmentAccess) error {
	return db.withTx(ctx, func(tx *sql.Tx) error {
		var room string
		err := tx.QueryRowContext(ctx, "SELECT room_number FROM apartments WHERE id=$1 FOR UPDATE", id).Scan(&room)
		if errors.Is(err, sql.ErrNoRows) {
			return erro.ErrNotFound
		}
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, "DELETE FROM apartment_access WHERE room_number=$1", room); err != nil {
			return err
		}
		return insertApartmentAccess(ctx, tx, room, access)
	})
}

func insertApartmentAccess(ctx context.Context, tx *sql.Tx, room string, access entities.ApartmentAccess) error {
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO apartment_access (room_number, role) SELECT $1, unnest($2::varchar[]) ON CONFLICT DO NOTHING`,
		room, pq.Array(access.Roles)); err != nil {
		return err
	}
	userIDs := make([]string, 0, len(access.UserIDs))
	for _, id := range access.UserIDs {
		userIDs = append(userIDs, id.String())
	}
	_, err := tx.ExecContext(ctx,
		`INSERT INTO apartment_access (room_number, user_id) SELECT $1, unnest($2::uuid[]) ON CONFLICT DO NOTHING`,
		room, pq.Array(userIDs))
	return err
}

// ─── Роли ───

func (db *Repository) ListRoles(ctx context.Context) ([]entities.Role, error) {
	rows, err := db.PostgreSQL.QueryContext(ctx,
		`SELECT r.name, r.description,
			COALESCE((SELECT array_agg(rp.permission ORDER BY rp.permission) FROM role_permissions rp WHERE rp.role = r.name), '{}'),
			(SELECT count(*) FROM users u WHERE u.role = r.name)
		FROM roles r
		ORDER BY r.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]entities.Role, 0)
	for rows.Next() {
		var role entities.Role
		if err := rows.Scan(&role.Name, &role.Description, pq.Array(&role.Permissions), &role.UserCount); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (db *Repository) CreateRole(ctx context.Context, name, description string) error {
	_, err := db.PostgreSQL.ExecContext(ctx, "INSERT INTO roles (name, description) VALUES ($1, $2)", name, description)
	return translateAdminPQ(err)
}

// UpdateRole меняет описание (nil — не менять) и права (nil — не менять, иначе заменяет целиком).
// Если среди прав есть apartments.assigned_only, выдача апартаментов этой роли удаляется.
func (db *Repository) UpdateRole(ctx context.Context, name string, description *string, permissions []string) error {
	return db.withTx(ctx, func(tx *sql.Tx) error {
		var exists bool
		if err := tx.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM roles WHERE name=$1)", name).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return erro.ErrNotFound
		}
		if description != nil {
			if _, err := tx.ExecContext(ctx, "UPDATE roles SET description=$2 WHERE name=$1", name, *description); err != nil {
				return err
			}
		}
		if permissions != nil {
			if _, err := tx.ExecContext(ctx, "DELETE FROM role_permissions WHERE role=$1", name); err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx,
				"INSERT INTO role_permissions (role, permission) SELECT $1, unnest($2::varchar[]) ON CONFLICT DO NOTHING",
				name, pq.Array(permissions)); err != nil {
				return err
			}
			// выдача апартаментов роли с apartments.assigned_only не действует — убираем её, чтобы не путала
			if slices.Contains(permissions, entities.PermApartmentsAssignedOnly) {
				if _, err := tx.ExecContext(ctx, "DELETE FROM apartment_access WHERE role=$1", name); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// DeleteRole удаляет роль (её права и выданный ей доступ удаляются каскадом).
// erro.ErrNotFound, если роли нет; erro.ErrInvalidReference, если у роли есть пользователи.
func (db *Repository) DeleteRole(ctx context.Context, name string) error {
	res, err := db.PostgreSQL.ExecContext(ctx, "DELETE FROM roles WHERE name=$1", name)
	if err != nil {
		return translateAdminPQ(err)
	}
	return requireAffected(res)
}
