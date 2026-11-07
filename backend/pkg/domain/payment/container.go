package payment

import (
	"context"
)

type Container struct {
	Repository
}

type Repository interface {
	CreatePayment(ctx context.Context, payment *Payment) (*Payment, error)
	GetAllPayments(ctx context.Context) ([]PaymentWithDebt, error)
	GetPaymentsByUserID(ctx context.Context, userID int) ([]PaymentWithDebt, error)
	GetPaymentByID(ctx context.Context, id int) (*Payment, error)
	DeletePayment(ctx context.Context, id int) error
}

type PaymentWithDebt struct {
	Payment
	DebtName              string
	DebtRemainingAmount   float64
	DebtInstallmentAmount float64
}
