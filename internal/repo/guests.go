package repo

import (
	"awesomeProject/internal/entities"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (db *Repository) UpdateGuest(ctx context.Context, g entities.Guest) (*entities.Guest, error) {
	tx := From(ctx)
	if tx == nil {
		return nil, errors.New("context doesn't contain transaction")
	}

	guest := new(entities.Guest)
	queryContext := tx.QueryRowContext(ctx,
		"UPDATE Guests SET name=$2, phone=$3, description=$4 where guest_id=$1 RETURNING *",
		g.GuestID, g.Name, Nullable(g.Phone), g.Description)

	var sPhone sql.NullString
	err := queryContext.Scan(&guest.GuestID, &guest.Name, &sPhone, &guest.Description)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if sPhone.Valid {
		guest.Phone = sPhone.String
	} else {
		guest.Phone = ""
	}

	return guest, nil
}

//func (db *Repository) UpdateGuest(ctx context.Context, g entities.Guest) (*entities.Guest, error) {
//	guest := new(entities.Guest)
//	queryContext := db.PostgreSQL.QueryRowContext(ctx,
//		"UPDATE Guests SET name=$2, phone=$3, description=$4 where guest_id=$1 RETURNING *",
//		g.GuestID, g.Name, Nullable(g.Phone), g.Description)
//
//	var sPhone sql.NullString
//	err := queryContext.Scan(&guest.GuestID, &guest.Name, &sPhone, &guest.Description)
//	if errors.Is(err, sql.ErrNoRows) {
//		return nil, nil
//	}
//	if err != nil {
//		return nil, err
//	}
//
//	if sPhone.Valid {
//		guest.Phone = sPhone.String
//	} else {
//		guest.Phone = ""
//	}
//
//	return guest, nil
//}

type queryer interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func getRunner(ctx context.Context, db *sql.DB) queryer {
	if tx := From(ctx); tx != nil {
		return tx
	}
	zap.L().Debug("getRunner", zap.Error(errors.New("context doesn't contain transaction")))
	return db
}

func (db *Repository) CreateGuest(ctx context.Context, g entities.Guest) (*entities.Guest, error) {
	runner := getRunner(ctx, db.PostgreSQL) // todo такое использование контекста надо переделать или ввести повсеместно

	guest := new(entities.Guest)
	queryContext := runner.QueryRowContext(ctx,
		"INSERT INTO Guests (guest_id, name, phone,description) VALUES ($1, $2, $3, $4) RETURNING *",
		g.GuestID, g.Name, Nullable(g.Phone), g.Description)

	var sPhone sql.NullString
	err := queryContext.Scan(&guest.GuestID, &guest.Name, &sPhone, &guest.Description)
	if err != nil {
		return nil, err
	}

	if sPhone.Valid {
		guest.Phone = sPhone.String
	} else {
		guest.Phone = ""
	}

	return guest, nil
}

//func (db *Repository) CreateGuest(ctx context.Context, g entities.Guest) (*entities.Guest, error) {
//	guest := new(entities.Guest)
//	queryContext := db.PostgreSQL.QueryRowContext(ctx,
//		"INSERT INTO Guests (guest_id, name, phone,description) VALUES ($1, $2, $3, $4) RETURNING *",
//		g.GuestID, g.Name, Nullable(g.Phone), g.Description)
//
//	var sPhone sql.NullString
//	err := queryContext.Scan(&guest.GuestID, &guest.Name, &sPhone, &guest.Description)
//	if err != nil {
//		return nil, err
//	}
//
//	if sPhone.Valid {
//		guest.Phone = sPhone.String
//	} else {
//		guest.Phone = ""
//	}
//
//	return guest, nil
//}

func (db *Repository) ReadGuest(ctx context.Context, guestID uuid.UUID) (*entities.Guest, error) {
	guest := new(entities.Guest)
	queryContext := db.PostgreSQL.QueryRowContext(ctx,
		"Select * from Guests where guest_id=$1",
		guestID)

	var sPhone sql.NullString
	err := queryContext.Scan(&guest.GuestID, &guest.Name, &sPhone, &guest.Description)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if sPhone.Valid {
		guest.Phone = sPhone.String
	} else {
		guest.Phone = ""
	}

	return guest, nil
}

func Nullable(field string) interface{} {
	if len(field) == 0 {
		return sql.NullString{}
	}
	return field
}

func (db *Repository) FindGuestByPhoneNumber(ctx context.Context, phone string) (*entities.Guest, error) {
	guest := new(entities.Guest)
	queryContext := db.PostgreSQL.QueryRowContext(ctx,
		"Select * from Guests where phone=$1",
		phone)

	err := queryContext.Scan(&guest.GuestID, &guest.Name, &guest.Phone, &guest.Description)
	if errors.Is(err, sql.ErrNoRows) {
		zap.L().Debug("FindGuestByPhoneNumber", zap.Error(err))
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return guest, nil
}
