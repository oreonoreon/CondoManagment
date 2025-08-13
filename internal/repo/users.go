package repo

import (
	"awesomeProject/internal/entities"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
)

func (db *Repository) CreateUser(ctx context.Context, user entities.User) (*entities.User, error) {
	u := new(entities.User)

	query := "INSERT INTO users (" +
		"id," +
		" name," +
		" passwordhash," +
		" phone," +
		" role" +
		") " +
		"VALUES ($1,$2,$3,$4,$5) Returning *"

	queryContext := db.PostgreSQL.QueryRowContext(
		ctx,
		query,
		user.ID,
		user.Name,
		user.PasswordHash,
		user.Phone,
		user.Role,
	)

	err := queryContext.Scan(
		&u.ID,
		&u.Name,
		&u.PasswordHash,
		&u.Phone,
		&u.Role,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (db *Repository) GetUser(ctx context.Context, username string) (*entities.User, error) {
	user := new(entities.User)
	queryContext := db.PostgreSQL.QueryRowContext(ctx,
		"Select * from users where name=$1",
		username)

	err := queryContext.Scan(&user.ID, &user.Name, &user.PasswordHash, &user.Phone, &user.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *Repository) UpdateUser(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return nil, nil
}

func (db *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return nil
}
