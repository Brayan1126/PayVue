package debt

import (
	"time"
)

type Debt struct {
	ID                int       `json:"id"`
	UserID            int       `json:"user_id"`
	Name              string    `json:"name"`
	TotalAmount       float64   `json:"total_amount"`
	RemainingAmount   float64   `json:"remaining_amount"`
	DueDate           time.Time `json:"due_date"`
	InterestRate      float64   `json:"interest_rate"`
	NumInstallments   int       `json:"num_installments"`
	InstallmentAmount float64   `json:"installment_amount"`
	PaymentDay        int       `json:"payment_day"`
	Paid              bool      `json:"paid"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CreateDebtRequest struct {
	UserID            int     `json:"user_id"`
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

type DebtListResponse struct {
	Debts []DebtResponse `json:"debts"`
}

type DebtResponse struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	TotalAmount       float64 `json:"total_amount"`
	RemainingAmount   float64 `json:"remaining_amount"`
	DueDate           string  `json:"due_date"`
	InterestRate      float64 `json:"interest_rate"`
	NumInstallments   int     `json:"num_installments"`
	InstallmentAmount float64 `json:"installment_amount"`
	PaymentDay        int     `json:"payment_day"`
	RemainingPayments int     `json:"remaining_payments"`
	Paid              bool    `json:"paid"`
}
