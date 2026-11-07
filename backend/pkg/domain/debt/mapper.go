package debt

import "math"

func ToDebtResponse(debt *Debt) DebtResponse {
	remainingPayments := 0
	if debt.InstallmentAmount > 0 {
		remainingPayments = int(math.Floor(debt.RemainingAmount / debt.InstallmentAmount))
	}

	return DebtResponse{
		ID:                debt.ID,
		Name:              debt.Name,
		TotalAmount:       debt.TotalAmount,
		RemainingAmount:   debt.RemainingAmount,
		DueDate:           debt.DueDate.Format("2006-01-02"),
		InterestRate:      debt.InterestRate,
		NumInstallments:   debt.NumInstallments,
		InstallmentAmount: debt.InstallmentAmount,
		PaymentDay:        debt.PaymentDay,
		RemainingPayments: remainingPayments,
		Paid:              debt.Paid,
	}
}

func ToDebtListResponse(debts []Debt) DebtListResponse {
	responses := make([]DebtResponse, len(debts))
	for i, debt := range debts {
		responses[i] = ToDebtResponse(&debt)
	}

	return DebtListResponse{
		Debts: responses,
	}
}
