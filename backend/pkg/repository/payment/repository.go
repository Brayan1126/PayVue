package payment

import (
	"context"
	"database/sql"
	"errors"

	"github.com/payvue/payvue-backend/pkg/domain/payment"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) payment.Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreatePayment(ctx context.Context, p *payment.Payment) (*payment.Payment, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, payment.ErrDatabaseError
	}
	defer tx.Rollback()

	// Insertar el pago
	query := `
		INSERT INTO payments (user_id, amount, debt_id, receipt_filename, date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := tx.ExecContext(ctx, query,
		p.UserID, p.Amount, p.DebtID, p.ReceiptFilename, p.Date, p.CreatedAt, p.UpdatedAt,
	)

	if err != nil {
		return nil, payment.ErrDatabaseError
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, payment.ErrDatabaseError
	}

	p.ID = int(id)

	// Actualizar la deuda
	updateQuery := `
		UPDATE debts 
		SET remaining_amount = CASE 
			WHEN remaining_amount - ? < 0 THEN 0 
			ELSE remaining_amount - ? 
		END,
		paid = CASE 
			WHEN remaining_amount - ? <= 0 THEN 1 
			ELSE paid 
		END,
		updated_at = ?
		WHERE id = ?
	`

	_, err = tx.ExecContext(ctx, updateQuery, p.Amount, p.Amount, p.Amount, p.UpdatedAt, p.DebtID)
	if err != nil {
		return nil, payment.ErrDatabaseError
	}

	if err := tx.Commit(); err != nil {
		return nil, payment.ErrDatabaseError
	}

	return p, nil
}

func (r *repository) GetAllPayments(ctx context.Context) ([]payment.PaymentWithDebt, error) {
	query := `
		SELECT 
			p.id, COALESCE(p.user_id, 0), p.amount, p.debt_id, p.receipt_filename, p.date, p.created_at, p.updated_at,
			d.name, d.remaining_amount, d.installment_amount
		FROM payments p
		INNER JOIN debts d ON p.debt_id = d.id
		ORDER BY p.date DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, payment.ErrDatabaseError
	}
	defer rows.Close()

	var payments []payment.PaymentWithDebt
	for rows.Next() {
		var pwd payment.PaymentWithDebt
		err := rows.Scan(
			&pwd.ID, &pwd.UserID, &pwd.Amount, &pwd.DebtID, &pwd.ReceiptFilename, &pwd.Date,
			&pwd.CreatedAt, &pwd.UpdatedAt,
			&pwd.DebtName, &pwd.DebtRemainingAmount, &pwd.DebtInstallmentAmount,
		)
		if err != nil {
			return nil, payment.ErrDatabaseError
		}
		payments = append(payments, pwd)
	}

	if err = rows.Err(); err != nil {
		return nil, payment.ErrDatabaseError
	}

	if payments == nil {
		payments = []payment.PaymentWithDebt{}
	}

	return payments, nil
}

func (r *repository) GetPaymentsByUserID(ctx context.Context, userID int) ([]payment.PaymentWithDebt, error) {
	query := `
		SELECT 
			p.id, COALESCE(p.user_id, 0), p.amount, p.debt_id, p.receipt_filename, p.date, p.created_at, p.updated_at,
			d.name, d.remaining_amount, d.installment_amount
		FROM payments p
		INNER JOIN debts d ON p.debt_id = d.id
		WHERE p.user_id = ?
		ORDER BY p.date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, payment.ErrDatabaseError
	}
	defer rows.Close()

	var payments []payment.PaymentWithDebt
	for rows.Next() {
		var pwd payment.PaymentWithDebt
		err := rows.Scan(
			&pwd.ID, &pwd.UserID, &pwd.Amount, &pwd.DebtID, &pwd.ReceiptFilename, &pwd.Date,
			&pwd.CreatedAt, &pwd.UpdatedAt,
			&pwd.DebtName, &pwd.DebtRemainingAmount, &pwd.DebtInstallmentAmount,
		)
		if err != nil {
			return nil, payment.ErrDatabaseError
		}
		payments = append(payments, pwd)
	}

	if err = rows.Err(); err != nil {
		return nil, payment.ErrDatabaseError
	}

	if payments == nil {
		payments = []payment.PaymentWithDebt{}
	}

	return payments, nil
}

func (r *repository) GetPaymentByID(ctx context.Context, id int) (*payment.Payment, error) {
	query := `
		SELECT id, COALESCE(user_id, 0), amount, debt_id, receipt_filename, date, created_at, updated_at
		FROM payments
		WHERE id = ?
	`

	var p payment.Payment
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.UserID, &p.Amount, &p.DebtID, &p.ReceiptFilename, &p.Date, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, payment.ErrPaymentNotFound
		}
		return nil, payment.ErrDatabaseError
	}

	return &p, nil
}

func (r *repository) DeletePayment(ctx context.Context, id int) error {
	query := `DELETE FROM payments WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return payment.ErrDatabaseError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return payment.ErrDatabaseError
	}

	if rowsAffected == 0 {
		return payment.ErrPaymentNotFound
	}

	return nil
}
