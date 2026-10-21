package services

import (
	"awesomeProject/internal/entities"
	"awesomeProject/internal/erro"
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type ServiceUsers struct {
	storageUsers StorageUsers
}

type StorageUsers interface {
	CreateUser(ctx context.Context, user entities.User) (*entities.User, error)
	GetUser(ctx context.Context, username string) (*entities.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID) (*entities.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ReadCurrentUser(ctx context.Context, id uuid.UUID) (*entities.CurrentUser, error)
	RoleExists(ctx context.Context, role string) (bool, error)
}

// defaultUserRole — роль нового пользователя, если она не указана явно.
const defaultUserRole = "manager"

func NewServiceUsers(storageUsers StorageUsers) *ServiceUsers {
	return &ServiceUsers{
		storageUsers: storageUsers,
	}
}

func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(b), err
}

// PrepareToCreateUser готовит пользователя к записи. Пустая роль заменяется на defaultUserRole,
// роль должна существовать в справочнике roles.
func (u *ServiceUsers) PrepareToCreateUser(ctx context.Context, username string, password string, phone string, role string) (*entities.User, error) {
	if role == "" {
		role = defaultUserRole
	}
	exists, err := u.storageUsers.RoleExists(ctx, role)
	if err != nil {
		zap.L().Error("PrepareToCreateUser: RoleExists", zap.Error(err))
		return nil, err
	}
	if !exists {
		return nil, erro.ErrUnknownRole
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		zap.L().Error("PrepareToCreateUser", zap.Error(err))
		return nil, err
	}

	user := entities.User{
		ID:           uuid.New(),
		Name:         username,
		PasswordHash: passwordHash,
		Phone:        phone,
		Role:         role,
	}
	return &user, nil
}

// GetCurrentUser возвращает пользователя с актуальными ролью и правами. nil, если пользователя нет.
func (u *ServiceUsers) GetCurrentUser(ctx context.Context, id uuid.UUID) (*entities.CurrentUser, error) {
	return u.storageUsers.ReadCurrentUser(ctx, id)
}

func (u *ServiceUsers) GetUser(ctx context.Context, username string) (*entities.User, error) {
	user, err := u.storageUsers.GetUser(ctx, username)
	if err != nil {
		zap.L().Error("GetUser", zap.Error(err))
	}
	if user == nil {
		return nil, erro.ErrWrongCreds
	}
	return user, nil
}

func (u *ServiceUsers) CreateUser(ctx context.Context, user entities.User) (*entities.User, error) {
	createdUser, err := u.storageUsers.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}
