package payment

import (
	"context"
	"errors"
	"time"
)

var (
	ErrPaymentNotFound    = errors.New("payment not found")
	ErrInvalidPaymentData = errors.New("invalid payment data")
	ErrDatabaseError      = errors.New("database error")
	ErrDebtNotFound       = errors.New("debt not found")
)

type Service interface {
	CreatePayment(ctx context.Context, request CreatePaymentRequest, filename string) (*Payment, error)
	GetAllPayments(ctx context.Context) ([]PaymentWithDebt, error)
	GetPaymentsByUserID(ctx context.Context, userID int) ([]PaymentWithDebt, error)
	GetPaymentByID(ctx context.Context, id int) (*Payment, error)
	DeletePayment(ctx context.Context, id int) error
}

type service struct {
	*Container
}

func New(container *Container) Service {
	return &service{
		Container: container,
	}
}

func (s *service) CreatePayment(ctx context.Context, request CreatePaymentRequest, filename string) (*Payment, error) {
	var date time.Time
	var err error

	if request.Date != "" {
		date, err = time.Parse("2006-01-02", request.Date)
		if err != nil {
			return nil, ErrInvalidPaymentData
		}
	} else {
		date = time.Now()
	}

	payment := &Payment{
		UserID:          request.UserID,
		Amount:          request.Amount,
		DebtID:          request.DebtID,
		ReceiptFilename: filename,
		Date:            date,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	createdPayment, err := s.Repository.CreatePayment(ctx, payment)
	if err != nil {
		return nil, err
	}

	return createdPayment, nil
}

func (s *service) GetAllPayments(ctx context.Context) ([]PaymentWithDebt, error) {
	payments, err := s.Repository.GetAllPayments(ctx)
	if err != nil {
		return nil, err
	}

	return payments, nil
}

func (s *service) GetPaymentsByUserID(ctx context.Context, userID int) ([]PaymentWithDebt, error) {
	payments, err := s.Repository.GetPaymentsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return payments, nil
}

func (s *service) GetPaymentByID(ctx context.Context, id int) (*Payment, error) {
	payment, err := s.Repository.GetPaymentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *service) DeletePayment(ctx context.Context, id int) error {
	_, err := s.Repository.GetPaymentByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.Repository.DeletePayment(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
