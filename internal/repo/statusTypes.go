package repo

import (
	"awesomeProject/internal/entities"
	"context"
)

// ListStatusTypes возвращает справочник активных типов статусов, отсортированных для UI.
func (db *Repository) ListStatusTypes(ctx context.Context) ([]entities.StatusType, error) {
	runner := getRunner(ctx, db.PostgreSQL)
	rows, err := runner.QueryContext(ctx,
		`SELECT id, code, name, color, sort_order, is_active
		 FROM status_types
		 WHERE is_active = TRUE
		 ORDER BY sort_order`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]entities.StatusType, 0)
	for rows.Next() {
		var st entities.StatusType
		if err := rows.Scan(&st.ID, &st.Code, &st.Name, &st.Color, &st.SortOrder, &st.IsActive); err != nil {
			return nil, err
		}
		result = append(result, st)
	}
	return result, rows.Err()
}

// CreateStatusType добавляет новый тип статуса в справочник.
func (db *Repository) CreateStatusType(ctx context.Context, st entities.StatusType) (*entities.StatusType, error) {
	runner := getRunner(ctx, db.PostgreSQL)
	result := new(entities.StatusType)

	query := `INSERT INTO status_types (code, name, color, sort_order, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, code, name, color, sort_order, is_active`

	row := runner.QueryRowContext(ctx, query, st.Code, st.Name, st.Color, st.SortOrder, st.IsActive)
	if err := row.Scan(&result.ID, &result.Code, &result.Name, &result.Color, &result.SortOrder, &result.IsActive); err != nil {
		return nil, err
	}
	return result, nil
}
