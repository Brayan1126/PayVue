package income

func ToIncomeResponse(income *Income) IncomeResponse {
	return IncomeResponse{
		ID:     income.ID,
		Amount: income.Amount,
		Source: income.Source,
		Date:   income.Date.Format("2006-01-02"),
	}
}

func ToIncomeListResponse(incomes []Income) IncomeListResponse {
	responses := make([]IncomeResponse, len(incomes))
	for i, income := range incomes {
		responses[i] = ToIncomeResponse(&income)
	}

	return IncomeListResponse{
		Incomes: responses,
	}
}
