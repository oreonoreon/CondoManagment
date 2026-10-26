package services

import (
	"awesomeProject/internal/entities"
	"awesomeProject/internal/erro"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeAdminStorage запоминает последние вызовы; методы, не нужные тестам, ничего не делают.
type fakeAdminStorage struct {
	roles       []entities.Role
	updatedRole *[]string
	deletedRole string
	deletedUser uuid.UUID
	updatedUser *entities.UserUpdate
	lastAccess  *entities.ApartmentAccess
}

func (f *fakeAdminStorage) ListUsers(context.Context) ([]entities.UserInfo, error) { return nil, nil }
func (f *fakeAdminStorage) UpdateUser(_ context.Context, id uuid.UUID, upd entities.UserUpdate) (*entities.UserInfo, error) {
	f.updatedUser = &upd
	return &entities.UserInfo{ID: id}, nil
}
func (f *fakeAdminStorage) DeleteUser(_ context.Context, id uuid.UUID) error {
	f.deletedUser = id
	return nil
}
func (f *fakeAdminStorage) ListApartmentsWithAccess(context.Context) ([]entities.ApartmentWithAccess, error) {
	return nil, nil
}
func (f *fakeAdminStorage) ReadApartmentWithAccess(context.Context, int) (*entities.ApartmentWithAccess, error) {
	return &entities.ApartmentWithAccess{}, nil
}
func (f *fakeAdminStorage) CreateApartment(_ context.Context, _ entities.Apartment, access entities.ApartmentAccess) (int, error) {
	f.lastAccess = &access
	return 1, nil
}
func (f *fakeAdminStorage) UpdateApartment(context.Context, int, entities.ApartmentUpdate) error {
	return nil
}
func (f *fakeAdminStorage) SetApartmentAccess(_ context.Context, _ int, access entities.ApartmentAccess) error {
	f.lastAccess = &access
	return nil
}
func (f *fakeAdminStorage) ListRoles(context.Context) ([]entities.Role, error) { return f.roles, nil }
func (f *fakeAdminStorage) CreateRole(context.Context, string, string) error   { return nil }
func (f *fakeAdminStorage) UpdateRole(_ context.Context, _ string, _ *string, permissions []string) error {
	f.updatedRole = &permissions
	return nil
}
func (f *fakeAdminStorage) DeleteRole(_ context.Context, name string) error {
	f.deletedRole = name
	return nil
}

func isValidation(err error) bool {
	var e *erro.Err
	return errors.As(err, &e) && e.Code == erro.CodeValidation
}

var adminUser = entities.CurrentUser{ID: uuid.New(), Role: "admin"}

func TestAdminCannotLockThemselvesOut(t *testing.T) {
	storage := &fakeAdminStorage{}
	s := NewServiceAdmin(storage, nil)
	ctx := context.Background()

	assert.True(t, isValidation(s.DeleteUser(ctx, adminUser, adminUser.ID)), "удалить себя")
	other := "manager"
	_, err := s.UpdateUser(ctx, adminUser, adminUser.ID, UserChanges{Role: &other})
	assert.True(t, isValidation(err), "сменить себе роль")
	assert.True(t, isValidation(s.DeleteRole(ctx, adminUser, "admin")), "удалить свою роль")

	noRolesManage := []string{entities.PermCalendarRead}
	assert.True(t, isValidation(s.UpdateRole(ctx, adminUser, "admin", RoleChanges{Permissions: &noRolesManage})),
		"снять roles.manage со своей роли")

	// с другими ролями и пользователями — можно
	require.NoError(t, s.UpdateRole(ctx, adminUser, "manager", RoleChanges{Permissions: &noRolesManage}))
	other2 := uuid.New()
	require.NoError(t, s.DeleteUser(ctx, adminUser, other2))
	assert.Equal(t, other2, storage.deletedUser)
}

func TestDeleteRoleWithUsers(t *testing.T) {
	storage := &fakeAdminStorage{roles: []entities.Role{{Name: "agent", UserCount: 2}, {Name: "empty"}}}
	s := NewServiceAdmin(storage, nil)

	assert.True(t, isValidation(s.DeleteRole(context.Background(), adminUser, "agent")))
	require.NoError(t, s.DeleteRole(context.Background(), adminUser, "empty"))
	assert.Equal(t, "empty", storage.deletedRole)
}

func TestRolePermissionsValidated(t *testing.T) {
	storage := &fakeAdminStorage{}
	s := NewServiceAdmin(storage, nil)
	ctx := context.Background()

	bad := []string{"calendar.read", "no.such.permission"}
	assert.True(t, isValidation(s.UpdateRole(ctx, adminUser, "agent", RoleChanges{Permissions: &bad})))

	dup := []string{"calendar.read", "booking.field.price", "calendar.read"}
	require.NoError(t, s.UpdateRole(ctx, adminUser, "agent", RoleChanges{Permissions: &dup}))
	assert.Equal(t, []string{"calendar.read", "booking.field.price"}, *storage.updatedRole)

	assert.True(t, isValidation(s.CreateRole(ctx, "1bad name", "", nil)), "имя роли")
	require.NoError(t, s.CreateRole(ctx, "cleaner", "Уборщица", []string{"cleaning.manage"}))
}

func TestApartmentAccessNormalized(t *testing.T) {
	storage := &fakeAdminStorage{}
	s := NewServiceAdmin(storage, nil)
	u := uuid.New()

	_, err := s.SetApartmentAccess(context.Background(), 1, entities.ApartmentAccess{
		Roles:   []string{"manager", "agent", "manager"},
		UserIDs: []uuid.UUID{u, u},
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"agent", "manager"}, storage.lastAccess.Roles)
	assert.Equal(t, []uuid.UUID{u}, storage.lastAccess.UserIDs)

	_, err = s.CreateApartment(context.Background(), entities.Apartment{RoomNumber: "  "}, entities.ApartmentAccess{})
	assert.True(t, isValidation(err), "пустой номер")
}

func TestUpdateUserValidation(t *testing.T) {
	storage := &fakeAdminStorage{}
	s := NewServiceAdmin(storage, nil)
	ctx := context.Background()
	id := uuid.New()

	short := "123"
	_, err := s.UpdateUser(ctx, adminUser, id, UserChanges{Password: &short})
	assert.True(t, isValidation(err), "короткий пароль")

	empty := "  "
	_, err = s.UpdateUser(ctx, adminUser, id, UserChanges{Name: &empty})
	assert.True(t, isValidation(err), "пустое имя")

	password, name := "secret123", "  Ivan "
	_, err = s.UpdateUser(ctx, adminUser, id, UserChanges{Password: &password, Name: &name})
	require.NoError(t, err)
	assert.Equal(t, "Ivan", *storage.updatedUser.Name)
	require.NotNil(t, storage.updatedUser.PasswordHash)
	assert.NotEqual(t, password, *storage.updatedUser.PasswordHash, "в БД уходит хеш, не пароль")
}
