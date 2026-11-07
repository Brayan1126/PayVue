package income

import (
	"context"
	"database/sql"
	"errors"

	"github.com/payvue/payvue-backend/pkg/domain/income"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) income.Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateIncome(ctx context.Context, i *income.Income) (*income.Income, error) {
	query := `
		INSERT INTO incomes (user_id, amount, source, date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		i.UserID, i.Amount, i.Source, i.Date, i.CreatedAt, i.UpdatedAt,
	)

	if err != nil {
		return nil, income.ErrDatabaseError
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, income.ErrDatabaseError
	}

	i.ID = int(id)
	return i, nil
}

func (r *repository) GetAllIncomes(ctx context.Context) ([]income.Income, error) {
	query := `
		SELECT id, COALESCE(user_id, 0), amount, source, date, created_at, updated_at
		FROM incomes
		ORDER BY date DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, income.ErrDatabaseError
	}
	defer rows.Close()

	var incomes []income.Income
	for rows.Next() {
		var i income.Income
		err := rows.Scan(
			&i.ID, &i.UserID, &i.Amount, &i.Source, &i.Date, &i.CreatedAt, &i.UpdatedAt,
		)
		if err != nil {
			return nil, income.ErrDatabaseError
		}
		incomes = append(incomes, i)
	}

	if err = rows.Err(); err != nil {
		return nil, income.ErrDatabaseError
	}

	if incomes == nil {
		incomes = []income.Income{}
	}

	return incomes, nil
}

func (r *repository) GetIncomesByUserID(ctx context.Context, userID int) ([]income.Income, error) {
	query := `
		SELECT id, COALESCE(user_id, 0), amount, source, date, created_at, updated_at
		FROM incomes
		WHERE user_id = ?
		ORDER BY date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, income.ErrDatabaseError
	}
	defer rows.Close()

	var incomes []income.Income
	for rows.Next() {
		var i income.Income
		err := rows.Scan(
			&i.ID, &i.UserID, &i.Amount, &i.Source, &i.Date, &i.CreatedAt, &i.UpdatedAt,
		)
		if err != nil {
			return nil, income.ErrDatabaseError
		}
		incomes = append(incomes, i)
	}

	if err = rows.Err(); err != nil {
		return nil, income.ErrDatabaseError
	}

	if incomes == nil {
		incomes = []income.Income{}
	}

	return incomes, nil
}

func (r *repository) GetIncomeByID(ctx context.Context, id int) (*income.Income, error) {
	query := `
		SELECT id, COALESCE(user_id, 0), amount, source, date, created_at, updated_at
		FROM incomes
		WHERE id = ?
	`

	var i income.Income
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&i.ID, &i.UserID, &i.Amount, &i.Source, &i.Date, &i.CreatedAt, &i.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, income.ErrIncomeNotFound
		}
		return nil, income.ErrDatabaseError
	}

	return &i, nil
}

func (r *repository) UpdateIncome(ctx context.Context, i *income.Income) (*income.Income, error) {
	query := `
		UPDATE incomes 
		SET amount = ?, source = ?, date = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		i.Amount, i.Source, i.Date, i.UpdatedAt, i.ID,
	)

	if err != nil {
		return nil, income.ErrDatabaseError
	}

	return i, nil
}

func (r *repository) DeleteIncome(ctx context.Context, id int) error {
	query := `DELETE FROM incomes WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return income.ErrDatabaseError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return income.ErrDatabaseError
	}

	if rowsAffected == 0 {
		return income.ErrIncomeNotFound
	}

	return nil
}
