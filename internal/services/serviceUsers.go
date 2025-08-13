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
}

func NewServiceUsers(storageUsers StorageUsers) *ServiceUsers {
	return &ServiceUsers{
		storageUsers: storageUsers,
	}
}

func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(b), err
}

func (u *ServiceUsers) PrepareToCreateUser(username string, password string, phone string) (*entities.User, error) {
	passwordHash, err := HashPassword(password)
	if err != nil {
		zap.L().Error("PrepareToCreateUser", zap.Error(err))
		return nil, err
	}
	id := uuid.New()

	role := "manager"

	user := entities.User{
		id,
		username,
		passwordHash,
		phone,
		role,
	}
	return &user, nil
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
