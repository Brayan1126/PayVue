package entities

import "github.com/payvue/payvue-backend/pkg/domain/income"

type CreateIncomeRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
	Source string  `json:"source" validate:"required"`
	Date   string  `json:"date" validate:"required"`
}

type UpdateIncomeRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
	Source string  `json:"source" validate:"required"`
	Date   string  `json:"date" validate:"required"`
}

func (r CreateIncomeRequest) ToDomain() income.CreateIncomeRequest {
	return income.CreateIncomeRequest{
		Amount: r.Amount,
		Source: r.Source,
		Date:   r.Date,
	}
}

func (r UpdateIncomeRequest) ToDomain() income.UpdateIncomeRequest {
	return income.UpdateIncomeRequest{
		Amount: r.Amount,
		Source: r.Source,
		Date:   r.Date,
	}
}
