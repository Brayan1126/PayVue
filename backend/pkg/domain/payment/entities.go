package payment

import (
	"time"
)

type Payment struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	Amount          float64   `json:"amount"`
	DebtID          int       `json:"debt_id"`
	ReceiptFilename string    `json:"receipt_filename"`
	Date            time.Time `json:"date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreatePaymentRequest struct {
	UserID int     `form:"user_id"`
	Amount float64 `form:"amount" validate:"required,gt=0"`
	DebtID int     `form:"debt_id" validate:"required,gt=0"`
	Date   string  `form:"date"`
}

type PaymentListResponse struct {
	Payments []PaymentResponse `json:"payments"`
}

type PaymentResponse struct {
	ID                    int     `json:"id"`
	DebtID                int     `json:"debt_id"`
	Amount                float64 `json:"amount"`
	Date                  string  `json:"date"`
	CreatedAt             string  `json:"created_at"`
	DebtName              string  `json:"debt_name"`
	RemainingInstallments int     `json:"remaining_installments"`
	RemainingAmount       float64 `json:"remaining_amount"`
	ReceiptURL            string  `json:"receipt_url"`
}
