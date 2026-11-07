package debt

import (
	"context"
	"errors"
	"time"
)

var (
	ErrDebtNotFound    = errors.New("debt not found")
	ErrInvalidDebtData = errors.New("invalid debt data")
	ErrDatabaseError   = errors.New("database error")
)

type Service interface {
	CreateDebt(ctx context.Context, request CreateDebtRequest) (*Debt, error)
	GetAllDebts(ctx context.Context) ([]Debt, error)
	GetDebtsByUserID(ctx context.Context, userID int) ([]Debt, error)
	GetDebtByID(ctx context.Context, id int) (*Debt, error)
	UpdateDebt(ctx context.Context, id int, request UpdateDebtRequest) (*Debt, error)
	DeleteDebt(ctx context.Context, id int) error
}

type service struct {
	*Container
}

func New(container *Container) Service {
	return &service{
		Container: container,
	}
}

func (s *service) CreateDebt(ctx context.Context, request CreateDebtRequest) (*Debt, error) {
	dueDate, err := time.Parse("2006-01-02", request.DueDate)
	if err != nil {
		return nil, ErrInvalidDebtData
	}

	debt := &Debt{
		UserID:            request.UserID,
		Name:              request.Name,
		TotalAmount:       request.TotalAmount,
		RemainingAmount:   request.RemainingAmount,
		DueDate:           dueDate,
		InterestRate:      request.InterestRate,
		NumInstallments:   request.NumInstallments,
		InstallmentAmount: request.InstallmentAmount,
		PaymentDay:        request.PaymentDay,
		Paid:              false,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	createdDebt, err := s.Repository.CreateDebt(ctx, debt)
	if err != nil {
		return nil, err
	}

	return createdDebt, nil
}

func (s *service) GetAllDebts(ctx context.Context) ([]Debt, error) {
	debts, err := s.Repository.GetAllDebts(ctx)
	if err != nil {
		return nil, err
	}

	return debts, nil
}

func (s *service) GetDebtsByUserID(ctx context.Context, userID int) ([]Debt, error) {
	debts, err := s.Repository.GetDebtsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return debts, nil
}

func (s *service) GetDebtByID(ctx context.Context, id int) (*Debt, error) {
	debt, err := s.Repository.GetDebtByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return debt, nil
}

func (s *service) UpdateDebt(ctx context.Context, id int, request UpdateDebtRequest) (*Debt, error) {
	existingDebt, err := s.Repository.GetDebtByID(ctx, id)
	if err != nil {
		return nil, err
	}

	dueDate, err := time.Parse("2006-01-02", request.DueDate)
	if err != nil {
		return nil, ErrInvalidDebtData
	}

	existingDebt.Name = request.Name
	existingDebt.TotalAmount = request.TotalAmount
	existingDebt.RemainingAmount = request.RemainingAmount
	existingDebt.DueDate = dueDate
	existingDebt.InterestRate = request.InterestRate
	existingDebt.NumInstallments = request.NumInstallments
	existingDebt.InstallmentAmount = request.InstallmentAmount
	existingDebt.PaymentDay = request.PaymentDay
	existingDebt.Paid = request.Paid
	existingDebt.UpdatedAt = time.Now()

	updatedDebt, err := s.Repository.UpdateDebt(ctx, existingDebt)
	if err != nil {
		return nil, err
	}

	return updatedDebt, nil
}

func (s *service) DeleteDebt(ctx context.Context, id int) error {
	_, err := s.Repository.GetDebtByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.Repository.DeleteDebt(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
