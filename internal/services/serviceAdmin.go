package services

import (
	"awesomeProject/internal/entities"
	"awesomeProject/internal/erro"
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/google/uuid"
)

// ServiceAdmin — админка: пользователи, апартаменты с доступом, роли и права.
type ServiceAdmin struct {
	storage StorageAdmin
	users   *ServiceUsers
}

type StorageAdmin interface {
	ListUsers(ctx context.Context) ([]entities.UserInfo, error)
	UpdateUser(ctx context.Context, id uuid.UUID, upd entities.UserUpdate) (*entities.UserInfo, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error

	ListApartmentsWithAccess(ctx context.Context) ([]entities.ApartmentWithAccess, error)
	ReadApartmentWithAccess(ctx context.Context, id int) (*entities.ApartmentWithAccess, error)
	CreateApartment(ctx context.Context, apartment entities.Apartment, access entities.ApartmentAccess) (int, error)
	UpdateApartment(ctx context.Context, id int, upd entities.ApartmentUpdate) error
	SetApartmentAccess(ctx context.Context, id int, access entities.ApartmentAccess) error

	ListRoles(ctx context.Context) ([]entities.Role, error)
	CreateRole(ctx context.Context, name, description string) error
	UpdateRole(ctx context.Context, name string, description *string, permissions []string) error
	DeleteRole(ctx context.Context, name string) error
}

func NewServiceAdmin(storage StorageAdmin, users *ServiceUsers) *ServiceAdmin {
	return &ServiceAdmin{storage: storage, users: users}
}

// ─── Пользователи ───

// UserChanges — изменения пользователя из админки; nil — не менять.
type UserChanges struct {
	Name     *string `json:"name"`
	Phone    *string `json:"phone"`
	Role     *string `json:"role"`
	Password *string `json:"password"`
}

func (s *ServiceAdmin) ListUsers(ctx context.Context) ([]entities.UserInfo, error) {
	return s.storage.ListUsers(ctx)
}

func (s *ServiceAdmin) CreateUser(ctx context.Context, name, password, phone, role string) (*entities.UserInfo, error) {
	user, err := s.users.PrepareToCreateUser(ctx, name, password, phone, role)
	if err != nil {
		return nil, err
	}
	created, err := s.users.CreateUser(ctx, *user)
	if err != nil {
		return nil, err
	}
	return &entities.UserInfo{ID: created.ID, Name: created.Name, Phone: created.Phone, Role: created.Role}, nil
}

func (s *ServiceAdmin) UpdateUser(ctx context.Context, current entities.CurrentUser, id uuid.UUID, ch UserChanges) (*entities.UserInfo, error) {
	upd := entities.UserUpdate{}
	if ch.Name != nil {
		name := strings.TrimSpace(*ch.Name)
		if err := ValidateUserName(name); err != nil {
			return nil, err
		}
		upd.Name = &name
	}
	if ch.Phone != nil {
		phone := strings.TrimSpace(*ch.Phone)
		if err := ValidatePhone(phone); err != nil {
			return nil, err
		}
		upd.Phone = &phone
	}
	if ch.Role != nil {
		if id == current.ID && *ch.Role != current.Role {
			return nil, erro.NewValidation("Нельзя сменить роль самому себе")
		}
		upd.Role = ch.Role
	}
	if ch.Password != nil {
		if err := ValidatePassword(*ch.Password); err != nil {
			return nil, err
		}
		hash, err := HashPassword(*ch.Password)
		if err != nil {
			return nil, err
		}
		upd.PasswordHash = &hash
	}
	return s.storage.UpdateUser(ctx, id, upd)
}

func (s *ServiceAdmin) DeleteUser(ctx context.Context, current entities.CurrentUser, id uuid.UUID) error {
	if id == current.ID {
		return erro.NewValidation("Нельзя удалить самого себя")
	}
	return s.storage.DeleteUser(ctx, id)
}

// ─── Апартаменты ───

func (s *ServiceAdmin) ListApartments(ctx context.Context) ([]entities.ApartmentWithAccess, error) {
	return s.storage.ListApartmentsWithAccess(ctx)
}

func (s *ServiceAdmin) CreateApartment(ctx context.Context, apartment entities.Apartment, access entities.ApartmentAccess) (*entities.ApartmentWithAccess, error) {
	apartment.RoomNumber = strings.TrimSpace(apartment.RoomNumber)
	if apartment.RoomNumber == "" {
		return nil, erro.NewValidation("Укажите номер апартамента")
	}
	id, err := s.storage.CreateApartment(ctx, apartment, normalizeAccess(access))
	if err != nil {
		return nil, err
	}
	return s.storage.ReadApartmentWithAccess(ctx, id)
}

func (s *ServiceAdmin) UpdateApartment(ctx context.Context, id int, ch entities.ApartmentUpdate) (*entities.ApartmentWithAccess, error) {
	err := s.storage.UpdateApartment(ctx, id, ch)
	if err != nil {
		return nil, err
	}
	return s.storage.ReadApartmentWithAccess(ctx, id)
}

func (s *ServiceAdmin) SetApartmentAccess(ctx context.Context, id int, access entities.ApartmentAccess) (*entities.ApartmentWithAccess, error) {
	if err := s.storage.SetApartmentAccess(ctx, id, normalizeAccess(access)); err != nil {
		return nil, err
	}
	return s.storage.ReadApartmentWithAccess(ctx, id)
}

// normalizeAccess убирает повторы и nil-срезы.
func normalizeAccess(access entities.ApartmentAccess) entities.ApartmentAccess {
	roles := slices.Compact(slices.Sorted(slices.Values(access.Roles)))
	if roles == nil {
		roles = []string{}
	}
	userIDs := make([]uuid.UUID, 0, len(access.UserIDs))
	for _, id := range access.UserIDs {
		if !slices.Contains(userIDs, id) {
			userIDs = append(userIDs, id)
		}
	}
	return entities.ApartmentAccess{Roles: roles, UserIDs: userIDs}
}

// ─── Роли ───

// roleNamePattern — имя роли: латиница, цифры и _, с буквы (как admin, apartmentOwner).
var roleNamePattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]{1,49}$`)

// RoleChanges — изменения роли из админки; nil — не менять (permissions заменяются целиком).
type RoleChanges struct {
	Description *string   `json:"description"`
	Permissions *[]string `json:"permissions"`
}

func (s *ServiceAdmin) ListRoles(ctx context.Context) ([]entities.Role, error) {
	return s.storage.ListRoles(ctx)
}

func (s *ServiceAdmin) CreateRole(ctx context.Context, name, description string, permissions []string) error {
	if !roleNamePattern.MatchString(name) {
		return erro.NewValidation("Имя роли: латиница, цифры и _, начинается с буквы, от 2 до 50 символов")
	}
	perms, err := validatePermissions(permissions)
	if err != nil {
		return err
	}
	if err := s.storage.CreateRole(ctx, name, strings.TrimSpace(description)); err != nil {
		return err
	}
	return s.storage.UpdateRole(ctx, name, nil, perms)
}

func (s *ServiceAdmin) UpdateRole(ctx context.Context, current entities.CurrentUser, name string, ch RoleChanges) error {
	var perms []string
	if ch.Permissions != nil {
		var err error
		perms, err = validatePermissions(*ch.Permissions)
		if err != nil {
			return err
		}
		// иначе администратор сам отрежет себе доступ к админке ролей
		if name == current.Role && !slices.Contains(perms, entities.PermRolesManage) {
			return erro.NewValidation("Нельзя снять право roles.manage со своей роли")
		}
	}
	var description *string
	if ch.Description != nil {
		d := strings.TrimSpace(*ch.Description)
		description = &d
	}
	return s.storage.UpdateRole(ctx, name, description, perms)
}

func (s *ServiceAdmin) DeleteRole(ctx context.Context, current entities.CurrentUser, name string) error {
	if name == current.Role {
		return erro.NewValidation("Нельзя удалить свою роль")
	}
	roles, err := s.storage.ListRoles(ctx)
	if err != nil {
		return err
	}
	for _, role := range roles {
		if role.Name == name && role.UserCount > 0 {
			return erro.NewValidation(fmt.Sprintf("У роли есть пользователи (%d) — сначала смените им роль", role.UserCount))
		}
	}
	return s.storage.DeleteRole(ctx, name)
}

// validatePermissions проверяет права по каталогу и убирает повторы. Пустой список — роль без прав.
func validatePermissions(permissions []string) ([]string, error) {
	result := make([]string, 0, len(permissions))
	for _, p := range permissions {
		if !entities.IsKnownPermission(p) {
			return nil, erro.NewValidation("Неизвестное право: " + p)
		}
		if !slices.Contains(result, p) {
			result = append(result, p)
		}
	}
	return result, nil
}

func (s *ServiceAdmin) PermissionCatalog() []entities.PermissionInfo {
	return entities.PermissionCatalog()
}
