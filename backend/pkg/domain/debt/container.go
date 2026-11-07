package debt

import (
	"context"
)

type Container struct {
	Repository
}

type Repository interface {
	CreateDebt(ctx context.Context, debt *Debt) (*Debt, error)
	GetAllDebts(ctx context.Context) ([]Debt, error)
	GetDebtsByUserID(ctx context.Context, userID int) ([]Debt, error)
	GetDebtByID(ctx context.Context, id int) (*Debt, error)
	UpdateDebt(ctx context.Context, debt *Debt) (*Debt, error)
	DeleteDebt(ctx context.Context, id int) error
}
