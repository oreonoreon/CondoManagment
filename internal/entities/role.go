package entities

import "github.com/google/uuid"

// Role — роль с её правами (для админки).
type Role struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
	UserCount   int      `json:"user_count"`
}

// UserInfo — пользователь без хеша пароля (для админки).
type UserInfo struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Phone string    `json:"phone"`
	Role  string    `json:"role"`
}

// UserUpdate — изменяемые поля пользователя; nil — не менять.
type UserUpdate struct {
	Name         *string
	Phone        *string
	Role         *string
	PasswordHash *string
}
