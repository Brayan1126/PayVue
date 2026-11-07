package income

import (
	"context"
	"errors"
	"time"
)

var (
	ErrIncomeNotFound    = errors.New("income not found")
	ErrInvalidIncomeData = errors.New("invalid income data")
	ErrDatabaseError     = errors.New("database error")
)

type Service interface {
	CreateIncome(ctx context.Context, request CreateIncomeRequest) (*Income, error)
	GetAllIncomes(ctx context.Context) ([]Income, error)
	GetIncomesByUserID(ctx context.Context, userID int) ([]Income, error)
	GetIncomeByID(ctx context.Context, id int) (*Income, error)
	UpdateIncome(ctx context.Context, id int, request UpdateIncomeRequest) (*Income, error)
	DeleteIncome(ctx context.Context, id int) error
}

type service struct {
	*Container
}

func New(container *Container) Service {
	return &service{
		Container: container,
	}
}

func (s *service) CreateIncome(ctx context.Context, request CreateIncomeRequest) (*Income, error) {
	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		return nil, ErrInvalidIncomeData
	}

	income := &Income{
		UserID:    request.UserID,
		Amount:    request.Amount,
		Source:    request.Source,
		Date:      date,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdIncome, err := s.Repository.CreateIncome(ctx, income)
	if err != nil {
		return nil, err
	}

	return createdIncome, nil
}

func (s *service) GetAllIncomes(ctx context.Context) ([]Income, error) {
	incomes, err := s.Repository.GetAllIncomes(ctx)
	if err != nil {
		return nil, err
	}

	return incomes, nil
}

func (s *service) GetIncomesByUserID(ctx context.Context, userID int) ([]Income, error) {
	incomes, err := s.Repository.GetIncomesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return incomes, nil
}

func (s *service) GetIncomeByID(ctx context.Context, id int) (*Income, error) {
	income, err := s.Repository.GetIncomeByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return income, nil
}

func (s *service) UpdateIncome(ctx context.Context, id int, request UpdateIncomeRequest) (*Income, error) {
	existingIncome, err := s.Repository.GetIncomeByID(ctx, id)
	if err != nil {
		return nil, err
	}

	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		return nil, ErrInvalidIncomeData
	}

	existingIncome.Amount = request.Amount
	existingIncome.Source = request.Source
	existingIncome.Date = date
	existingIncome.UpdatedAt = time.Now()

	updatedIncome, err := s.Repository.UpdateIncome(ctx, existingIncome)
	if err != nil {
		return nil, err
	}

	return updatedIncome, nil
}

func (s *service) DeleteIncome(ctx context.Context, id int) error {
	_, err := s.Repository.GetIncomeByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.Repository.DeleteIncome(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
