package income

import (
	"context"
)

type Container struct {
	Repository
}

type Repository interface {
	CreateIncome(ctx context.Context, income *Income) (*Income, error)
	GetAllIncomes(ctx context.Context) ([]Income, error)
	GetIncomesByUserID(ctx context.Context, userID int) ([]Income, error)
	GetIncomeByID(ctx context.Context, id int) (*Income, error)
	UpdateIncome(ctx context.Context, income *Income) (*Income, error)
	DeleteIncome(ctx context.Context, id int) error
}
