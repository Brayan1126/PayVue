package income

import (
	"time"
)

type Income struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Amount    float64   `json:"amount"`
	Source    string    `json:"source"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateIncomeRequest struct {
	UserID int     `json:"user_id"`
	Amount float64 `json:"amount" validate:"required,gt=0"`
	Source string  `json:"source" validate:"required"`
	Date   string  `json:"date" validate:"required"`
}

type UpdateIncomeRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
	Source string  `json:"source" validate:"required"`
	Date   string  `json:"date" validate:"required"`
}

type IncomeListResponse struct {
	Incomes []IncomeResponse `json:"incomes"`
}

type IncomeResponse struct {
	ID     int     `json:"id"`
	Amount float64 `json:"amount"`
	Source string  `json:"source"`
	Date   string  `json:"date"`
}
