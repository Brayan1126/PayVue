package debt

import (
	"context"
	"database/sql"
	"errors"

	"github.com/payvue/payvue-backend/pkg/domain/debt"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) debt.Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateDebt(ctx context.Context, d *debt.Debt) (*debt.Debt, error) {
	query := `
		INSERT INTO debts (user_id, name, total_amount, remaining_amount, due_date, interest_rate, 
		                   num_installments, installment_amount, payment_day, paid, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		d.UserID, d.Name, d.TotalAmount, d.RemainingAmount, d.DueDate,
		d.InterestRate, d.NumInstallments, d.InstallmentAmount,
		d.PaymentDay, d.Paid, d.CreatedAt, d.UpdatedAt,
	)

	if err != nil {
		return nil, debt.ErrDatabaseError
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, debt.ErrDatabaseError
	}

	d.ID = int(id)
	return d, nil
}

func (r *repository) GetAllDebts(ctx context.Context) ([]debt.Debt, error) {
	query := `
		SELECT id, COALESCE(user_id, 0), name, total_amount, remaining_amount, due_date, interest_rate,
		       num_installments, installment_amount, payment_day, paid, created_at, updated_at
		FROM debts
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, debt.ErrDatabaseError
	}
	defer rows.Close()

	var debts []debt.Debt
	for rows.Next() {
		var d debt.Debt
		err := rows.Scan(
			&d.ID, &d.UserID, &d.Name, &d.TotalAmount, &d.RemainingAmount, &d.DueDate,
			&d.InterestRate, &d.NumInstallments, &d.InstallmentAmount,
			&d.PaymentDay, &d.Paid, &d.CreatedAt, &d.UpdatedAt,
		)
		if err != nil {
			return nil, debt.ErrDatabaseError
		}
		debts = append(debts, d)
	}

	if err = rows.Err(); err != nil {
		return nil, debt.ErrDatabaseError
	}

	if debts == nil {
		debts = []debt.Debt{}
	}

	return debts, nil
}

func (r *repository) GetDebtsByUserID(ctx context.Context, userID int) ([]debt.Debt, error) {
	query := `
		SELECT id, COALESCE(user_id, 0), name, total_amount, remaining_amount, due_date, interest_rate,
		       num_installments, installment_amount, payment_day, paid, created_at, updated_at
		FROM debts
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, debt.ErrDatabaseError
	}
	defer rows.Close()

	var debts []debt.Debt
	for rows.Next() {
		var d debt.Debt
		err := rows.Scan(
			&d.ID, &d.UserID, &d.Name, &d.TotalAmount, &d.RemainingAmount, &d.DueDate,
			&d.InterestRate, &d.NumInstallments, &d.InstallmentAmount,
			&d.PaymentDay, &d.Paid, &d.CreatedAt, &d.UpdatedAt,
		)
		if err != nil {
			return nil, debt.ErrDatabaseError
		}
		debts = append(debts, d)
	}

	if err = rows.Err(); err != nil {
		return nil, debt.ErrDatabaseError
	}

	if debts == nil {
		debts = []debt.Debt{}
	}

	return debts, nil
}

func (r *repository) GetDebtByID(ctx context.Context, id int) (*debt.Debt, error) {
	query := `
		SELECT id, COALESCE(user_id, 0), name, total_amount, remaining_amount, due_date, interest_rate,
		       num_installments, installment_amount, payment_day, paid, created_at, updated_at
		FROM debts
		WHERE id = ?
	`

	var d debt.Debt
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&d.ID, &d.UserID, &d.Name, &d.TotalAmount, &d.RemainingAmount, &d.DueDate,
		&d.InterestRate, &d.NumInstallments, &d.InstallmentAmount,
		&d.PaymentDay, &d.Paid, &d.CreatedAt, &d.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, debt.ErrDebtNotFound
		}
		return nil, debt.ErrDatabaseError
	}

	return &d, nil
}

func (r *repository) UpdateDebt(ctx context.Context, d *debt.Debt) (*debt.Debt, error) {
	query := `
		UPDATE debts 
		SET name = ?, total_amount = ?, remaining_amount = ?, due_date = ?,
		    interest_rate = ?, num_installments = ?, installment_amount = ?,
		    payment_day = ?, paid = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		d.Name, d.TotalAmount, d.RemainingAmount, d.DueDate,
		d.InterestRate, d.NumInstallments, d.InstallmentAmount,
		d.PaymentDay, d.Paid, d.UpdatedAt, d.ID,
	)

	if err != nil {
		return nil, debt.ErrDatabaseError
	}

	return d, nil
}

func (r *repository) DeleteDebt(ctx context.Context, id int) error {
	query := `DELETE FROM debts WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return debt.ErrDatabaseError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return debt.ErrDatabaseError
	}

	if rowsAffected == 0 {
		return debt.ErrDebtNotFound
	}

	return nil
}
