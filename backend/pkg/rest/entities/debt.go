package entities

import "github.com/payvue/payvue-backend/pkg/domain/debt"

type CreateDebtRequest struct {
	Name              string  `json:"name" validate:"required"`
	TotalAmount       float64 `json:"total_amount" validate:"required,gt=0"`
	RemainingAmount   float64 `json:"remaining_amount" validate:"required,gte=0"`
	DueDate           string  `json:"due_date" validate:"required"`
	InterestRate      float64 `json:"interest_rate" validate:"gte=0"`
	NumInstallments   int     `json:"num_installments" validate:"required,gt=0"`
	InstallmentAmount float64 `json:"installment_amount" validate:"required,gt=0"`
	PaymentDay        int     `json:"payment_day" validate:"required,min=1,max=31"`
}

type UpdateDebtRequest struct {
	Name              string  `json:"name" validate:"required"`
	TotalAmount       float64 `json:"total_amount" validate:"required,gt=0"`
	RemainingAmount   float64 `json:"remaining_amount" validate:"required,gte=0"`
	DueDate           string  `json:"due_date" validate:"required"`
	InterestRate      float64 `json:"interest_rate" validate:"gte=0"`
	NumInstallments   int     `json:"num_installments" validate:"required,gt=0"`
	InstallmentAmount float64 `json:"installment_amount" validate:"required,gt=0"`
	PaymentDay        int     `json:"payment_day" validate:"required,min=1,max=31"`
	Paid              bool    `json:"paid"`
}

func (r CreateDebtRequest) ToDomain() debt.CreateDebtRequest {
	return debt.CreateDebtRequest{
		Name:              r.Name,
		TotalAmount:       r.TotalAmount,
		RemainingAmount:   r.RemainingAmount,
		DueDate:           r.DueDate,
		InterestRate:      r.InterestRate,
		NumInstallments:   r.NumInstallments,
		InstallmentAmount: r.InstallmentAmount,
		PaymentDay:        r.PaymentDay,
	}
}

func (r UpdateDebtRequest) ToDomain() debt.UpdateDebtRequest {
	return debt.UpdateDebtRequest{
		Name:              r.Name,
		TotalAmount:       r.TotalAmount,
		RemainingAmount:   r.RemainingAmount,
		DueDate:           r.DueDate,
		InterestRate:      r.InterestRate,
		NumInstallments:   r.NumInstallments,
		InstallmentAmount: r.InstallmentAmount,
		PaymentDay:        r.PaymentDay,
		Paid:              r.Paid,
	}
}
